package statsd_interceptor

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	monitor "github.com/DarkMetrix/gofra/grpc-utils/monitor/statsd"
)

func GetClientInterceptor() grpc.UnaryClientInterceptor {
	return GofraClientInterceptor
}

func GetServerInterceptor() grpc.UnaryServerInterceptor {
	return GofraServerInterceptor
}

// statsd client interceptor
var GofraClientInterceptor grpc.UnaryClientInterceptor = GofraClientInterceptorFunc

func GofraClientInterceptorFunc(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	//Monitor method enter total
	monitor.Increment(method + ",type=Client.Total")

	// Invoke remote
	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		//Monitor method fail total
		monitor.Increment(method + ",type=Client.Fail")
	} else {
		//Monitor method success total
		monitor.Increment(method + ",type=Client.Success")
	}

	return err
}

// statsd server interceptor
var GofraServerInterceptor grpc.UnaryServerInterceptor = GofraServerInterceptorFunc

func GofraServerInterceptorFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
	//Monitor method enter total
	monitor.Increment(info.FullMethod + ",type=Server.Total")

	// Process
	reply, err = handler(ctx, req)

	if err != nil {
		//Monitor method fail total
		monitor.Increment(info.FullMethod + ",type=Server.Fail")
	} else {
		//Monitor method success total
		monitor.Increment(info.FullMethod + ",type=Server.Success")
	}

	return reply, err
}
