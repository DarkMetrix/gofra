package null_interceptor

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func GetClientInterceptor() grpc.UnaryClientInterceptor {
	return NullClientInterceptor
}

func GetServerInterceptor() grpc.UnaryServerInterceptor {
	return NullServerInterceptor
}

// null client interceptor
var NullClientInterceptor grpc.UnaryClientInterceptor = NullClientInterceptorFunc

func NullClientInterceptorFunc(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	// Invoke remote
	return invoker(ctx, method, req, reply, cc, opts...)
}

// null server interceptor
var NullServerInterceptor grpc.UnaryServerInterceptor = NullServerInterceptorFunc

func NullServerInterceptorFunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply interface{}, err error) {
	// Process
	return handler(ctx, req)
}
