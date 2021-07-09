package grpc

import (
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/iancoleman/strcase"
	"golang.org/x/xerrors"

	"github.com/DarkMetrix/gofra/internal/pkg/templates"
)

// ServiceInfo represents the gRPC service information
type ServiceInfo struct {
	Opts        *option.Options
	ServiceName string
}

// NewServiceInfo returns a new ServiceInfo pointer
func NewServiceInfo(opts ...option.Option) *ServiceInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &ServiceInfo{
		Opts:        newOpts,
		ServiceName: newOpts.ServiceName,
	}
}

// RenderFile render template and output to file
func (service *ServiceInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, service.Opts.Override, service.Opts.IgnoreExist,
		"template-service", ServiceTemplate, service); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

var ServiceTemplate string = `package {{.ServiceName}}

// Implementation implements {{.ServiceName}} interface
type Implementation struct{}
`

// RPCInfo represents the gRPC RPC information
type RPCInfo struct {
	Opts                *option.Options
	PackageName         string
	ServiceName         string
	ImportedPackageName string
	RPCName             string
	Request             string
	Response            string
}

// NewRPCInfo returns a new RPCInfo pointer
func NewRPCInfo(opts ...option.Option) *RPCInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &RPCInfo{
		Opts:                newOpts,
		PackageName:         newOpts.PackageName,
		ServiceName:         strcase.ToCamel(newOpts.ServiceName),
		ImportedPackageName: newOpts.ImportedPackageName,
		RPCName:             strcase.ToCamel(newOpts.RPCName),
		Request:             strcase.ToCamel(newOpts.RequestName),
		Response:            strcase.ToCamel(newOpts.ResponseName),
	}
}

// RenderFile render template and output to file
func (rpc *RPCInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, rpc.Opts.Override, rpc.Opts.IgnoreExist,
		"template-rpc", RPCTemplate, rpc); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

var RPCTemplate string = `package {{.PackageName}}

import (
    "context"

    pb "{{.ImportedPackageName}}"
)

// {{.RPCName}} implements {{.ServiceName}} interface 
func (service Implementation) {{.RPCName}} (ctx context.Context, req *pb.{{.Request}}) (*pb.{{.Response}}, error) {
    resp := &pb.{{.Response}}{}

    // TODO: implementation
    return resp, nil
}
`
