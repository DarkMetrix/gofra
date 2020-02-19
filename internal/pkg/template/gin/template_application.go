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
	"time"
	"context"

	"os"
	"os/signal"

	"net/http"
	_ "net/http/pprof"

	log "github.com/cihub/seelog"

	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"

	//Component
	logger "github.com/DarkMetrix/gofra/pkg/logger/seelog"
	monitor "{{.MonitorPackage}}"
	tracing "{{.TracingPackage}}"
	performance "github.com/DarkMetrix/gofra/pkg/performance"

	//Gin relative
	"github.com/gin-gonic/gin"
	logMiddleware "github.com/DarkMetrix/gofra/pkg/gin-utils/middleware/log_middleware/seelog"
	monitorMiddleware "github.com/DarkMetrix/gofra/pkg/gin-utils/middleware/monitor_middleware/statsd"
	recoveryMiddleware "github.com/DarkMetrix/gofra/pkg/gin-utils/middleware/recovery_middleware/recovery"

	//gRPC relative
	logInterceptor "github.com/DarkMetrix/gofra/pkg/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "{{.MonitorInterceptorPackage}}"
	tracingInterceptor "{{.TracingInterceptorPackage}}"

	//Common
	pool "github.com/DarkMetrix/gofra/pkg/grpc-utils/pool"
	commonUtils "github.com/DarkMetrix/gofra/pkg/utils"

	"{{.WorkingPathRelative}}/internal/pkg/common"
	"{{.WorkingPathRelative}}/internal/pkg/config"

	//Http handler
	httpHandler "{{.WorkingPathRelative}}/internal/http_handler"
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
	err := logger.Init("../configs/log.config", common.ProjectName)

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
	httpClose, err := app.runHttpServer(address)

	if err != nil {
		log.Warnf("app.runHttpServer failed! error:%v", err.Error())
		return err
	}

	defer httpClose()

	// deal with signals, when interrupt was notified, server will stop gracefully
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	signalOccur := <- signalChannel

	log.Infof("Signal occured, signal:%v", signalOccur.String())

	return nil
}

type httpCloseFunc func()

func (app *Application) runHttpServer(address string) (httpCloseFunc, error) {
	//Set release or debug mode
	if config.GetConfig().Server.GinDebug == 0 {
		gin.SetMode(gin.ReleaseMode)
		log.Infof("gin runs in release mode!")
	} else {
		gin.SetMode(gin.DebugMode)
		log.Infof("gin runs in debug mode!")
	}

	// init engine
	engine := gin.Default()

	group := engine.Group("/",
		recoveryMiddleware.GetMiddleware(),
		logMiddleware.GetMiddleware(),
		monitorMiddleware.GetMiddleware())

	// add http handler
	//!!!DO NOT EDIT!!!
	group.POST("/health", httpHandler.HEALTH)
	/*@REGISTER_HTTP_STUB*/

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

	return func(){
		// stop http service gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		err := httpServer.Shutdown(ctx)

		if err != nil {
			log.Warnf("Http server stopped gracefully failed! error:%s", err.Error())
		} else {
			log.Infof("Http server stopped gracefully!")
		}
	}, nil
}
`
