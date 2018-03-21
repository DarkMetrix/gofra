package main

import (
	"fmt"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	pb "learn/test_grpc/basic/proto"

	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"

	log_interceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitor_interceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/statsd_interceptor"
	tracing_interceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/zipkin_interceptor"

	monitor "github.com/DarkMetrix/gofra/grpc-utils/monitor/statsd"
	tracing "github.com/DarkMetrix/gofra/grpc-utils/tracing/zipkin"
	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"
)

func main() {
    // init statsd
    monitor.InitStatsd("172.16.101.128:8125")

	// init tracing
	tracing.Init("127.0.0.1", "false", "0.0.0.0:0", "test_client")

	// dial remote server
	clientOpts := make([]grpc.DialOption, 0)

	clientOpts = append(clientOpts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			log_interceptor.GetClientInterceptor(),
			monitor_interceptor.GetClientInterceptor(),
			tracing_interceptor.GetClientInterceptor())))

	pool.GetConnectionPool().Init(clientOpts, 5, 10, time.Second * 10)

    // RPC call
    req := new(pb.AddUserRequest)
    req.Session = new(pb.Session)
    req.User = new(pb.User)

    req.Session.Seq = "12345678"
    req.User.Name = "techieliu"
    req.User.Sex = 1

    for index := 0; index < 10000; index++ {
        // get remote server connection
        conn, err := pool.GetConnectionPool().GetConnection(context.Background(),":58888")

        // new client
        c := pb.NewUserServiceClient(conn.ClientConn)

        if err != nil {
            fmt.Printf("Get connection failed! error%v", err.Error())
            continue
        }

        resp, err := c.AddUser(context.Background(), req)

        if err != nil {
            stat, ok := status.FromError(err)

            if ok {
                    fmt.Printf("AddUser failed! code:%d, message:%v\r\n",
                                    stat.Code(), stat.Message())
            } else {
                    fmt.Printf("AddUser failed! err:%v\r\n", err.Error())
            }

            conn.Unhealhty()

            return
        }

        conn.Close()

        fmt.Println(resp.String())

        time.Sleep(time.Second * 5)
    }

    time.Sleep(time.Second * 1)
}
