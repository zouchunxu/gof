package client

import (
	"github.com/zouchunxu/gof/middlewares/prometheus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"strings"
)

//NewConnect etcd:///foo/bar/my-service
func NewConnect(target string) (*grpc.ClientConn, error) {
	//TODO trace
	if strings.Contains(target, "etcd") {
		cli, cerr := clientv3.NewFromURL("http://127.0.0.1:2379")
		if cerr != nil {
			return nil, cerr
		}
		etcdResolver, err := resolver.NewBuilder(cli)
		if err != nil {
			return nil, err
		}
		conn, gerr := grpc.Dial(target,
			grpc.WithResolvers(etcdResolver),
			grpc.WithInsecure(),
			grpc.WithChainUnaryInterceptor(
				prometheus.UnaryClientInterceptor,
			),
		)
		if gerr != nil {
			return nil, gerr
		}
		return conn, nil
	}

	return grpc.Dial(target,
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(
			prometheus.UnaryClientInterceptor,
		),
	)
}
