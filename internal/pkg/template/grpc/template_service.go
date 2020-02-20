package grpc

type ServiceInfo struct {
	Author string
	Time string

	ServiceName string
}

var ServiceTemplate string = `
/**********************************
 * Author : {{.Author}}
 * Time : {{.Time}}
 **********************************/

package {{.ServiceName}}

// Implement {{.ServiceName}} interface
type {{.ServiceName}}Impl struct{}
`

type RpcInfo struct {
	Author string
	Time string
	WorkingPathRelative string

	ServiceName string
	FileNamePrefix string
	RpcName string
	Request string
	Response string

	MonitorPackage string
	TracingPackage string
}

var ServiceRpcTemplate string = `
/**********************************
 * Author : {{.Author}}
 * Time : {{.Time}}
 **********************************/

package {{.ServiceName}} 

import (
	"context"

	// log package
	// log "github.com/cihub/seelog"

	// monitor package
	// monitor "{{.MonitorPackage}}"

	// tracing package
	// tracing "{{.TracingPackage}}"

	pb "{{.WorkingPathRelative}}/api/protobuf_spec/{{.FileNamePrefix}}"
)

func (this {{.ServiceName}}Impl) {{.RpcName}} (ctx context.Context, req *pb.{{.Request}}) (*pb.{{.Response}}, error) {
	// log Example:traceid must be logged
	// log.Infof("{{.RpcName}} begin, traceid=%v, req=%v", tracing.GetTracingId(ctx), req)

	resp := new(pb.{{.Response}})

	return resp, nil
}
`

type HealthCheckServiceProtoInfo struct {
	Author string
	Time string
}

var HealthCheckServiceProtoTemplate string = `
// Author : {{.Author}}
// Time : {{.Time}}

syntax = "proto3";

package common.health.check;

// the health check service definition.
service HealthCheckService {
    // sends a health check request
    rpc HealthCheck (HealthCheckRequest) returns (HealthCheckResponse) {}
}

message HealthCheckRequest {
    string message = 1;
}

message HealthCheckResponse {
    string message = 1;
}
`
