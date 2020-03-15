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
	"time"
	"os"
	"os/signal"
	"net"

	"net/http"
	_ "net/http/pprof"

	log "github.com/cihub/seelog"
	viper "github.com/spf13/viper"

	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"

	logger "github.com/DarkMetrix/gofra/pkg/logger/seelog"
	monitor "{{.MonitorPackage}}"
	tracing "{{.TracingPackage}}"
	performance "github.com/DarkMetrix/gofra/pkg/performance"

	recoverInterceptor "github.com/DarkMetrix/gofra/pkg/grpc-utils/interceptor/recover_interceptor"
	logInterceptor "github.com/DarkMetrix/gofra/pkg/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "{{.MonitorInterceptorPackage}}"
	tracingInterceptor "{{.TracingInterceptorPackage}}"

	pool "github.com/DarkMetrix/gofra/pkg/grpc-utils/pool"
	commonUtils "github.com/DarkMetrix/gofra/pkg/utils"

	"{{.WorkingPathRelative}}/internal/pkg/common"
	"{{.WorkingPathRelative}}/internal/pkg/config"

	// !!!DO NOT EDIT!!!
	/*@PROTO_STUB*/
	/*@HANDLER_STUB*/
)

var globalApplication *Application

type Application struct {
	ServerOpts []grpc.ServerOption
	ClientOpts []grpc.DialOption
}

// new Application
func newApplication() *Application {
	return &Application{}
}

// get singleton application
func GetApplication() *Application {
	if globalApplication == nil {
		globalApplication = newApplication()
	}

	return globalApplication
}

// init application
func (app *Application) Init(conf *config.Config) error {
	// process conf.Server.Addr
	conf.Server.Addr = commonUtils.GetRealAddrByNetwork(conf.Server.Addr)

	// init log
	err := logger.Init(viper.GetString("log.config.path"), common.ProjectName)

	if err != nil {
		log.Warnf("Init logger failed! error:%v", err.Error())
	}

	// init pprof
	if conf.Pprof.Active != 0 {
		go func() {
			log.Infof("Begin pprof at addr:%v", conf.Pprof.Addr)
			err = http.ListenAndServe(conf.Pprof.Addr, nil)

			if err != nil {
				log.Warnf("Pprof http.ListenAndServe failed! error:%v", err.Error())
			}
		}()
	}

	// init monitor
	if conf.Monitor.Active != 0 {
    	err = monitor.Init(conf.Monitor.Params...)
    
    	if err != nil {
    		log.Warnf("Init monitor failed! error:%v", err.Error())
    	}
    }

	// init tracing
	if conf.Tracing.Active != 0 {
    	err = tracing.Init(conf.Tracing.Params...)
    
    	if err != nil {
    		log.Warnf("Init tracing failed! error:%v", err.Error())
    	}
    }

	// set server interceptor
	var serverInterceptors []grpc.UnaryServerInterceptor
	serverInterceptors = append(serverInterceptors, recoverInterceptor.GetServerInterceptor())
	serverInterceptors = append(serverInterceptors, logInterceptor.GetServerInterceptor())

	if conf.Monitor.Active != 0 {
		serverInterceptors = append(serverInterceptors, monitorInterceptor.GetServerInterceptor())
	}

	if conf.Tracing.Active != 0 {
		serverInterceptors = append(serverInterceptors, tracingInterceptor.GetServerInterceptor())
	}

	app.ServerOpts = append(app.ServerOpts, 
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(serverInterceptors...)))

	// set client interceptor
    var clientInterceptors []grpc.UnaryClientInterceptor
	clientInterceptors = append(clientInterceptors, logInterceptor.GetClientInterceptor())
	
	if conf.Monitor.Active != 0 {
		clientInterceptors = append(clientInterceptors, monitorInterceptor.GetClientInterceptor())
	}
	
	if conf.Tracing.Active != 0 {
		clientInterceptors = append(clientInterceptors, tracingInterceptor.GetClientInterceptor())
	}

	app.ClientOpts = append(app.ClientOpts, 
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(clientInterceptors...)), 
		grpc.WithInsecure())

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

// run application
func (app *Application) Run(address string) error {
	defer log.Flush()
	defer tracing.Close()
	defer pool.GetConnectionPool().Close()

	// run to serve grpc
	grpcClose, err := app.runGRPCServer(address)

	if err != nil {
		log.Warnf("app.runGRPCServer failed! error:%v", err.Error())
		return err
	}

	defer grpcClose()

	// deal with signals, when interrupt was notified, server will stop gracefully
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	signalOccur := <- signalChannel

	log.Infof("Signal occurred, signal:%v", signalOccur.String())

	return nil
}

type grpcCloseFunc func()

func (app *Application) runGRPCServer(address string) (grpcCloseFunc, error) {
	// listen
	listen, err := net.Listen("tcp", address)

	if err != nil {
		return nil, err
	}

	// init grpc server
	s := grpc.NewServer(app.ServerOpts ...)

	// register services
	// !!!DO NOT EDIT!!!
	health_check.RegisterHealthCheckServiceServer(s, HealthCheckServiceHandler.HealthCheckServiceImpl{})
	/*@REGISTER_STUB*/

	// run to serve
	go func() {
		err = s.Serve(listen)

		if err != nil {
			log.Errorf("Serve gRPC failed! error:%v", err.Error())

			time.Sleep(time.Second)
			os.Exit(-2)
		} else {
			log.Infof("Serve gRPC quit!")
		}
	}()

	return func() {
		// stop grpc service gracefully
		s.GracefulStop()

		log.Infof("gRPC server stopped gracefully!")
	}, nil
}
`
