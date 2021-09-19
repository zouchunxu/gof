package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/apimachinery/pkg/util/json"
)

type watcher struct {
	key       string
	ctx       context.Context
	cancel    context.CancelFunc
	watchChan clientv3.WatchChan
	watcher   clientv3.Watcher
	kv        clientv3.KV
	first     bool
}

func newWatcher(ctx context.Context, key string, client *clientv3.Client) *watcher {
	w := &watcher{
		key:     key,
		watcher: clientv3.NewWatcher(client),
		kv:      clientv3.NewKV(client),
		first:   true,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.watchChan = w.watcher.Watch(w.ctx, key, clientv3.WithPrefix(), clientv3.WithRev(0))
	w.watcher.RequestProgress(context.Background())
	return w
}

func (w *watcher) Next() ([]string, error) {
	if w.first {
		item, err := w.getIps()
		w.first = false
		return item, err
	}
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.watchChan:
		return w.getIps()
	}
}

func (w *watcher) Stop() error {
	w.ctx.Done()
	return w.watcher.Close()
}

func (w *watcher) getIps() ([]string, error) {
	fmt.Println("key: " + w.key)
	resp, err := w.kv.Get(w.ctx, w.key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var res []string
	for _, kv := range resp.Kvs {
		var ips []string
		err = json.Unmarshal(kv.Value, &ips)
		fmt.Println(ips)
		if err != nil {
			return nil, err
		}
		res = append(res, ips...)
	}
	return res, nil
}
