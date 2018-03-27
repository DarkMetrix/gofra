package application

import (
	"net"
	"time"
	"google.golang.org/grpc"
	logger "github.com/DarkMetrix/gofra/grpc-utils/logger/seelog"
	monitor "github.com/DarkMetrix/gofra/grpc-utils/monitor/statsd"
	tracing "github.com/DarkMetrix/gofra/grpc-utils/tracing/zipkin"

	recover_interceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/recover_interceptor"
	log_interceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitor_interceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/statsd_interceptor"
	tracing_interceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/zipkin_interceptor"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"

	"github.com/DarkMetrix/gofra/test/server/src/config"

	pb "github.com/DarkMetrix/gofra/test/proto"
	UserServiceHandler "github.com/DarkMetrix/gofra/test/server/src/handler/UserService"
	"github.com/DarkMetrix/gofra/test/generate/src/common"
)

type Application struct {
	ServerOpts []grpc.ServerOption
	ClientOpts []grpc.DialOption
}

//Init application
func (app *Application) Init(conf *config.Config) error {
	// init log
	logger.Init("../conf/log.config", common.ProjectName)

	// init statsd
	monitor.Init("127.0.0.1:8125")

	// init tracing
	tracing.Init("127.0.0.1:9411", "false", ":58888", "test_server")

	// set server interceptor
	app.ServerOpts = append(app.ServerOpts, grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			recover_interceptor.GetServerInterceptor(),
			tracing_interceptor.GetServerInterceptor(),
			log_interceptor.GetServerInterceptor(),
			monitor_interceptor.GetServerInterceptor())))

	// set client interceptor
	app.ClientOpts = append(app.ClientOpts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			tracing_interceptor.GetClientInterceptor(),
			log_interceptor.GetClientInterceptor(),
			monitor_interceptor.GetClientInterceptor())))

	err := pool.GetConnectionPool().Init(app.ClientOpts,
		conf.Client.Pool.InitConns,
		conf.Client.Pool.MaxConns,
		time.Second * time.Duration(conf.Client.Pool.IdleTime))

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