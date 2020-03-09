package jaeger

import (
	"fmt"
	"errors"
	"io"

	"github.com/uber/jaeger-client-go"

	"github.com/opentracing/opentracing-go"
)

var globalCloser io.Closer = nil

func Init(args... string) error {
	if len(args) < 2 {
		return errors.New(fmt.Sprintf("param invalid! args:%v", args))
	}

	addr := args[0]
	serviceName := args[1]

	err := InitJaeger(addr, serviceName)

	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	if globalCloser != nil {
		return globalCloser.Close()
	}

	return nil
}

func InitJaeger(addr string,  serviceName string) error {
	// create sampler
	sampler := jaeger.NewConstSampler(true)

	// create transport
	transport, err := jaeger.NewUDPTransport(addr, 64 * 1024)

	if err != nil {
		return err
	}

	// create report
	reporter := jaeger.NewRemoteReporter(transport)

	// new tracer
	tracer, closer := jaeger.NewTracer(serviceName, sampler, reporter)
	globalCloser = closer

	opentracing.SetGlobalTracer(tracer)

	return nil
}

func GetTracer() opentracing.Tracer {
	return opentracing.GlobalTracer()
}
