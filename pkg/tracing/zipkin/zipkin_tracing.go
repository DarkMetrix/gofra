package zipkin

import (
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go"
	reporter "github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
)

var globalReporter reporter.Reporter

func Init(args... string) error {
	if len(args) < 3 {
		return errors.New(fmt.Sprintf("param invalid! args:%v", args))
	}

	addr := args[0]
	hostPort := args[1]
	serviceName := args[2]

	err := InitZipkin(addr, hostPort, serviceName)

	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	globalReporter.Close()
	return nil
}

func InitZipkin(addr string, hostPort string, serviceName string) error {
	// set up a span reporter
	globalReporter = zipkinhttp.NewReporter(addr)

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(serviceName, hostPort)

	if err != nil {
		return err
	}

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(globalReporter, zipkin.WithLocalEndpoint(endpoint))

	if err != nil {
		return err
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)

	// optionally set as Global OpenTracing tracer instance
	opentracing.SetGlobalTracer(tracer)

	return nil
}

func GetTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}
