package jaeger

import (

"errors"
"fmt"
"io"


	"github.com/uber/jaeger-client-go"

	"github.com/opentracing/opentracing-go"
)

var tracer opentracing.Tracer
var closer io.Closer

func Init(args... string) error {
	if len(args) < 1 {
		return errors.New(fmt.Sprintf("param invalid! args:%v", args))
	}

	addr := args[0]
	serviceName := args[1]

	err := InitJaeger(addr, serviceName)

	if err != nil {
		return err
	}

	opentracing.InitGlobalTracer(tracer)

	return nil
}

func Close() error {
	return closer.Close()
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
	tracer, _ = jaeger.NewTracer(serviceName, sampler, reporter)

	return nil
}

func GetTracer() opentracing.Tracer {
	return tracer
}
