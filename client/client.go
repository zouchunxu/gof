package client

import (
	"github.com/zouchunxu/gof/middlewares/prometheus"
	"google.golang.org/grpc"
)

func NewConnect(target string) (*grpc.ClientConn, error) {
	//TODO trace
	return grpc.Dial(target,
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(
			prometheus.UnaryClientInterceptor,
		),
	)
}
