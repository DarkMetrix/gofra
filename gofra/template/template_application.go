package template

type ApplicationInfo struct {
	Author string
	Time string

	WorkingPathRelative string

	MonitorPackage string
	MonitorInitParam string

	InterceptorPackage string
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

	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"
	logger "github.com/DarkMetrix/gofra/grpc-utils/logger/seelog"
	monitor "{{.MonitorPackage}}"
	interceptor "{{.InterceptorPackage}}"

	"{{.WorkingPathRelative}}/src/config"

	//!!!DO NOT EDIT!!!
	/*@PROTO_STUB*/
	/*@HANDLER_STUB*/
)

type Application struct {
	ServerOpts []grpc.ServerOption
}

//Init application
func (app *Application) Init(conf *config.Config) error {
	// init log
	logger.Init("../conf/log.config")

	// init monitor
	monitor.Init("{{.MonitorInitParam}}")

	// set server interceptor
	app.ServerOpts = append(app.ServerOpts, grpc.UnaryInterceptor(interceptor.GetServerInterceptor()))

	// set client interceptor
	err := pool.GetConnectionPool().Init(interceptor.GetClientInterceptor(),
		conf.Client.Pool.InitConns, conf.Client.Pool.MaxConns, time.Second * time.Duration(conf.Client.Pool.IdleTime))

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