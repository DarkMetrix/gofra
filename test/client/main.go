package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "learn/test_grpc/basic/proto"

	"github.com/DarkMetrix/gofra/grpc-utils/interceptor"
	"github.com/DarkMetrix/gofra/grpc-utils/monitor"
)

func main() {
	// init statsd
	monitor.InitStatsd("172.16.101.128:8125")

	// dial remote server
	conn, err := grpc.Dial(":58888", grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.GofraClientInterceptor))

	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	// new client
	c := pb.NewUserServiceClient(conn)

	// RPC call
	req := new(pb.AddUserRequest)
	req.Session = new(pb.Session)
	req.User = new(pb.User)

	req.Session.Seq = "12345678"
	req.User.Name = "techieliu"
	req.User.Sex = 1

	for index := 0; index < 10000; index++ {
		resp, err := c.AddUser(context.Background(), req)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(resp.String())

		time.Sleep(time.Microsecond * 10)
	}

	time.Sleep(time.Second * 1)
}
