package recover_interceptor

import (
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

func GetServerInterceptor() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor()
}
