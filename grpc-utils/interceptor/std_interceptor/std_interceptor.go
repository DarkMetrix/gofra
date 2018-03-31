package std_interceptor

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func GetClientInterceptor() grpc.UnaryClientInterceptor {
	return StdClientInterceptor
}

func GetServerInterceptor() grpc.UnaryServerInterceptor {
	return StdServerInterceptor
}

// std client interceptor
var StdClientInterceptor grpc.UnaryClientInterceptor = StdClientInterceptorFunc

func StdClientInterceptorFunc(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Invoke remote
	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		fmt.Printf("context=%v, req=%v, invoke failed!!! error:%v\r\n", ctx, req, err)
	} else {
		fmt.Printf("context=%v, req=%v, invoke success!!! reply:%v\r\n", ctx, req, reply)
	}

	return err
}

// std server interceptor
var StdServerInterceptor grpc.UnaryServerInterceptor = StdServerInterceptorFunc

func StdServerInterceptorFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
	// Process
	reply, err = handler(ctx, req)

	if err != nil {
		fmt.Printf("context=%v, req=%v, invoke failed!!! error:%v\r\n", ctx, req, err)
	} else {
		fmt.Printf("context=%v, req=%v, invoke success!!! reply:%v\r\n", ctx, req, reply)
	}

	return reply, err
}
