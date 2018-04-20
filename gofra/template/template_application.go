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

	log "github.com/cihub/seelog"

	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"

	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"
	naming "github.com/DarkMetrix/gofra/common/naming"

	logger "github.com/DarkMetrix/gofra/common/logger/seelog"
	monitor "{{.MonitorPackage}}"
	tracing "{{.TracingPackage}}"

	recoverInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/recover_interceptor"
	logInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "{{.MonitorInterceptorPackage}}"
	tracingInterceptor "{{.TracingInterceptorPackage}}"

	commonUtils "github.com/DarkMetrix/gofra/common/utils"

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
	// process conf.Server.Addr
	conf.Server.Addr = commonUtils.GetRealAddrByNetwork(conf.Server.Addr)

	// init log
	err := logger.Init("../conf/log.config", common.ProjectName)

	if err != nil {
		log.Warnf("Init logger failed! error:%v", err.Error())
	}

	// init monitor
	monitorParams := append(config.GetConfig().Monitor.Params, common.ProjectName)
	err = monitor.Init(monitorParams...)

	if err != nil {
		log.Warnf("Init monitor failed! error:%v", err.Error())
	}

	// init tracing
	tracingParams := append(config.GetConfig().Tracing.Params, common.ProjectName)
	err = tracing.Init(tracingParams...)

	if err != nil {
		log.Warnf("Init tracing failed! error:%v", err.Error())
	}

	// set server interceptor
	app.ServerOpts = append(app.ServerOpts, grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			recoverInterceptor.GetServerInterceptor(),
			tracingInterceptor.GetServerInterceptor(),
			logInterceptor.GetServerInterceptor(),
			monitorInterceptor.GetServerInterceptor())))

	// set client interceptor
	app.ClientOpts = append(app.ClientOpts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			tracingInterceptor.GetClientInterceptor(),
			logInterceptor.GetClientInterceptor(),
			monitorInterceptor.GetClientInterceptor())))

	// init pool
	err = pool.GetConnectionPool().Init(app.ClientOpts,
		conf.Client.Pool.InitConns,
		conf.Client.Pool.MaxConns,
		time.Second * time.Duration(conf.Client.Pool.IdleTime))

	if err != nil {
		log.Warnf("Init pool failed! error:%v", err.Error())
		return err
	}

	// init naming
	err = naming.Init("../conf/naming.json")

	if err != nil {
		log.Warnf("Init naming failed! error:%v", err.Error())
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
