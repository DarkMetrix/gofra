package grpc

import (
	"golang.org/x/xerrors"

	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/templates"
)

// HealthCheckProtoInfo represents the health check protobuf file information
type HealthCheckProtoInfo struct {
	Opts *option.Options
}

// NewHealthCheckProtoInfo returns a new HealthCheckProtoInfo pointer
func NewHealthCheckProtoInfo(opts ...option.Option) *HealthCheckProtoInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &HealthCheckProtoInfo{Opts: newOpts}
}

// RenderFile render template and output to file
func (proto *HealthCheckProtoInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, proto.Opts.Override, proto.Opts.IgnoreExist,
		"template-health-check-proto", HealthCheckProtoTemplate, proto); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

var HealthCheckProtoTemplate string = `syntax = "proto3";

package common.health.check;

// the health check service definition.
service HealthCheck {
    // sends a health check request
    rpc HealthCheck (HealthCheckRequest) returns (HealthCheckResponse) {}
}

// health check request
message HealthCheckRequest {
    string message = 1;
}

// health check response
message HealthCheckResponse {
    string message = 1;
}
`
