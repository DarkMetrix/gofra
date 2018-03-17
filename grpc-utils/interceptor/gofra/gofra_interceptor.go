package gofra

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	monitor "github.com/DarkMetrix/gofra/grpc-utils/monitor/statsd"
	log "github.com/cihub/seelog"
)

func GetClientInterceptor() grpc.UnaryClientInterceptor {
	return GofraClientInterceptor
}

func GetServerInterceptor() grpc.UnaryServerInterceptor {
	return GofraServerInterceptor
}

// gofra client interceptor
var GofraClientInterceptor grpc.UnaryClientInterceptor = GofraClientInterceptorFunc

func GofraClientInterceptorFunc(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Debugf("====== Enter std client interceptor ======")
	defer log.Debugf("====== Leave std client interceptor ======")

	log.Debugf("context:%v", ctx)
	log.Debugf("method:%v", method)
	log.Debugf("options:%v", opts)
	log.Debugf("req:%v", req)

	//Monitor method enter total
	monitor.Increment(method + ",type=Client.Total")

	// Invoke remote
	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		log.Warnf("invoke failed!!! error:%v", err)

		//Monitor method fail total
		monitor.Increment(method + ",type=Client.Fail")
	} else {
		log.Debugf("reply:%v", reply)

		//Monitor method success total
		monitor.Increment(method + ",type=Client.Success")
	}

	return err
}

// gofra server interceptor
var GofraServerInterceptor grpc.UnaryServerInterceptor = GofraServerInterceptorFunc

func GofraServerInterceptorFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
	log.Debugf("====== Enter std server interceptor ======")
	defer log.Debugf("====== Leave std server interceptor ======")

	log.Debugf("context:%v", ctx)
	log.Debugf("method:%v", info.FullMethod)
	log.Debugf("server:%v", info.Server)
	log.Debugf("req:%v", req)

	//Monitor method enter total
	monitor.Increment(info.FullMethod + ",type=Server.Total")

	// Process
	reply, err = handler(ctx, req)

	if err != nil {
		log.Warnf("handle failed!!! error:%v", err)

		//Monitor method fail total
		monitor.Increment(info.FullMethod + ",type=Server.Fail")
	} else {
		log.Debugf("reply:%v", reply)

		//Monitor method success total
		monitor.Increment(info.FullMethod + ",type=Server.Success")
	}

	return reply, err
}
