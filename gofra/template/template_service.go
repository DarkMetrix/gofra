package template

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
}

var ServiceRpcTemplate string = `
/**********************************
 * Author : {{.Author}}
 * Time : {{.Time}}
 **********************************/

package {{.ServiceName}} 

import (
	context "golang.org/x/net/context"

	pb "{{.WorkingPathRelative}}/src/proto/{{.FileNamePrefix}}"
)

func (this {{.ServiceName}}Impl) {{.RpcName}} (ctx context.Context, req *pb.{{.Request}}) (*pb.{{.Response}}, error) {
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

// The health check service definition.
service HealthCheckService {
    // Sends a health check request
    rpc HealthCheck (HealthCheckRequest) returns (HealthCheckResponse) {}
}

message HealthCheckRequest {
    string message = 1;
}

message HealthCheckResponse {
    string message = 1;
}
`
