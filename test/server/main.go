package main

import (
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "learn/test_grpc/basic/proto"

	"github.com/DarkMetrix/gofra/grpc-utils/interceptor"
	"github.com/DarkMetrix/gofra/grpc-utils/monitor"
)

// Implement UserService interface
type UserService struct{}

// UserService ...
var userService = UserService{}

func (h UserService) AddUser (ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	resp := new(pb.AddUserResponse)
	resp.Session = new(pb.Session)
	resp.Result = new(pb.Result)

	//Set session
	resp.Session = req.Session

	//Set result
	resp.GetResult().Code = 0
	resp.GetResult().Message = "Success"

	return resp, nil
}

func main() {
	// init statsd
	monitor.InitStatsd("172.16.101.128:8125")

	// start
	fmt.Println("====== Test grpc server ======")

	listen, err := net.Listen("tcp", ":58888")
	if err != nil {
		fmt.Println("failed to listen: ", err.Error())
		return
	}

	// set interceptor
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(interceptor.GofraServerInterceptor))

	// init grpc Server
	s := grpc.NewServer(opts...)

	// register UserService
	pb.RegisterUserServiceServer(s, userService)

	fmt.Println("Listen on :58888")

	s.Serve(listen)
}
