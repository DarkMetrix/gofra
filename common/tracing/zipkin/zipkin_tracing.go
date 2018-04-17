package zipkin

import (
	"fmt"
	"errors"
	"context"

	"google.golang.org/grpc/metadata"

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

func GetTracingId(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return "TracingIdNotFound"
	}

	data, ok := md["x-b3-traceid"]

	if !ok {
		return "TracingIdNotFound"
	}

	for _, value := range data {
		if len(data) != 0 {
			return value
		}
	}

	return "TracingIdNotFound"
}

func GetSpanId(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return "SpanIdNotFound"
	}

	data, ok := md["x-b3-spanid"]

	if !ok {
		return "SpanIdNotFound"
	}

	for _, value := range data {
		if len(data) != 0 {
			return value
		}
	}

	return "SpanIdNotFound"
}

func GetTracer() opentracing.Tracer {
	return tracer
}
