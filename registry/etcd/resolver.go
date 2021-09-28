package etcd

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"time"
)

type builder struct {
	client *clientv3.Client
}

func NewBuilder(client *clientv3.Client) *builder {
	return &builder{
		client: client,
	}
}

func (b *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	var (
		err error
		w   *watcher
	)
	done := make(chan bool, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		key := fmt.Sprintf("%s/%s", namespace, target.Endpoint)
		w = newWatcher(ctx, key, b.client)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second * 10):
		err = errors.New("timeout")
	}
	if err != nil {
		cancel()
		return nil, err
	}
	r := &gofResolver{
		w:      w,
		cc:     cc,
		ctx:    ctx,
		cancel: cancel,
	}
	go r.watch()
	return r, nil
}

// Scheme return scheme of discovery
func (*builder) Scheme() string {
	return "etcd"
}

type gofResolver struct {
	ctx    context.Context
	cancel context.CancelFunc
	w      *watcher
	cc     resolver.ClientConn
}

func (g *gofResolver) watch() {
	for {
		select {
		case <-g.ctx.Done():
			return
		default:
		}
		ips, err := g.w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			time.Sleep(time.Second)
			continue
		}
		//fmt.Printf("ips: %+v\n", ips)
		g.update(ips)
	}
}

func (g *gofResolver) update(ips []string) {
	addrs := make([]resolver.Address, 0)
	for _, ip := range ips {
		//endpoint, err := endpoint2.ParseEndpoint([]string{ip}, "grpc", false)
		//fmt.Println("endpoint: " + endpoint)
		//if err != nil {
		//	log.Error(err)
		//	continue
		//}
		//if endpoint == "" {
		//	continue
		//}
		addr := resolver.Address{
			ServerName: ip,
			Addr:       ip,
		}
		addrs = append(addrs, addr)
	}
	if len(addrs) == 0 {
		return
	}
	err := g.cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		log.Error(err)
	}
}

func (g *gofResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (g *gofResolver) Close() {
	g.cancel()
	err := g.w.Stop()
	if err != nil {
		log.Error(err)
	}
}
