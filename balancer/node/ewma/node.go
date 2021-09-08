package ewma

import (
	"container/list"
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"errors"
	"github.com/zouchunxu/gof/balancer"
)

const (

	// The mean lifetime of `cost`, it reaches its half-life after Tau*ln(2).
	tau = int64(time.Millisecond * 600)
	// if statistic not collected,we add a big lag penalty to endpoint
	penalty = uint64(time.Second * 10)
)

var (
	_ balancer.Node        = &node{}
	_ balancer.NodeBuilder = &Builder{}
)

// node is endpoint instance
type node struct {
	addr     string
	version  string
	metadata balancer.Metadata

	// client statistic data
	lag       int64
	success   uint64
	inflight  int64
	inflights *list.List
	// last collected timestamp
	stamp     int64
	predictTs int64
	predict   int64
	// request number in a period time
	reqs int64
	// last lastPick timestamp
	lastPick int64

	errHandler func(err error) (isErr bool)
	lk         sync.RWMutex
}

// Builder is ewma node builder
type Builder struct {
	ErrHandler func(err error) (isErr bool)
}

// Build create node
func (b *Builder) Build(addr string, initWeight float64, metadata balancer.Metadata) balancer.Node {
	s := &node{
		addr:       addr,
		metadata:   metadata,
		lag:        0,
		success:    1000,
		inflight:   1,
		inflights:  list.New(),
		errHandler: b.ErrHandler,
	}
	return s
}

func (n *node) health() uint64 {
	return atomic.LoadUint64(&n.success)
}

func (n *node) load() (load uint64) {
	now := time.Now().UnixNano()
	avgLag := atomic.LoadInt64(&n.lag)

	lastPredictTs := atomic.LoadInt64(&n.predictTs)
	predicInterval := avgLag / 5
	if predicInterval < int64(time.Millisecond*5) {
		predicInterval = int64(time.Millisecond * 5)
	} else if predicInterval > int64(time.Millisecond*200) {
		predicInterval = int64(time.Millisecond * 200)
	}
	if now-lastPredictTs > predicInterval {
		if atomic.CompareAndSwapInt64(&n.predictTs, lastPredictTs, now) {
			var (
				total   int64
				count   int
				predict int64
			)
			n.lk.RLock()
			first := n.inflights.Front()
			for {
				if first == nil {
					break
				}
				lag := now - first.Value.(int64)
				if lag > avgLag {
					count++
					total += lag
				}
				first = first.Next()
			}
			if count > (n.inflights.Len()/2 + 1) {
				predict = total / int64(count)
			}
			n.lk.RUnlock()
			atomic.StoreInt64(&n.predict, predict)
		}
	}

	if avgLag == 0 {
		// penalty是node刚启动时没有数据时的惩罚值，默认为1e9 * 10
		load = penalty * uint64(atomic.LoadInt64(&n.inflight))
	} else {
		predict := atomic.LoadInt64(&n.predict)
		if predict > avgLag {
			avgLag = predict
		}
		load = uint64(avgLag) * uint64(atomic.LoadInt64(&n.inflight))
	}
	return
}

// Pick choose node
func (n *node) Pick() balancer.Done {
	now := time.Now().UnixNano()
	atomic.StoreInt64(&n.lastPick, now)
	atomic.AddInt64(&n.inflight, 1)
	atomic.AddInt64(&n.reqs, 1)
	n.lk.Lock()
	e := n.inflights.PushBack(now)
	n.lk.Unlock()

	return func(ctx context.Context, di balancer.DoneInfo) {
		n.lk.Lock()
		n.inflights.Remove(e)
		n.lk.Unlock()
		atomic.AddInt64(&n.inflight, -1)

		now := time.Now().UnixNano()
		// get moving average ratio w
		stamp := atomic.SwapInt64(&n.stamp, now)
		td := now - stamp
		if td < 0 {
			td = 0
		}
		w := math.Exp(float64(-td) / float64(tau))

		start := e.Value.(int64)
		lag := now - start
		if lag < 0 {
			lag = 0
		}
		oldLag := atomic.LoadInt64(&n.lag)
		if oldLag == 0 {
			w = 0.0
		}
		lag = int64(float64(oldLag)*w + float64(lag)*(1.0-w))
		atomic.StoreInt64(&n.lag, lag)

		success := uint64(1000) // error value ,if error set 1
		if di.Err != nil {
			if n.errHandler != nil {
				if n.errHandler(di.Err) {
					success = 0
				}
			} else if errors.Is(context.DeadlineExceeded, di.Err) || errors.Is(context.Canceled, di.Err) {
				success = 0
			}
		}
		oldSuc := atomic.LoadUint64(&n.success)
		success = uint64(float64(oldSuc)*w + float64(success)*(1.0-w))
		atomic.StoreUint64(&n.success, success)
	}
}

// Weight is node effective weight
func (n *node) Weight() (weight float64) {
	weight = float64(n.health()*uint64(time.Second)) / float64(n.load())
	return
}

func (n *node) PickElapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - atomic.LoadInt64(&n.lastPick))
}

func (n *node) Address() string {
	return n.addr
}

func (n *node) Version() string {
	return n.version
}

func (n *node) Metadata() balancer.Metadata {
	return n.metadata
}
