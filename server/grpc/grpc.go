package grpc

import (
	"github.com/zouchunxu/gof/internal/endpoint"
	"github.com/zouchunxu/gof/internal/host"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
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
