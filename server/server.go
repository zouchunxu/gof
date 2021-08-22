package server

import (
	"github.com/opentracing/opentracing-go"
	opentracing3 "github.com/zouchunxu/gof/middlewares/opentracing"
	"github.com/zouchunxu/gof/middlewares/prometheus"
	"google.golang.org/grpc"
	"net"
)

var s *grpc.Server

var mid []grpc.UnaryServerInterceptor

func Init() *grpc.Server {
	s = grpc.NewServer()
	mid = append(mid, opentracing3.OpentracingServerInterceptor(opentracing.GlobalTracer()))
	mid = append(mid, prometheus.UnaryServerInterceptor)
	grpc.ChainUnaryInterceptor(mid...)
	return s
}

func Run(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}
