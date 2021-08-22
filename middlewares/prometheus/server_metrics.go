package prometheus

import (
	"context"
	prom "github.com/prometheus/client_golang/prometheus"
	grpcstatus2 "github.com/zouchunxu/gof/middlewares/prometheus/grpcstatus"
	"google.golang.org/grpc"
)

// ServerMetrics represents a collection of metrics to be registered on a
// Prometheus metrics registry for a gRPC server.
type ServerMetrics struct {
	serverStartedCounter          *prom.CounterVec
	serverHandledCounter          *prom.CounterVec
	serverStreamMsgReceived       *prom.CounterVec
	serverStreamMsgSent           *prom.CounterVec
	serverHandledHistogramEnabled bool
	serverHandledHistogramOpts    prom.HistogramOpts
	serverHandledHistogram        *prom.HistogramVec
}

func NewServerMetrics(counterOpts ...CounterOption) *ServerMetrics {
	opts := counterOptions(counterOpts)
	return &ServerMetrics{
		serverStartedCounter: prom.NewCounterVec(
			opts.apply(prom.CounterOpts{
				Name: "grpc_server_started_total",
				Help: "Total number of RPCs started on the server.",
			}), []string{"grpc_type", "grpc_service", "grpc_method"}),
		serverHandledCounter: prom.NewCounterVec(
			opts.apply(prom.CounterOpts{
				Name: "grpc_server_handled_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}), []string{"grpc_type", "grpc_service", "grpc_method", "grpc_code"}),
		serverStreamMsgReceived: prom.NewCounterVec(
			opts.apply(prom.CounterOpts{
				Name: "grpc_server_msg_received_total",
				Help: "Total number of RPC stream messages received on the server.",
			}), []string{"grpc_type", "grpc_service", "grpc_method"}),
		serverStreamMsgSent: prom.NewCounterVec(
			opts.apply(prom.CounterOpts{
				Name: "grpc_server_msg_sent_total",
				Help: "Total number of gRPC stream messages sent by the server.",
			}), []string{"grpc_type", "grpc_service", "grpc_method"}),
		serverHandledHistogramEnabled: false,
		serverHandledHistogramOpts: prom.HistogramOpts{
			Name:    "grpc_server_handling_seconds",
			Help:    "Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.",
			Buckets: prom.DefBuckets,
		},
		serverHandledHistogram: nil,
	}
}

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
func (m *ServerMetrics) UnaryServerInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		monitor := newServerReporter(m, Unary, info.FullMethod)
		monitor.ReceivedMessage()
		resp, err := handler(ctx, req)
		st, _ := grpcstatus2.FromError(err)
		monitor.Handled(st.Code())
		if err == nil {
			monitor.SentMessage()
		}
		return resp, err
	}
}

// InitializeMetrics initializes all metrics, with their appropriate null
// value, for all gRPC methods registered on a gRPC server. This is useful, to
// ensure that all metrics exist when collecting and querying.
func (m *ServerMetrics) InitializeMetrics(server *grpc.Server) {
	serviceInfo := server.GetServiceInfo()
	for serviceName, info := range serviceInfo {
		for _, mInfo := range info.Methods {
			preRegisterMethod(m, serviceName, &mInfo)
		}
	}
}

// preRegisterMethod is invoked on Register of a Server, allowing all gRPC services labels to be pre-populated.
func preRegisterMethod(metrics *ServerMetrics, serviceName string, mInfo *grpc.MethodInfo) {
	methodName := mInfo.Name
	methodType := string(typeFromMethodInfo(mInfo))
	// These are just references (no increments), as just referencing will create the labels but not set values.
	metrics.serverStartedCounter.GetMetricWithLabelValues(methodType, serviceName, methodName)
	metrics.serverStreamMsgReceived.GetMetricWithLabelValues(methodType, serviceName, methodName)
	metrics.serverStreamMsgSent.GetMetricWithLabelValues(methodType, serviceName, methodName)
	if metrics.serverHandledHistogramEnabled {
		metrics.serverHandledHistogram.GetMetricWithLabelValues(methodType, serviceName, methodName)
	}
	for _, code := range allCodes {
		metrics.serverHandledCounter.GetMetricWithLabelValues(methodType, serviceName, methodName, code.String())
	}
}

// EnableHandlingTimeHistogram enables histograms being registered when
// registering the ServerMetrics on a Prometheus registry. Histograms can be
// expensive on Prometheus servers. It takes options to configure histogram
// options such as the defined buckets.
func (m *ServerMetrics) EnableHandlingTimeHistogram(opts ...HistogramOption) {
	for _, o := range opts {
		o(&m.serverHandledHistogramOpts)
	}
	if !m.serverHandledHistogramEnabled {
		m.serverHandledHistogram = prom.NewHistogramVec(
			m.serverHandledHistogramOpts,
			[]string{"grpc_type", "grpc_service", "grpc_method"},
		)
	}
	m.serverHandledHistogramEnabled = true
}
