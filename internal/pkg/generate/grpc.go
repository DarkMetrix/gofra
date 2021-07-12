package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/DarkMetrix/gofra/internal/pkg/directory"
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/pb"
	"github.com/DarkMetrix/gofra/internal/pkg/templates/general"
	"github.com/DarkMetrix/gofra/internal/pkg/templates/grpc"
	"github.com/DarkMetrix/gofra/pkg/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"golang.org/x/xerrors"
)

// InitGRPCDirectoryStructure initializes the gRPC service directory to output path
func InitGRPCDirectoryStructure(opts ...option.Option) (directory.GRPCServiceLayout, error) {
	layout := directory.NewGRPCLayout(opts...)

	// save directory structure
	if err := layout.Save(); err != nil {
		return nil, xerrors.Errorf("layout.Save failed! error:%w", err)
	}
	return layout, nil
}

// InitGRPCService initializes all the files needed for a basic gRPC service
func InitGRPCService(layout directory.GRPCServiceLayout, opts ...option.Option) error {
	// create go module file
	if err := general.NewGoModuleInfo(opts...).RenderFile(layout.GetGoModuleFilePath()); err != nil {
		return xerrors.Errorf("create go module file failed! error:%w", err)
	}

	// create main file
	options := option.NewOptions(opts...)
	opts = append(opts, option.WithConfigPackagePath(layout.GetConfigPackageName(options.GoModule)))
	if err := grpc.NewMainInfo(opts...).RenderFile(layout.GetMainFilePath()); err != nil {
		return xerrors.Errorf("create main file failed! error:%w", err)
	}

	// create config file
	if err := grpc.NewConfigInfo(opts...).RenderFile(layout.GetConfigGoFilePath()); err != nil {
		return xerrors.Errorf("create config file failed! error:%w", err)
	}

	// create gofra.yaml file
	if err := grpc.NewConfigYAMLInfo(opts...).RenderFile(layout.GetConfigYAMLFilePath()); err != nil {
		return xerrors.Errorf("create config YAML file failed! error:%w", err)
	}

	// create health check proto file
	if err := utils.CreatePaths(true, filepath.Join(layout.GetAPIProtobufBasePath(), "health_check")); err != nil {
		return xerrors.Errorf("create health check directory failed! error:%w", err)
	}
	healthCheckProtoPath := filepath.Join(layout.GetAPIProtobufBasePath(), "health_check", "health_check.proto")
	if err := grpc.NewHealthCheckProtoInfo(opts...).RenderFile(healthCheckProtoPath); err != nil {
		return xerrors.Errorf("create health check proto file failed! error:%w", err)
	}

	// add service health check
	if err := AddGRPCService(healthCheckProtoPath, layout, opts...); err != nil {
		return xerrors.Errorf("add gRPC service failed! error:%w", err)
	}
	return nil
}

// AddGRPCService adds a gRPC service according to proto file
func AddGRPCService(protoPath string, layout directory.GRPCServiceLayout, opts ...option.Option) error {
	opts = append(opts, option.WithIgnoreExist(false))
	return generateGRPCService(protoPath, layout, opts...)
}

// UpdateGRPCService updates a gRPC service according to proto file
func UpdateGRPCService(protoPath string, layout directory.GRPCServiceLayout, opts ...option.Option) error {
	opts = append(opts, option.WithIgnoreExist(true))
	return generateGRPCService(protoPath, layout, opts...)
}

// generateGRPCService generates a gRPC service according to proto file
func generateGRPCService(protoPath string, layout directory.GRPCServiceLayout, opts ...option.Option) error {
	options := option.NewOptions(opts...)

	// compile gRPC service
	protoFileIncludePath := append(options.ProtoFileIncludePath, options.OutputPath)
	if err := pb.CompileGRPC(options.ProtocPath, protoPath, protoFileIncludePath); err != nil {
		return xerrors.Errorf("pb.CompileGRPC failed! error:%w", err)
	}

	// parse .proto file
	parser := protoparse.Parser{
		ImportPaths: protoFileIncludePath,
	}
	fileDescs, err := parser.ParseFiles(protoPath)
	if err != nil {
		return xerrors.Errorf("Unable to parse proto file! proto file:%v, error:%w", protoPath, err)
	}

	for _, fileDesc := range fileDescs {
		for _, serviceDesc := range fileDesc.GetServices() {
			// generate service related files
			if err := generateServiceFiles(protoPath, true, serviceDesc, layout, opts...); err != nil {
				return xerrors.Errorf("generateServiceFiles failed! error:%w", err)
			}

			// add service import to main.go
			if err := addServiceImportToMain(layout, options.GoModule, serviceDesc.GetName()); err != nil {
				return xerrors.Errorf("add service import to main.go failed! error:%w", err)
			}

			// add service register to main.go
			if err := addServiceRegisterToMain(layout, serviceDesc.GetName(), protoPath); err != nil {
				return xerrors.Errorf("add service import to main.go failed! error:%w", err)
			}
		}
	}

	// add service proto import to application
	if err := addServiceProtoImportToMain(layout, options.GoModule, protoPath); err != nil {
		return xerrors.Errorf("add service proto import to main.go failed! error:%w", err)
	}
	return nil
}

// generateServiceFiles generates gRPC service related files
func generateServiceFiles(protoPath string, update bool, serviceDesc *desc.ServiceDescriptor,
	layout directory.GRPCServiceLayout, opts ...option.Option) error {
	opts = append(opts, option.WithServiceName(serviceDesc.GetName()))
	options := option.NewOptions(opts...)
	serviceInfo := grpc.NewServiceInfo(opts...)

	// create path
	handlerPath := layout.GetGRPCServicePath(serviceInfo.ServiceName)
	if err := utils.CreatePath(handlerPath, options.Override); err != nil {
		return xerrors.Errorf("create service path failed! error:%w", err)
	}

	// create implementation file
	if err := grpc.NewServiceInfo(opts...).RenderFile(layout.GetGRPCServiceFilePath(serviceDesc.GetName())); err != nil {
		return xerrors.Errorf("create service implementation file failed! error:%w", err)
	}

	// create RPC file
	opts = append(opts,
		option.WithPackageName(serviceDesc.GetName()),
		option.WithImportedPackageName(layout.GetProtoPackagePath(options.GoModule, protoPath)),
	)

	for _, rpcDesc := range serviceDesc.GetMethods() {
		if err := generateRPCFiles(update, serviceDesc, rpcDesc, layout, opts...); err != nil {
			return xerrors.Errorf("generateRPCFiles failed! error:%w", err)
		}
	}
	return nil
}

// generateRPCFiles generates gRPC RPC method related files
func generateRPCFiles(update bool, serviceDesc *desc.ServiceDescriptor, rpcDesc *desc.MethodDescriptor,
	layout directory.GRPCServiceLayout, opts ...option.Option) error {
	// generate
	opts = append(opts,
		option.WithRPCName(rpcDesc.GetName()),
		option.WithRequestName(rpcDesc.GetInputType().GetName()),
		option.WithResponseName(rpcDesc.GetOutputType().GetName()),
	)

	if err := grpc.NewRPCInfo(opts...).RenderFile(
		layout.GetGRPCRPCFilePath(serviceDesc.GetName(), rpcDesc.GetName())); err != nil {
		return xerrors.Errorf("create RPC implementation file failed! error:%w", err)
	}
	return nil
}

// addServiceImportToMain adds service import to main.go
func addServiceImportToMain(layout directory.GRPCServiceLayout, goModule, serviceName string) error {
	// read file content
	mainFilePath := layout.GetMainFilePath()
	mainContent, err := ioutil.ReadFile(mainFilePath)
	if err != nil {
		return xerrors.Errorf("ioutil.ReadFile failed! error:%w", err)
	}

	// generate
	protoImport := layout.GetServiceImportStub(goModule, serviceName)
	protoImportStub := fmt.Sprintf("%v\r\n	/*@HANDLER_STUB*/", protoImport)
	if strings.Contains(string(mainContent), protoImport) {
		return nil
	}
	mainContent = []byte(strings.Replace(string(mainContent), "/*@HANDLER_STUB*/", protoImportStub, 1))

	// write to file
	if err := ioutil.WriteFile(mainFilePath, mainContent, os.ModePerm); err != nil {
		return xerrors.Errorf("ioutil.WriteFile failed! error:%w", err)
	}
	return nil
}

// addServiceRegisterToMain adds service register to main.go
func addServiceRegisterToMain(layout directory.GRPCServiceLayout, serviceName, protoPath string) error {
	// read file content
	mainFilePath := layout.GetMainFilePath()
	mainContent, err := ioutil.ReadFile(mainFilePath)
	if err != nil {
		return xerrors.Errorf("ioutil.ReadFile failed! error:%w", err)
	}

	// generate register stub
	// something like: healthcheck.RegisterHealthCheckService(server, healthcheck.Implenmentation)
	protoImport := layout.GetServiceRegisterStub(protoPath, serviceName)
	protoImportStub := fmt.Sprintf("%v\r\n	/*@REGISTER_STUB*/", protoImport)
	if strings.Contains(string(mainContent), protoImport) {
		return nil
	}
	mainContent = []byte(strings.Replace(string(mainContent), "/*@REGISTER_STUB*/", protoImportStub, 1))

	// write to file
	if err = ioutil.WriteFile(mainFilePath, mainContent, os.ModePerm); err != nil {
		return xerrors.Errorf("ioutil.WriteFile failed! error:%w", err)
	}
	return nil
}

// addServiceProtoImportToMain adds service proto import to main.go
func addServiceProtoImportToMain(layout directory.GRPCServiceLayout, goModule, protoPath string) error {
	// read file content
	mainFilePath := layout.GetMainFilePath()
	mainContent, err := ioutil.ReadFile(mainFilePath)
	if err != nil {
		return xerrors.Errorf("ioutil.ReadFile failed! error:%w", err)
	}

	// generate proto stub
	// something like: healthCheck "git.code.oa.com/foo/api/protobuf-spec/health-check"
	protoImport := layout.GetProtoImportStub(goModule, protoPath)
	protoImportStub := fmt.Sprintf("%v\r\n	/*@PROTO_STUB*/", protoImport)
	if strings.Contains(string(mainContent), protoImport) {
		return nil
	}
	mainContent = []byte(strings.Replace(string(mainContent), "/*@PROTO_STUB*/", protoImportStub, 1))

	// write to file
	if err = ioutil.WriteFile(mainFilePath, mainContent, os.ModePerm); err != nil {
		return xerrors.Errorf("ioutil.WriteFile failed! error:%w", err)
	}
	return nil
}
