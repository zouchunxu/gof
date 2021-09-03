package server

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	cfg "github.com/zouchunxu/gof/config"
	opentracing3 "github.com/zouchunxu/gof/middlewares/opentracing"
	"github.com/zouchunxu/gof/middlewares/prometheus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

//App server
type App struct {
	GrpcSever        *grpc.Server
	Mid              []grpc.UnaryServerInterceptor
	Log              *logrus.Logger
	DB               *gorm.DB
	prrofServer      *http.Server
	prometheusServer *http.Server
	c                cfg.System
	path             string
	G                *gin.Engine
}

//New init
func New(path string) *App {
	s := &App{path: path}
	s.initConfig()
	//s.c = cfg.System{}
	s.GrpcSever = grpc.NewServer(
		grpc.ChainUnaryInterceptor(s.Mid...),
	)
	s.Mid = append(s.Mid, opentracing3.OpentracingServerInterceptor(opentracing.GlobalTracer()))
	s.Mid = append(s.Mid, prometheus.UnaryServerInterceptor)
	s.initLog()
	s.initJaeger()
	s.initPprof()
	s.initPrometheus()
	s.initDB()
	s.initHttpServer()
	return s
}

//Run 运行server
func (s *App) Run() error {
	go func() {
		lis, _ := net.Listen("tcp", s.c.PrometheusHost)
		s.prometheusServer.Serve(lis)
	}()
	go func() {
		lis, _ := net.Listen("tcp", s.c.PprofHost)
		s.prrofServer.Serve(lis)
	}()
	go func() {
		s.G.Run(s.c.HttpPort)
	}()
	lis, err := net.Listen("tcp", s.c.ServerPort)
	if err != nil {
		return err
	}
	return s.GrpcSever.Serve(lis)
}

//initJaeger 初始化jaeger
func (s *App) initJaeger() {
	sampler := &config.SamplerConfig{
		Type:  "const",
		Param: 1,
	}
	jcfg := &config.Configuration{
		ServiceName: s.c.Name,
		Sampler:     sampler,
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			LocalAgentHostPort:  s.c.JaegerHost,
			BufferFlushInterval: 1 * time.Second,
			QueueSize:           200,
		},
	}
	tracer, _, err := jcfg.NewTracer(
		config.Logger(jaeger.NullLogger),
		config.Metrics(metrics.NullFactory),
	)
	if err != nil {
		panic(err.Error())
	}
	opentracing.SetGlobalTracer(tracer)
}

//initPprof 初始化pprof
func (s *App) initPprof() {
	mux := &http.ServeMux{}
	s.prrofServer = &http.Server{
		Addr:    s.c.PprofHost,
		Handler: mux,
	}
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

//initPrometheus 初始化prometheus
func (s *App) initPrometheus() {
	mux := &http.ServeMux{}
	s.prometheusServer = &http.Server{
		Addr:    s.c.PrometheusHost,
		Handler: mux,
	}
	mux.Handle("/metrics", promhttp.Handler())
}

//initLog 初始化日志
func (s *App) initLog() {
	s.Log = logrus.New()
	s.Log.SetFormatter(&logrus.JSONFormatter{})
	s.Log.SetOutput(os.Stdout)
	s.Log.SetLevel(logrus.InfoLevel)
}

//initDB 初始化数据库
func (s *App) initDB() {
	//db, err := gorm.Open(mysql.New(mysql.Config{
	//	DSN:                       s.c.DSN,
	//	DefaultStringSize:         256,
	//	DisableDatetimePrecision:  true,
	//	DontSupportRenameIndex:    true,
	//	DontSupportRenameColumn:   true,
	//	SkipInitializeWithVersion: false,
	//}), &gorm.Config{})
	//if err != nil {
	//	panic("failed to connect database")
	//}
	//s.DB = db
}

func (s *App) initConfig() {
	viper.SetConfigFile(s.path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}
	if err := viper.Unmarshal(&s.c); err != nil {
		panic(err.Error())
	}
}

func (s *App) initHttpServer() {
	s.G = gin.Default()
}
