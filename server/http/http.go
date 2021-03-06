package http

import (
	"context"
	"github.com/zouchunxu/gof/internal/endpoint"
	"github.com/zouchunxu/gof/internal/host"
	"net"
	"net/http"
	"net/url"
	"sync"
)

type HttpServer struct {
	s        *http.Server
	mux      *http.ServeMux
	once     sync.Once
	endpoint *url.URL
	addr     string
}

func NewHttpServer(addr string) *HttpServer {
	h := &HttpServer{
		s:    &http.Server{},
		mux:  http.NewServeMux(),
		addr: addr,
	}
	h.s.Handler = h.mux
	return h
}

func (h *HttpServer) AddRouterFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	h.mux.HandleFunc(pattern, handler)
}

func (h *HttpServer) AddRouter(pattern string, handler http.Handler) {
	h.mux.Handle(pattern, handler)
}

func (h *HttpServer) Start(context.Context) error {
	l, err := net.Listen("tcp", h.addr)
	if err != nil {
		return err
	}
	return h.s.Serve(l)
}

func (h *HttpServer) Stop(context.Context) error {
	return h.s.Shutdown(context.Background())
}

// Endpoint return a real address to registry endpoint.
// examples:
//   grpc://127.0.0.1:9000?isSecure=false
func (h *HttpServer) Endpoint() (*url.URL, error) {
	var err error
	h.once.Do(func() {
		if h.endpoint != nil {
			return
		}
		lis, err := net.Listen("tcp", h.addr)
		if err != nil {
			return
		}
		addr, err := host.Extract(h.addr, lis)
		if err != nil {
			lis.Close()
			return
		}
		h.endpoint = endpoint.NewEndpoint("http", addr, false)
	})
	return h.endpoint, err
}
