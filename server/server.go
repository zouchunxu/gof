package server

import (
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	opentracing3 "github.com/zouchunxu/gof/middlewares/opentracing"
	"github.com/zouchunxu/gof/middlewares/prometheus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

//Server server，框架目前未引入配置，所以一些值先固定死
type Server struct {
	GrpcSever        *grpc.Server
	Mid              []grpc.UnaryServerInterceptor
	Name             string
	Log              *logrus.Logger
	prrofServer      *http.Server
	prometheusServer *http.Server
}

//Init init
func (s *Server) Init(name string) *grpc.Server {
	s.Name = name
	s.GrpcSever = grpc.NewServer(
		grpc.ChainUnaryInterceptor(s.Mid...),
	)
	s.Mid = append(s.Mid, opentracing3.OpentracingServerInterceptor(opentracing.GlobalTracer()))
	s.Mid = append(s.Mid, prometheus.UnaryServerInterceptor)
	s.initLog()
	s.initJaeger()
	s.initPprof()
	s.initPrometheus()
	return s.GrpcSever
}

//Run 运行server
func (s *Server) Run(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return s.GrpcSever.Serve(lis)
}

//initJaeger 初始化jaeger
func (s *Server) initJaeger() {
	sampler := &config.SamplerConfig{
		Type:  "const",
		Param: 1,
	}
	cfg := &config.Configuration{
		ServiceName: s.Name,
		Sampler:     sampler,
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			LocalAgentHostPort:  "127.0.0.1:9502",
			BufferFlushInterval: 1 * time.Second,
			QueueSize:           200,
		},
	}
	tracer, _, err := cfg.NewTracer(
		config.Logger(jaeger.NullLogger),
		config.Metrics(metrics.NullFactory),
	)
	if err != nil {
		panic(err.Error())
	}
	opentracing.SetGlobalTracer(tracer)
}

//initPprof 初始化pprof
func (s *Server) initPprof() {
	mux := &http.ServeMux{}
	s.prrofServer = &http.Server{
		Addr:    ":9909",
		Handler: mux,
	}
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

//initPrometheus 初始化prometheus
func (s *Server) initPrometheus() {
	mux := &http.ServeMux{}
	s.prometheusServer = &http.Server{
		Addr:    ":9910",
		Handler: mux,
	}
	mux.Handle("/metrics", promhttp.Handler())
}

//initLog 初始化日志
func (s *Server) initLog() {
	s.Log = logrus.New()
	s.Log.SetFormatter(&logrus.JSONFormatter{})
	s.Log.SetOutput(os.Stdout)
	s.Log.SetLevel(logrus.InfoLevel)
}
