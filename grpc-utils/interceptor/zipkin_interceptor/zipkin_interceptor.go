package zipkin_interceptor

import (
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
)

func GetClientInterceptor() grpc.UnaryClientInterceptor {
	return grpc_opentracing.UnaryClientInterceptor()
}

func GetServerInterceptor() grpc.UnaryServerInterceptor {
	return grpc_opentracing.UnaryServerInterceptor()
}
