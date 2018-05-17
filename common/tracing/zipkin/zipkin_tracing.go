package zipkin

import (
	"fmt"
	"errors"

	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"

	commonUtils "github.com/DarkMetrix/gofra/common/utils"
)

var tracer opentracing.Tracer

func Init(args... string) error {
	if len(args) < 1 {
		return errors.New(fmt.Sprintf("param invalid! args:%v", args))
	}

	addr := args[0]
	debugStr := args[1]
	hostPort := args[2]
	serviceName := args[3]

	debug := false

	if debugStr == "true" {
		debug = true
	}

	err := InitZipkin(addr, debug, hostPort, serviceName)

	if err != nil {
		return err
	}

	opentracing.InitGlobalTracer(tracer)

	return nil
}

func InitZipkin(addr string, debug bool, hostPort string, serviceName string) error {
	// create collector.
	collector, err := zipkin.NewHTTPCollector(addr)

	if err != nil {
		return err
	}

	hostPort = commonUtils.GetRealAddrByNetwork(hostPort)

	// create recorder.
	recorder := zipkin.NewRecorder(collector, debug, hostPort, serviceName)

	// create tracer.
	tracer, err = zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)

	if err != nil {
		return err
	}

	return nil
}

func GetTracer() opentracing.Tracer {
	return tracer
}
