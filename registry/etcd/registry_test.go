package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestHeartBeat(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	ctx := context.Background()

	r := New(client)
	r.lease = clientv3.NewLease(client)
	leaseID, err := r.registerWithKV(ctx, "a", "b")

	fmt.Println(leaseID)
	fmt.Println(r.GetService(ctx, "a"))
	r.heartBeat(ctx, leaseID, "a", "c")
}
