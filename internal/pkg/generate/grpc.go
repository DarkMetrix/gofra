package generate

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DarkMetrix/gofra/internal/pkg/directory"
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/templates/general"
	"github.com/DarkMetrix/gofra/internal/pkg/templates/grpc"
	"github.com/DarkMetrix/gofra/pkg/utils"
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
	options := option.NewOptions(opts...)

	// build args which includes proto file include path
	protoFileIncludePath := append(options.ProtoFileIncludePath, layout.GetOutputPath())
	args := []string{}
	for _, path := range protoFileIncludePath {
		arg := fmt.Sprintf("--proto_path=%v", path)
		args = append(args, arg)
	}
	args = append(args, "--go_out=plugins=grpc:.")
	args = append(args, protoPath)

	// execute protoc to generate .pb.go file
	shellCmd := exec.Command(options.ProtocPath, args...)
	if err := shellCmd.Run(); err != nil {
		return xerrors.Errorf("%v %v failed! error:%v", options.ProtocPath, args, err.Error())
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
		serviceDescs := fileDesc.GetServices()

		for _, serviceDesc := range serviceDescs {
			opts = append(opts, option.WithServiceName(serviceDesc.GetName()))
			serviceInfo := grpc.NewServiceInfo(opts...)

			// create path
			handlerPath := layout.GetGRPCServicePath(serviceInfo.ServiceName)
			if err := utils.CreatePath(handlerPath, options.Override); err != nil {
				return xerrors.Errorf("create service path failed! error:%w", err)
			}

			// create implementation file
			if err := grpc.NewServiceInfo(opts...).RenderFile(
				layout.GetGRPCServiceFilePath(serviceDesc.GetName())); err != nil {
				return xerrors.Errorf("create service implementation file failed! error:%w", err)
			}

			// create RPC file
			opts = append(opts,
				option.WithPackageName(serviceDesc.GetName()),
				option.WithImportedPackageName(layout.GetProtoPackagePath(options.GoModule, protoPath)),
			)

			for _, rpcDesc := range serviceDesc.GetMethods() {
				opts = append(opts,
					option.WithRPCName(rpcDesc.GetName()),
					option.WithRequestName(rpcDesc.GetInputType().GetName()),
					option.WithResponseName(rpcDesc.GetOutputType().GetName()),
				)

				if err := grpc.NewRPCInfo(opts...).RenderFile(
					layout.GetGRPCRPCFilePath(serviceDesc.GetName(), rpcDesc.GetName())); err != nil {
					return xerrors.Errorf("create RPC implementation file failed! error:%w", err)
				}
			}

			// add service import to main.go
			if err := AddServiceImportToMain(layout, options.GoModule, serviceDesc.GetName()); err != nil {
				return xerrors.Errorf("add service import to main.go failed! error:%w", err)
			}

			// add service register to main.go
			if err := AddServiceRegisterToMain(layout, serviceDesc.GetName(), protoPath); err != nil {
				return xerrors.Errorf("add service import to main.go failed! error:%w", err)
			}
		}
	}

	// add service proto import to application
	if err := AddServiceProtoImportToMain(layout, options.GoModule, protoPath); err != nil {
		return xerrors.Errorf("add service proto import to main.go failed! error:%w", err)
	}

	return nil
}

// AddServiceImportToMain adds service import to main.go
func AddServiceImportToMain(layout directory.GRPCServiceLayout, goModule, serviceName string) error {
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

// AddServiceRegisterToMain adds service register to main.go
func AddServiceRegisterToMain(layout directory.GRPCServiceLayout, serviceName, protoPath string) error {
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

// AddServiceProtoImportToMain adds service proto import to main.go
func AddServiceProtoImportToMain(layout directory.GRPCServiceLayout, goModule, protoPath string) error {
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
