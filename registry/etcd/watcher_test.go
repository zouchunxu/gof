package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	w := newWatcher(context.Background(), "zcx", client)
	w.Next()
	for {
		ips, err := w.Next()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(ips)
	}
}
