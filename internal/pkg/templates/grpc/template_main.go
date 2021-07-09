package grpc

import (
	"golang.org/x/xerrors"

	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/templates"
)

// MainInfo represents the main file information
type MainInfo struct {
	Opts              *option.Options
	ConfigPackagePath string
}

// NewMainInfo returns a new MainInfo pointer
func NewMainInfo(opts ...option.Option) *MainInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &MainInfo{Opts: newOpts, ConfigPackagePath: newOpts.ConfigPackagePath}
}

// RenderFile render template and output to file
func (main *MainInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, main.Opts.Override, main.Opts.IgnoreExist,
		"template-main", MainTemplate, main); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

var MainTemplate string = `package main

import (
	"net"
	"os"
	"os/signal"

	recoverInterceptor "github.com/DarkMetrix/gofra/pkg/grpc-utils/interceptor/recover_interceptor"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
	pflag "github.com/spf13/pflag"
	viper "github.com/spf13/viper"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"

	config "{{.ConfigPackagePath}}"
    // !!!DO NOT EDIT!!!
    /*@PROTO_STUB*/
    /*@HANDLER_STUB*/
)

func main() {
	log.Info("====== Server [default] Start ======")

	// init config
	conf, err := initConfig()
	if err != nil {
		log.Fatalf("initConfig failed! error:%+v", err)
	}

	// run to serve grpc
	closeFunc, err := startGRPCServer(conf)
	if err != nil {
		log.Fatalf("runGRPCServer failed! error:%v", err)
	}
	defer closeFunc()

	// deal with signals, when interrupt was notified, server will stop gracefully
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	signalOccur := <- signalChannel
	log.Infof("Signal occurred, signal:%v", signalOccur.String())
}

func initConfig() (*config.Config, error) {
	// init flags
	pflag.String("config.path", "../configs/gofra.yaml", "Config file path, default '../configs/gofra.yaml'")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, xerrors.Errorf("viper.BindPFlags failed! error:%w", err)
	}

	// init config
	conf := config.GetConfig()
	if err := conf.Init(viper.GetString("config.path")); err != nil {
	    return nil, xerrors.Errorf("conf.Init failed! error:%w", err)
	}
	log.Infof("config:%v", conf)
	return conf, nil
}

func startGRPCServer(conf *config.Config) (func(), error) {
	// set server interceptor
	var serverInterceptors []grpc.UnaryServerInterceptor
	serverInterceptors = append(serverInterceptors, recoverInterceptor.GetServerInterceptor())

	serverOpts := make([]grpc.ServerOption, 0)
	serverOpts = append(serverOpts,
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(serverInterceptors...)))

	// listen
	listen, err := net.Listen("tcp", conf.Server.Addr)
	if err != nil {
		return nil, xerrors.Errorf("net.Listen failed! error:%w", err)
	}

	// init grpc server
	server := grpc.NewServer(serverOpts...)

	// register services
	// !!!DO NOT EDIT!!!
	/*@REGISTER_STUB*/

	// run to serve
	go func() {
		if err := server.Serve(listen); err != nil {
			log.Fatalf("server.Serve failed! error:%v", err.Error())
		}
		log.Infof("server quit!")
	}()

	return func() {
		// stop grpc service gracefully
		server.GracefulStop()
		log.Infof("gRPC server stopped gracefully!")
	}, nil
}
`
