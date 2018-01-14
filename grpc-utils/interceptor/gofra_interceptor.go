package interceptor

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/DarkMetrix/gofra/grpc-utils/monitor"
	log "github.com/cihub/seelog"
)

// gofra client interceptor
var GofraClientInterceptor grpc.UnaryClientInterceptor = GofraClientInterceptorFunc

func GofraClientInterceptorFunc(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Infof("====== Enter std client interceptor ======")
	log.Infof("context:%v", ctx)
	log.Infof("method:%v", method)
	log.Infof("options:%v", opts)
	log.Infof("req:%v", req)

	//Monitor method enter total
	monitor.GetStatsd().Increment(method + ",type=Client.Total")

	// Invoke remote
	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		log.Infof("invoke failed!!! error:%v", err)

		//Monitor method fail total
		monitor.GetStatsd().Increment(method + ",type=Client.Fail")
	} else {
		log.Infof("reply:%v", reply)

		//Monitor method success total
		monitor.GetStatsd().Increment(method + ",type=Client.Success")
	}

	log.Infof("====== Leave std client interceptor ======")

	return nil
}

// gofra server interceptor
var GofraServerInterceptor grpc.UnaryServerInterceptor = GofraServerInterceptorFunc

func GofraServerInterceptorFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
	log.Infof("====== Enter std server interceptor ======")
	log.Infof("context:%v", ctx)
	log.Infof("method:%v", info.FullMethod)
	log.Infof("server:%v", info.Server)
	log.Infof("req:%v", req)

	//Monitor method enter total
	monitor.GetStatsd().Increment(info.FullMethod + ",type=Server.Total")

	// Process
	reply, err = handler(ctx, req)

	if err != nil {
		log.Infof("handle failed!!! error:%v", err)

		//Monitor method fail total
		monitor.GetStatsd().Increment(info.FullMethod + ",type=Server.Fail")
	} else {
		log.Infof("reply:%v", reply)

		//Monitor method success total
		monitor.GetStatsd().Increment(info.FullMethod + ",type=Server.Success")
	}

	log.Infof("====== Leave std server interceptor ======")

	return reply, err
}
