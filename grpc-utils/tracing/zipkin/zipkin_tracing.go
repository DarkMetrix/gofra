package zipkin

import (
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

var tracer opentracing.Tracer

func Init(args... string) {
	if len(args) < 1 {
		panic("Init args length < 1")
	}

	addr := args[0]
	debugStr := args[1]
	hostPort := args[2]
	serviceName := args[3]

	debug := false

	if debugStr == "true" {
		debug = true
	}

	InitZipkin(addr, debug, hostPort, serviceName)

	opentracing.InitGlobalTracer(tracer)
}

func InitZipkin(addr string, debug bool, hostPort string, serviceName string) {
	// create collector.
	collector, err := zipkin.NewHTTPCollector(addr)

	if err != nil {
		panic(err)
	}

	// create recorder.
	recorder := zipkin.NewRecorder(collector, debug, hostPort, serviceName)

	// create tracer.
	tracer, err = zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)

	if err != nil {
		panic(err)
	}
}
