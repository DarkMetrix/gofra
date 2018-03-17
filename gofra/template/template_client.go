package template

type TestClientInfo struct {
	Author string
	Time string

	Project string

	Addr string

	WorkingPathRelative string

	MonitorPackage string
	MonitorInitParam string

	InterceptorPackage string
}

var TestClientTemplate string = `
package main

import (
	"fmt"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"

	health_check "{{.WorkingPathRelative}}/src/proto/health_check"

	interceptor "{{.InterceptorPackage}}"
	monitor "{{.MonitorPackage}}"
	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"
)

func main() {
	fmt.Println("====== Test [{{.Project}}] begin ======")
	defer fmt.Println("====== Test [{{.Project}}] end ======")

	// init statsd
	monitor.Init("{{.MonitorInitParam}}")

	// dial remote server
	pool.GetConnectionPool().Init(interceptor.GetClientInterceptor(), 5, 10, time.Second * 10)

	// RPC call
	req := new(health_check.HealthCheckRequest)
	req.Message = "ping"

	for index := 0; index < 1; index++ {
		// get remote server connection
		conn, err := pool.GetConnectionPool().GetConnection(context.Background(),"{{.Addr}}")
		defer conn.Close()

		// new client
		c := health_check.NewHealthCheckServiceClient(conn.ClientConn)

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

	time.Sleep(time.Second * 1)
}`
