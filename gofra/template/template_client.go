package template

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

	logInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "{{.MonitorInterceptorPackage}}"
	tracingInterceptor "{{.TracingInterceptorPackage}}"

    logger "github.com/DarkMetrix/gofra/common/logger/seelog"
	monitor "{{.MonitorPackage}}"
    tracing "{{.TracingPackage}}"

	health_check "{{.WorkingPathRelative}}/src/proto/health_check"
)

func main() {
	defer log.Flush()

    // init log
    err := logger.Init("../conf/log.config", "{{.Project}}_test")

	if err != nil {
		log.Warnf("Init logger failed! error:%v", err.Error())
	}

	log.Info("====== Test [{{.Project}}] begin ======")
	defer log.Info("====== Test [{{.Project}}] end ======")

	// init monitor
	err = monitor.Init({{.MonitorInitParam}})

	if err != nil {
		log.Warnf("Init monitor failed! error:%v", err.Error())
	}

    // init tracing
    err = tracing.Init({{.TracingInitParam}})

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

	// init conn
	conn, err := grpc.Dial("{{.Addr}}", clientOpts...)

	if err != nil {
		log.Warnf("grpc.Dial failed! error:%v", err.Error())
		return
	}

	// begin test
	testHealthCheck(conn)

	time.Sleep(time.Second * 1)
}

func testHealthCheck(conn *grpc.ClientConn) {
	// rpc call
	req := new(health_check.HealthCheckRequest)
	req.Message = "ping"

	for index := 0; index < 1; index++ {
		c := health_check.NewHealthCheckServiceClient(conn)

		_, err := c.HealthCheck(context.Background(), req)

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
