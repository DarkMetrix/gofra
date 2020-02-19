package grpc

type TestClientInfo struct {
	Author string
	Time string

	Project string

	Addr string

	WorkingPathRelative string

	MonitorPackage string
	MonitorInitParam string

    TracingPackage string
    TracingInitParam string

	MonitorInterceptorPackage string
	TracingInterceptorPackage string
}

var TestClientTemplate string = `
package main

import (
	"time"
	"context"

	log "github.com/cihub/seelog"

    "google.golang.org/grpc"
	"google.golang.org/grpc/status"

    "github.com/grpc-ecosystem/go-grpc-middleware"

	logInterceptor "github.com/DarkMetrix/gofra/pkg/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "{{.MonitorInterceptorPackage}}"
	tracingInterceptor "{{.TracingInterceptorPackage}}"

    logger "github.com/DarkMetrix/gofra/pkg/logger/seelog"
	monitor "{{.MonitorPackage}}"
    tracing "{{.TracingPackage}}"
    pool "github.com/DarkMetrix/gofra/pkg/grpc-utils/pool"

	health_check "{{.WorkingPathRelative}}/api/protobuf_spec/health_check"
)

func main() {
	defer log.Flush()

    // init log
    err := logger.Init("../configs/log.config", "{{.Project}}_test")

	if err != nil {
		log.Warnf("Init logger failed! error:%v", err.Error())
	}

	log.Info("====== Test [{{.Project}}] begin ======")
	defer log.Info("====== Test [{{.Project}}] end ======")

	// init monitor
	err = monitor.Init("127.0.0.1:8125", "{{.Project}}_test")

	if err != nil {
		log.Warnf("Init monitor failed! error:%v", err.Error())
	}

    // init tracing
    err = tracing.Init("127.0.0.1:6831", "{{.Project}}_test")

	if err != nil {
		log.Warnf("Init tracing failed! error:%v", err.Error())
	}

	// dial remote server
	clientOpts := make([]grpc.DialOption, 0)

	clientOpts = append(clientOpts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			tracingInterceptor.GetClientInterceptor(),
			logInterceptor.GetClientInterceptor(),
			monitorInterceptor.GetClientInterceptor())), grpc.WithInsecure())

	// init grpc connection pool
	err = pool.GetConnectionPool().Init(clientOpts)

	if err != nil {
		log.Warnf("Init grpc pool failed! error:%v", err.Error())
		return
	}

	// begin test
	testHealthCheck()

	time.Sleep(time.Second * 1)
}

func testHealthCheck() {
	// rpc call
	req := new(health_check.HealthCheckRequest)
	req.Message = "ping"

	for index := 0; index < 1; index++ {
		// get connection
		conn, err := pool.GetConnectionPool().GetConnection(context.Background(), "{{.Addr}}")

		if err != nil {
			log.Warnf("pool.GetConnection failed! error:%v", err.Error())
			continue
		}

		// call
		c := health_check.NewHealthCheckServiceClient(conn)

		_, err = c.HealthCheck(context.Background(), req)

		if err != nil {
			stat, ok := status.FromError(err)

			if ok {
				log.Warnf("HealthCheck request failed! code:%d, message:%v",
					stat.Code(), stat.Message())
			} else {
				log.Warnf("HealthCheck request failed! err:%v", err.Error())
			}

			return
		}
	}
}
`
