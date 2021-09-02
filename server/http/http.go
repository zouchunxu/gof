package http

import (
	"context"
	"net/http"
)

type HttpServer struct {
	s   *http.Server
	mux *http.ServeMux
}

func NewHttpServer(addr string) {
	h := HttpServer{
		s:   &http.Server{},
		mux: http.NewServeMux(),
	}
	h.s.Handler = h.mux
}

func (h *HttpServer) AddRouter(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	h.mux.HandleFunc(pattern, handler)
}

func (h *HttpServer) Start() error {
	return h.s.ListenAndServe()
}

func (h *HttpServer) Stop() error {
	return h.s.Shutdown(context.Background())
}
