package server

import (
	"context"
	"errors"
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
	zgrpc "github.com/zouchunxu/gof/server/grpc"
	zhttp "github.com/zouchunxu/gof/server/http"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

//App server
type App struct {
	GrpcSever        *zgrpc.GrpcServer
	Mid              []grpc.UnaryServerInterceptor
	Log              *logrus.Logger
	DB               *gorm.DB
	prrofServer      *http.Server
	prometheusServer *http.Server
	c                cfg.System
	path             string
	G                *gin.Engine
	servers          []Server
}

//New init
func New(path string) *App {
	s := &App{path: path}
	s.initConfig()
	//s.c = cfg.System{}
	s.Mid = append(s.Mid, opentracing3.OpentracingServerInterceptor(opentracing.GlobalTracer()))
	s.Mid = append(s.Mid, prometheus.UnaryServerInterceptor)
	s.GrpcSever = zgrpc.NewGrpcServer(s.c.ServerPort, s.Mid...)
	s.servers = append(s.servers, s.GrpcSever)
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
	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	for _, srv := range s.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			return srv.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}
	wg.Wait()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				err := s.Stop()
				if err != nil {
					s.Log.Errorf("failed to app stop: %v", err)
				}
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

//Stop todo
func (s *App) Stop() error {
	return nil
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
	svr := zhttp.NewHttpServer(s.c.PprofHost)
	svr.AddRouterFunc("/debug/pprof/", pprof.Index)
	svr.AddRouterFunc("/debug/pprof/cmdline", pprof.Cmdline)
	svr.AddRouterFunc("/debug/pprof/profile", pprof.Profile)
	svr.AddRouterFunc("/debug/pprof/symbol", pprof.Symbol)
	svr.AddRouterFunc("/debug/pprof/trace", pprof.Trace)
	s.servers = append(s.servers, svr)
}

//initPrometheus 初始化prometheus
func (s *App) initPrometheus() {
	svr := zhttp.NewHttpServer(s.c.PrometheusHost)
	svr.AddRouter("/metrics", promhttp.Handler())
	s.servers = append(s.servers, svr)
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
