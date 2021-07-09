package directory

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/pkg/utils"
	"github.com/iancoleman/strcase"
	"golang.org/x/xerrors"
)

/*********************************************************************************
gRPC layout follows the rules of https://github.com/golang-standards/project-layout

The layout looks like below

.

 *********************************************************************************/

// GRPCServiceLayout defines the layout interface of a gRPC service without
type GRPCServiceLayout interface {
	ServiceLayout

	// GetMainFilePath returns the main.go file path, e.g.: /.../cmd/main.go
	GetMainFilePath() string
	// GetAPIProtobufBasePath returns the API protobuf root directory path, e.g.: /.../api/protobuf_spec
	GetAPIProtobufBasePath() string
	// GetAPIProtobufPath returns the specific API directory path, e.g.: /.../api/protobuf_spec/health_check
	GetAPIProtobufPath(protoPath string) string
	// GetAPIProtobufFilePath returns the specific API file path, e.g.: /.../api/protobuf_spec/health_check/health_check.proto
	GetAPIProtobufFilePath(protoFile string) string

	// GetGRPCServiceBasePath returns the gRPC gateway directory path, e.g.: /.../internal/service/grpc
	GetGRPCServiceBasePath() string
	// GetGRPCServicePath returns the specific gRPC gateway directory path, e.g.: /.../internal/service/grpc/health_check
	GetGRPCServicePath(service string) string
	// GetGRPCServiceFilePath returns the specific gRPC gateway service file path, e.g.: /.../internal/service/grpc/health_check/health_check.go
	GetGRPCServiceFilePath(service string) string
	// GetGRPCRPCFilePath returns the specific gRPC gateway RPC file path, e.g.: /.../internal/service/grpc/health_check/check.go
	GetGRPCRPCFilePath(service, RPC string) string

	// GetConfigGoBasePath returns the config path, e.g.: /.../internal/config
	GetConfigGoBasePath() string
	// GetConfigGoFilePath returns the config.go file path, e.g.: /.../internal/config/config.go
	GetConfigGoFilePath() string
	// GetConfigYAMLFilePath returns the gofra.yaml file path, e.g.: /.../configs/gofra.yaml
	GetConfigYAMLFilePath() string
	// GetConfigPackageName returns the package name of config, e.g.: github.com/foo/bar/internal/config
	GetConfigPackageName(goModule string) string

	// GetPackageAlias returns the alias of specific package name
	GetPackageAlias(packageName string) string
	// GetServicePackagePath returns the package path of specific service
	GetServicePackagePath(goModule, service string) string
	// GetProtoPackagePath returns the package path of specific proto file
	GetProtoPackagePath(goModule, protoFile string) string
	// GetServiceRegisterStub returns the stub of specific service of proto
	GetServiceRegisterStub(protoFile, serviceName string) string
	// GetServiceImportStub returns the service import line of code
	GetServiceImportStub(goModule, service string) string
	// GetProtoImportStub returns the proto import line of code
	GetProtoImportStub(goModule, protoPath string) string
}

// GRPCLayout implements Layout interface to generate gRPC service directory structure
type GRPCLayout struct {
	BasicServiceLayout
}

// NewGRPCLayout returns the gRPC directory layout pointer
func NewGRPCLayout(option ...option.Option) *GRPCLayout {
	basicServiceLayout := NewBasicServiceLayout(option...)
	return &GRPCLayout{*basicServiceLayout}
}

// String returns the JSON format of directory information
func (layout *GRPCLayout) String() string {
	jsonBuf, err := json.Marshal(layout)
	if err != nil {
		return ""
	}
	return string(jsonBuf)
}

// Save creates all the generated directory structure
func (layout *GRPCLayout) Save() error {
	// initialize basic service layout
	if err := layout.BasicServiceLayout.Save(); err != nil {
		return xerrors.Errorf("layout.BasicServiceLayout.Save failed! error:%w", err)
	}

	// create proto base & service base directories
	if err := utils.CreatePaths(
		layout.Options.Override,
		layout.GetAPIProtobufBasePath(),
		layout.GetGRPCServiceBasePath(),
		layout.GetConfigGoBasePath()); err != nil {
		return xerrors.Errorf("utils.CreatePaths failed! "+
			"api base:%v, service base:%v, error:%w", layout.GetAPIProtobufBasePath(), layout.GetServiceBasePath(), err)
	}
	return nil
}

// GetMainFilePath returns the main.go file path
func (layout *GRPCLayout) GetMainFilePath() string {
	return filepath.Join(layout.GetCommandBasePath(), "main.go")
}

// GetAPIProtobufBasePath returns the API protobuf base path
func (layout *GRPCLayout) GetAPIProtobufBasePath() string {
	return filepath.Join(layout.GetAPIBasePath(), "protobuf_spec")
}

// GetAPIProtobufPath returns the the API protobuf path connected with protobuf file path
func (layout *GRPCLayout) GetAPIProtobufPath(protoFile string) string {
	filename := strings.TrimSuffix(protoFile, filepath.Ext(protoFile))
	return filepath.Join(layout.GetAPIProtobufBasePath(), strcase.ToSnake(filename))
}

// GetAPIProtobufFilePath returns the the API protobuf path connected with protobuf file path and protobuf file name
func (layout *GRPCLayout) GetAPIProtobufFilePath(protoFile string) string {
	filename := strings.TrimSuffix(protoFile, filepath.Ext(protoFile))
	return filepath.Join(layout.GetAPIProtobufBasePath(), strcase.ToSnake(filename), protoFile)
}

// GetGRPCServiceBasePath returns the gateway gRPC base path
func (layout *GRPCLayout) GetGRPCServiceBasePath() string {
	return filepath.Join(layout.GetServiceBasePath(), "grpc")
}

// GetGRPCServicePath returns the the gateway gRPC path connected with service path
func (layout *GRPCLayout) GetGRPCServicePath(service string) string {
	return filepath.Join(layout.GetGRPCServiceBasePath(), strcase.ToSnake(service))
}

// GetGRPCServiceFilePath returns the the gateway tRPC path connected with service path and file name
func (layout *GRPCLayout) GetGRPCServiceFilePath(service string) string {
	return filepath.Join(layout.GetGRPCServiceBasePath(), strcase.ToSnake(service), "implementation.go")
}

// GetGRPCRPCFilePath returns the the gateway tRPC path connected with service path and RPC file name
func (layout *GRPCLayout) GetGRPCRPCFilePath(service, RPC string) string {
	return filepath.Join(layout.GetGRPCServiceBasePath(), strcase.ToSnake(service), strcase.ToSnake(RPC+".go"))
}

// GetConfigGoBasePath returns the config base path
func (layout *GRPCLayout) GetConfigGoBasePath() string {
	return filepath.Join(layout.GetInternalBasePath(), "config")
}

// GetConfigGoFilePath returns the config.go file path
func (layout *GRPCLayout) GetConfigGoFilePath() string {
	return filepath.Join(layout.GetConfigGoBasePath(), "config.go")
}

// GetConfigYAMLFilePath returns the gofra.yaml file path
func (layout *GRPCLayout) GetConfigYAMLFilePath() string {
	return filepath.Join(layout.GetConfigYAMLBasePath(), "gofra.yaml")
}

// GetConfigPackageName returns the config package name
func (layout *GRPCLayout) GetConfigPackageName(goModule string) string {
	relativePath := strings.TrimPrefix(layout.GetConfigGoBasePath(), layout.GetOutputPath())
	return filepath.Join(goModule, relativePath)
}

// GetPackageAlias returns the package name alias
func (layout *GRPCLayout) GetPackageAlias(packageName string) string {
	return strcase.ToLowerCamel(packageName)
}

// GetServicePackagePath returns the import package path according to go module path and service name
func (layout *GRPCLayout) GetServicePackagePath(goModule, service string) string {
	// generate
	relativePath := strings.TrimPrefix(layout.GetGRPCServiceBasePath(), layout.GetOutputPath())
	return filepath.Join(goModule, relativePath, strcase.ToSnake(service))
}

// GetProtoPackagePath returns the import protobuf package path according to go module path and proto file path
func (layout *GRPCLayout) GetProtoPackagePath(goModule, protoFile string) string {
	fileName := strings.TrimSuffix(filepath.Base(protoFile), filepath.Ext(protoFile))

	// generate proto stub
	// something like: healthCheck "git.code.oa.com/foo/api/protobuf-spec/health-check"
	relativePath := strings.TrimSuffix(layout.GetAPIProtobufBasePath(), layout.GetOutputPath())
	return filepath.Join(goModule, relativePath, strcase.ToSnake(fileName))
}

// GetServiceRegisterStub returns the service register line of code
func (layout *GRPCLayout) GetServiceRegisterStub(protoPath, serviceName string) string {
	fileName := strings.TrimSuffix(filepath.Base(protoPath), filepath.Ext(protoPath))

	// generate register stub
	// something like: healthCheckProto.RegisterHealthCheckServer(server, healthCheck.Implementation)
	return fmt.Sprintf("%v.Register%vServer(server, %v.Implementation{})",
		layout.GetPackageAlias(fileName)+"Proto", serviceName, strcase.ToLowerCamel(serviceName))
}

// GetServiceImportStub returns the service import line of code
func (layout *GRPCLayout) GetServiceImportStub(goModule, service string) string {
	return fmt.Sprintf("%v \"%v\"",
		layout.GetPackageAlias(service), layout.GetServicePackagePath(goModule, service))
}

// GetProtoImportStub returns the proto import line of code
func (layout *GRPCLayout) GetProtoImportStub(goModule, protoPath string) string {
	fileName := strings.TrimSuffix(filepath.Base(protoPath), filepath.Ext(protoPath))
	return fmt.Sprintf("%v \"%v\"",
		layout.GetPackageAlias(fileName)+"Proto", layout.GetProtoPackagePath(goModule, protoPath))
}
