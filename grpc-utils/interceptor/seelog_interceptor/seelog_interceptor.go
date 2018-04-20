package seelog_interceptor

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	log "github.com/cihub/seelog"

	tracing "github.com/DarkMetrix/gofra/common/tracing/zipkin"
)

func GetClientInterceptor() grpc.UnaryClientInterceptor {
	return GofraClientInterceptor
}

func GetServerInterceptor() grpc.UnaryServerInterceptor {
	return GofraServerInterceptor
}

// seelog client interceptor
var GofraClientInterceptor grpc.UnaryClientInterceptor = GofraClientInterceptorFunc

func GofraClientInterceptorFunc(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Invoke remote
	err := invoker(ctx, method, req, reply, cc, opts...)

	log.Tracef("context:%v", ctx)

	if err != nil {
		log.Warnf("invoke failed! trace id=%v, span id=%v, req=%v, error:%v", tracing.GetTracingId(ctx), tracing.GetSpanId(ctx), req, err)
	} else {
		log.Debugf("invoke success! trace id=%v, span id=%v, req=%v, reply:%v", tracing.GetTracingId(ctx), tracing.GetSpanId(ctx), req, reply)
	}

	return err
}

// seelog server interceptor
var GofraServerInterceptor grpc.UnaryServerInterceptor = GofraServerInterceptorFunc

func GofraServerInterceptorFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
	// Process
	reply, err = handler(ctx, req)

	log.Tracef("context:%v", ctx)

	if err != nil {
		log.Warnf("handle failed! trace id=%v, span id=%v, req=%v, error:%v", tracing.GetTracingId(ctx), tracing.GetSpanId(ctx), req, err)
	} else {
		log.Debugf("handle success! trace id=%v, span id=%v, req=%v, reply:%v", tracing.GetTracingId(ctx), tracing.GetSpanId(ctx), req, reply)
	}

	return reply, err
}
