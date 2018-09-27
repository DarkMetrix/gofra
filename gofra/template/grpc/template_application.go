package grpc

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
	"os"
	"os/signal"
	"net"

	"net/http"
	_ "net/http/pprof"

	log "github.com/cihub/seelog"

	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"

	logger "github.com/DarkMetrix/gofra/common/logger/seelog"
	monitor "{{.MonitorPackage}}"
	tracing "{{.TracingPackage}}"
	performance "github.com/DarkMetrix/gofra/common/performance"

	recoverInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/recover_interceptor"
	logInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "{{.MonitorInterceptorPackage}}"
	tracingInterceptor "{{.TracingInterceptorPackage}}"

	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"
	commonUtils "github.com/DarkMetrix/gofra/common/utils"

	"{{.WorkingPathRelative}}/src/common"
	"{{.WorkingPathRelative}}/src/config"

	//!!!DO NOT EDIT!!!
	/*@PROTO_STUB*/
	/*@HANDLER_STUB*/
)

var globalApplication *Application

type Application struct {
	ServerOpts []grpc.ServerOption
	ClientOpts []grpc.DialOption
}

//New Application
func newApplication() *Application {
	return &Application{}
}

//Get singleton application
func GetApplication() *Application {
	if globalApplication == nil {
		globalApplication = newApplication()
	}

	return globalApplication
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

	//Init pprof
	if conf.Pprof.Active != 0 {
		go func() {
			log.Infof("Begin pprof at addr:%v", conf.Pprof.Addr)
			http.ListenAndServe(conf.Pprof.Addr, nil)
		}()
	}

	// init monitor
	err = monitor.Init(config.GetConfig().Monitor.Params...)

	if err != nil {
		log.Warnf("Init monitor failed! error:%v", err.Error())
	}

	// init tracing
	err = tracing.Init(config.GetConfig().Tracing.Params...)

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
			monitorInterceptor.GetClientInterceptor())), grpc.WithInsecure())

	pool.GetConnectionPool().Init(app.ClientOpts)

	// init performance
	if conf.Performance.Active != 0 {
		switch conf.Performance.Type {
		case "log":
			go performance.BeginMemoryPerformanceMonitorWithLog()
			go performance.BeginGoroutinePerformanceMonitorWithLog()
		case "statsd":
			go performance.BeginMemoryPerformanceMonitorWithStatsd()
			go performance.BeginGoroutinePerformanceMonitorWithStatsd()
		default:
			log.Warnf("Performance type not found! Type:%v", conf.Performance.Type)
		}
	}

	return nil
}

//Run application
func (app *Application) Run(address string) error {
	defer log.Flush()
	defer tracing.Close()

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
	go func() {
		// deal with signals, when interrupt was notified, server will stop gracefully
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt)

		signalOccur := <- signalChannel

		log.Infof("Signal occured, signal:%v", signalOccur.String())

		// stop server gracefully
		s.GracefulStop()

		log.Infof("Server stopped gracefully!")
	}()

	err = s.Serve(listen)

	if err != nil {
		log.Errorf("Serve failed! error:%v", err.Error())
	} else {
		log.Infof("Serve quit!")
	}

	return nil
}
`
