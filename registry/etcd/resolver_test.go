package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"sync"
	"testing"
	"time"
)

func TestConn(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	r := New(client)
	b := NewBuilder(r)

	_, err = b.Build(resolver.Target{
		Endpoint: "zcx",
	}, &testClientConn{}, resolver.BuildOptions{})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)
}

type testClientConn struct {
	resolver.ClientConn // For unimplemented functions
	target              string
	m1                  sync.Mutex
	state               resolver.State
	updateStateCalls    int
	errChan             chan error
	updateStateErr      error
}

func (t *testClientConn) UpdateState(s resolver.State) error {
	t.m1.Lock()
	defer t.m1.Unlock()
	t.state = s
	t.updateStateCalls++
	// This error determines whether DNS Resolver actually decides to exponentially backoff or not.
	// This can be any error.
	return t.updateStateErr
}

func (t *testClientConn) getState() (resolver.State, int) {
	t.m1.Lock()
	defer t.m1.Unlock()
	return t.state, t.updateStateCalls
}
