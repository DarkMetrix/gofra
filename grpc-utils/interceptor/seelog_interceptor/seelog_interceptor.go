package seelog_interceptor

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	log "github.com/cihub/seelog"
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

	ctxInfo := fmt.Sprintf("%v", ctx)
	log.Tracef("context:%v", ctxInfo)

	if err != nil {
		log.Warnf("invoke failed! req=%v, error:%v", req, err)
	} else {
		log.Debugf("invoke success! req=%v, reply:%v", req, reply)
	}

	return err
}

// seelog server interceptor
var GofraServerInterceptor grpc.UnaryServerInterceptor = GofraServerInterceptorFunc

func GofraServerInterceptorFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
	// Process
	reply, err = handler(ctx, req)

	ctxInfo := fmt.Sprintf("%v", ctx)
	log.Tracef("context:%v", ctxInfo)

	if err != nil {
		log.Warnf("handle failed! req=%v, error:%v", req, err)
	} else {
		log.Debugf("handle success! req=%v, reply:%v", req, reply)
	}

	return reply, err
}
