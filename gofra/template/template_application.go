package template

type ApplicationInfo struct {
	Author string
	Time string
	Project string

	WorkingPathRelative string

	MonitorPackage string
	MonitorInitParam string

    TracingPackage string
    TracingInitParam string

	MonitorInterceptorPackage string
	TracingInterceptorPackage string
}

var ApplicationTemplate string = `
/**********************************
 * Author : {{.Author}}
 * Time : {{.Time}}
 **********************************/

package application

import (
	"net"
	"time"

	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"

	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"
	logger "github.com/DarkMetrix/gofra/grpc-utils/logger/seelog"
	monitor "{{.MonitorPackage}}"
	tracing "{{.TracingPackage}}"

	recover_interceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/recover_interceptor"
	log_interceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitor_interceptor "{{.MonitorInterceptorPackage}}"
	tracing_interceptor "{{.TracingInterceptorPackage}}"

	helper "github.com/DarkMetrix/gofra/grpc-utils/helper"

	"{{.WorkingPathRelative}}/src/common"
	"{{.WorkingPathRelative}}/src/config"

	//!!!DO NOT EDIT!!!
	/*@PROTO_STUB*/
	/*@HANDLER_STUB*/
)

type Application struct {
	ServerOpts []grpc.ServerOption
	ClientOpts []grpc.DialOption
}

//Init application
func (app *Application) Init(conf *config.Config) error {
	//Process conf.Server.Addr
	conf.Server.Addr = helper.GetRealAddrByNetwork(conf.Server.Addr)

	// init log
	logger.Init("../conf/log.config", common.ProjectName)

	// init monitor
	monitor.Init({{.MonitorInitParam}}, common.ProjectName)

	// init tracing
	tracing.Init({{.TracingInitParam}})

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
func (app *Application) Run(address string) error {
	listen, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	// init grpc server
	s := grpc.NewServer(app.ServerOpts ...)

	// register services
	//!!!DO NOT EDIT!!!
	/*@REGISTER_STUB*/

	// run to serve
	err = s.Serve(listen)

	if err != nil {
		return err
	}

	return nil
}
`
