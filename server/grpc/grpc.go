package grpc

import (
	"context"
	"github.com/zouchunxu/gof/internal/endpoint"
	"github.com/zouchunxu/gof/internal/host"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"
	"net"
	"net/url"
	"sync"
)

type GrpcServer struct {
	*grpc.Server
	addr     string
	health   *health.Server
	once     sync.Once
	endpoint *url.URL
}

func NewGrpcServer(addr string, mid ...grpc.UnaryServerInterceptor) *GrpcServer {
	g := &GrpcServer{
		Server: grpc.NewServer(grpc.ChainUnaryInterceptor(mid...)),
		addr:   addr,
		health: health.NewServer(),
	}
	reflection.Register(g.Server)
	return g
}

func (g *GrpcServer) Start(context.Context) error {
	lis, err := net.Listen("tcp", g.addr)
	if err != nil {
		return err
	}
	g.health.Resume()
	return g.Serve(lis)
}

func (g *GrpcServer) Stop(context.Context) error {
	g.GracefulStop()
	g.health.Shutdown()
	return nil
}

// Endpoint return a real address to registry endpoint.
// examples:
//   grpc://127.0.0.1:9000?isSecure=false
func (g *GrpcServer) Endpoint() (*url.URL, error) {
	var err error
	g.once.Do(func() {
		if g.endpoint != nil {
			return
		}
		lis, err := net.Listen("tcp", g.addr)
		if err != nil {
			return
		}
		addr, err := host.Extract(g.addr, lis)
		if err != nil {
			lis.Close()
			return
		}
		g.endpoint = endpoint.NewEndpoint("grpc", addr, false)
	})
	return g.endpoint, err
}
