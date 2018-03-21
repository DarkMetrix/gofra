package seelog_interceptor

import (
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
	log.Debugf("====== Enter seelog client interceptor ======")
	defer log.Debugf("====== Leave seelog client interceptor ======")

	log.Debugf("context:%v", ctx)
	log.Debugf("method:%v", method)
	log.Debugf("options:%v", opts)
	log.Debugf("req:%v", req)

	// Invoke remote
	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		log.Warnf("invoke failed!!! error:%v", err)
	} else {
		log.Debugf("reply:%v", reply)
	}

	return err
}

// seelog server interceptor
var GofraServerInterceptor grpc.UnaryServerInterceptor = GofraServerInterceptorFunc

func GofraServerInterceptorFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
	log.Debugf("====== Enter seelog server interceptor ======")
	defer log.Debugf("====== Leave seelog server interceptor ======")

	log.Debugf("context:%v", ctx)
	log.Debugf("method:%v", info.FullMethod)
	log.Debugf("server:%v", info.Server)
	log.Debugf("req:%v", req)

	// Process
	reply, err = handler(ctx, req)

	if err != nil {
		log.Warnf("handle failed!!! error:%v", err)
	} else {
		log.Debugf("reply:%v", reply)
	}

	return reply, err
}