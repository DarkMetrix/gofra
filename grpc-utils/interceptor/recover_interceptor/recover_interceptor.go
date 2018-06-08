package recover_interceptor

import (
	"runtime/debug"

	log "github.com/cihub/seelog"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

func GetServerInterceptor() grpc.UnaryServerInterceptor {
	recoverFunc := func (p interface{}) error {
		log.Errorf("Got panic! error:%v, stack:%v", p, string(debug.Stack()))
		return nil
	}

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoverFunc),
	}

	return grpc_recovery.UnaryServerInterceptor(opts...)
}
