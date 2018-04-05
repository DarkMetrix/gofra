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
	"fmt"
	"time"
	"context"

    "google.golang.org/grpc"
	"google.golang.org/grpc/status"

    "github.com/grpc-ecosystem/go-grpc-middleware"

	health_check "{{.WorkingPathRelative}}/src/proto/health_check"

	logInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "{{.MonitorInterceptorPackage}}"
	tracingInterceptor "{{.TracingInterceptorPackage}}"

    logger "github.com/DarkMetrix/gofra/common/logger/seelog"
	monitor "{{.MonitorPackage}}"
    tracing "{{.TracingPackage}}"

	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"
	commonUtils "github.com/DarkMetrix/gofra/common/utils"
)

func main() {
	fmt.Println("====== Test [{{.Project}}] begin ======")
	defer fmt.Println("====== Test [{{.Project}}] end ======")

    // init log
    logger.Init("../conf/log.config", "{{.Project}}_test")

	// init monitor
	monitor.Init({{.MonitorInitParam}})

    // init tracing
    tracing.Init({{.TracingInitParam}})

	// dial remote server
	clientOpts := make([]grpc.DialOption, 0)

	clientOpts = append(clientOpts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			tracingInterceptor.GetClientInterceptor(),
			logInterceptor.GetClientInterceptor(),
			monitorInterceptor.GetClientInterceptor())))

	pool.GetConnectionPool().Init(clientOpts, 5, 10, time.Second * 10)

	addr := commonUtils.GetRealAddrByNetwork("{{.Addr}}")

	// begin test
	testHealthCheck(addr)

	time.Sleep(time.Second * 1)
}

func testHealthCheck(addr string) {
	// rpc call
	req := new(health_check.HealthCheckRequest)
	req.Message = "ping"

	for index := 0; index < 1; index++ {
		// get remote server connection
		conn, err := pool.GetConnectionPool().GetConnection(context.Background(), addr)
		defer conn.Recycle()

		// new client
		c := health_check.NewHealthCheckServiceClient(conn.Get())

		if err != nil {
			fmt.Printf("HealthCheck get connection failed! error%v", err.Error())
			continue
		}

		_, err = c.HealthCheck(context.Background(), req)

		if err != nil {
			stat, ok := status.FromError(err)

			if ok {
				fmt.Printf("HealthCheck request failed! code:%d, message:%v\r\n",
					stat.Code(), stat.Message())
			} else {
				fmt.Printf("HealthCheck request failed! err:%v\r\n", err.Error())
			}

			conn.Unhealhty()

			return
		}
	}
}
`
