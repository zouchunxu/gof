package opentracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

//MDReaderWriter metadata Reader and Writer
type MDReaderWriter struct {
	metadata.MD
}

var (
	grpcTag = opentracing.Tag{Key: string(ext.Component), Value: "gRPC"}
)

// OpentracingServerInterceptor rewrite server's interceptor with open tracing
func OpentracingServerInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		parentSpanContext, err := tracer.Extract(opentracing.HTTPHeaders, metadataTextMap(md))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			grpclog.Infof("grpc_opentracing: failed parsing trace information: %v", err)
		}
		serverSpan := tracer.StartSpan(info.FullMethod,
			// 这里实现的关联父级span，如果partentSpanContext为空，则是顶级
			ext.RPCServerOption(parentSpanContext),
			grpcTag,
		)
		injectOpentracingIdsToTags(serverSpan, Extract(ctx))
		newCtx := opentracing.ContextWithSpan(ctx, serverSpan)
		rsp, err := handler(newCtx, req)
		tags := Extract(newCtx)
		for key, val := range tags.Values() {
			// Don't tag errors, log them instead.
			if vErr, ok := val.(error); ok {
				serverSpan.LogKV(key, vErr.Error())

			} else {
				serverSpan.SetTag(key, val)
			}
		}
		if err != nil {
			ext.Error.Set(serverSpan, true)
			serverSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}
		return rsp, err
	}
}
