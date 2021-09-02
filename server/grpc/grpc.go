package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"net"
)

type GrpcServer struct {
	*grpc.Server
	addr   string
	health *health.Server
}

func NewGrpcServer(addr string) *GrpcServer {
	g := &GrpcServer{
		Server: grpc.NewServer(),
		addr:   addr,
		health: health.NewServer(),
	}
	return g
}

func (g *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", g.addr)
	if err != nil {
		return err
	}
	g.health.Resume()
	return g.Serve(lis)
}

func (g *GrpcServer) Stop() error {
	g.GracefulStop()
	g.health.Shutdown()
	return nil
}
