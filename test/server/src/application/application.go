package application

import (
	"net"
	"time"

	"google.golang.org/grpc"

	logger "github.com/DarkMetrix/gofra/grpc-utils/logger/seelog"
	monitor "github.com/DarkMetrix/gofra/grpc-utils/monitor/statsd"
	interceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/gofra"
	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"

	"github.com/DarkMetrix/gofra/test/server/src/config"

	pb "github.com/DarkMetrix/gofra/test/proto"
	UserServiceHandler "github.com/DarkMetrix/gofra/test/server/src/handler/UserService"
)

type Application struct {
	ServerOpts []grpc.ServerOption
}

//Init application
func (app *Application) Init(conf *config.Config) error {
	// init log
	logger.Init("../conf/log.config")

	// init statsd
	monitor.Init("127.0.0.1:8125")

	// set server interceptor
	app.ServerOpts = append(app.ServerOpts, grpc.UnaryInterceptor(interceptor.GofraServerInterceptor))

	// set client interceptor
	err := pool.GetConnectionPool().Init(interceptor.GofraClientInterceptor,
		conf.Client.Pool.InitConns, conf.Client.Pool.MaxConns, time.Second * time.Duration(conf.Client.Pool.IdleTime))

	if err != nil {
		return err
	}

	return nil
}

//Run application
func (app *Application) Run(network, address string) error {
	listen, err := net.Listen(network, address)

	if err != nil {
		return err
	}

	// init grpc Server
	s := grpc.NewServer(app.ServerOpts ...)

	// register UserService
	pb.RegisterUserServiceServer(s, UserServiceHandler.UserServiceImpl{})

	// run to servce
	err = s.Serve(listen)

	if err != nil {
		return err
	}

	return nil
}