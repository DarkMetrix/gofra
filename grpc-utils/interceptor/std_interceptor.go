package interceptor

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// std client interceptor
var StdClientInterceptor grpc.UnaryClientInterceptor = StdClientInterceptorFunc

func StdClientInterceptorFunc(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("====== Enter std client interceptor ======")
	fmt.Printf("context:%v\r\n", ctx)
	fmt.Printf("method:%v\r\n", method)
	fmt.Printf("options:%v\r\n", opts)
	fmt.Printf("req:%v\r\n", req)

	// Invoke remote
	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		fmt.Printf("invoke failed!!! error:%v\r\n", err)
	} else {
		fmt.Printf("reply:%v\r\n", reply)
	}

	fmt.Println("====== Leave std client interceptor ======")

	return err
}

// std server interceptor
var StdServerInterceptor grpc.UnaryServerInterceptor = StdServerInterceptorFunc

func StdServerInterceptorFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
	fmt.Println("====== Enter std server interceptor ======")
	fmt.Printf("context:%v\r\n", ctx)
	fmt.Printf("method:%v\r\n", info.FullMethod)
	fmt.Printf("server:%v\r\n", info.Server)
	fmt.Printf("req:%v\r\n", req)

	// Process
	reply, err = handler(ctx, req)

	if err != nil {
		fmt.Printf("handle failed!!! error:%v\r\n", err)
	} else {
		fmt.Printf("reply:%v\r\n", reply)
	}

	fmt.Println("====== Leave std server interceptor ======")

	return reply, err
}
