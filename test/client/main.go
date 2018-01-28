package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/status"

	pb "learn/test_grpc/basic/proto"

	"github.com/DarkMetrix/gofra/grpc-utils/interceptor"
	"github.com/DarkMetrix/gofra/grpc-utils/monitor"
	"github.com/DarkMetrix/gofra/grpc-utils/pool"
)

func main() {
	// init statsd
	monitor.InitStatsd("172.16.101.128:8125")

	// dial remote server
	pool.GetConnectionPool().Init(interceptor.GofraClientInterceptor, 5, 10, time.Second * 10)

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
