package gin

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
	"time"
	"context"

	"net/http"
	_ "net/http/pprof"

	log "github.com/cihub/seelog"

	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"

	//Component
	logger "github.com/DarkMetrix/gofra/common/logger/seelog"
	monitor "{{.MonitorPackage}}"
	tracing "{{.TracingPackage}}"
	performance "github.com/DarkMetrix/gofra/common/performance"

	//Gin relative
	"github.com/gin-gonic/gin"
	logMiddleware "github.com/DarkMetrix/gofra/gin-utils/middleware/log_middleware/seelog"
	monitorMiddleware "github.com/DarkMetrix/gofra/gin-utils/middleware/monitor_middleware/statsd"

	//gRPC relative
	logInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "{{.MonitorInterceptorPackage}}"
	tracingInterceptor "{{.TracingInterceptorPackage}}"

	//Common
	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"
	commonUtils "github.com/DarkMetrix/gofra/common/utils"

	"{{.WorkingPathRelative}}/src/common"
	"{{.WorkingPathRelative}}/src/config"

	//Http handler
	httpHandler "{{.WorkingPathRelative}}/src/http_handler"
)

var globalApplication *Application

type Application struct {
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
	// process conf.Server.HttpAddr
	conf.Server.HttpAddr = commonUtils.GetRealAddrByNetwork(conf.Server.HttpAddr)

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

	// set client interceptor
	app.ClientOpts = append(app.ClientOpts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			tracingInterceptor.GetClientInterceptor(),
			logInterceptor.GetClientInterceptor(),
			monitorInterceptor.GetClientInterceptor())), grpc.WithInsecure())

	// init pool
	err = pool.GetConnectionPool().Init(app.ClientOpts)

	if err != nil {
		log.Warnf("Init pool failed! error:%v", err.Error())
		return err
	}

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

	// run to serve http
	engine := gin.Default()

	group := engine.Group("/", gin.Recovery(),
		logMiddleware.GetMiddleware(),
		monitorMiddleware.GetMiddleware())

	// add http handler
	//!!!DO NOT EDIT!!!
	/*@REGISTER_HTTP_STUB*/

	//gin.SetMode(gin.ReleaseMode)
	httpServer := &http.Server{
		Addr: address,
		Handler: engine,
	}

	go func() {
		err := httpServer.ListenAndServe()

		if err != nil {
			log.Errorf("Serve http failed! error:%v", err.Error())

			time.Sleep(time.Second)
			os.Exit(-2)
		} else {
			log.Infof("Serve http quit")
		}
	}()

	// deal with signals, when interrupt was notified, server will stop gracefully
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	signalOccur := <- signalChannel

	log.Infof("Signal occured, signal:%v", signalOccur.String())

	// stop http service gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err := httpServer.Shutdown(ctx)

	if err != nil {
		log.Warnf("http shutdown failed! error:%s", err.Error())
	}

	log.Infof("Server stopped gracefully!")

	return nil
}
`
