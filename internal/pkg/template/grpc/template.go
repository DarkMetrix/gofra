package grpc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/tallstoat/pbparser"

	"github.com/DarkMetrix/gofra/internal/pkg/utils/pbparser_import_provider"
	commonUtils "github.com/DarkMetrix/gofra/pkg/utils"
	gofraTemplate "github.com/DarkMetrix/gofra/internal/pkg/template"
)

//protoc command path
var ProtocPath string = "protoc"

//Generate template.json
func GenerateTemplateJsonFile(workingPath string, override bool, jsonInfo gofraTemplate.JsonInfo) error {
	filePath := filepath.Join(workingPath, "template.json")

	//Check file is exist or not
	isExist, err := commonUtils.CheckPathExists(filePath)

	if err != nil {
		return err
	}

	if isExist && !override {
		filePathRel, err := filepath.Rel(workingPath, filePath)

		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will override it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	jsonTemplate, err := template.New("template_json").Parse(gofraTemplate.JsonTemplate)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	//Render template to file
	err = jsonTemplate.Execute(file, jsonInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate common.go
func GenerateCommonFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "internal", "pkg", "common", "common.go")

	//Check file is exist or not
	isExist, err := commonUtils.CheckPathExists(filePath)

	if err != nil {
		return err
	}

	if isExist && !override {
		filePathRel, err := filepath.Rel(workingPath, filePath)

		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will override it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	commonTemplate, err := template.New("common").Parse(CommonTemplate)

	if err != nil {
		return err
	}

	commonInfo := &CommonInfo{
		Author: info.Author,
		Time: time.Now().Format("2006-01-02 15:04:05"),
		Project: info.Project,
		Version: info.Version,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = commonTemplate.Execute(file, commonInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate config.go
func GenerateConfigFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "internal", "pkg", "config", "config.go")

	//Check file is exist or not
	isExist, err := commonUtils.CheckPathExists(filePath)

	if err != nil {
		return err
	}

	if isExist && !override {
		filePathRel, err := filepath.Rel(workingPath, filePath)

		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will override it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	configTemplate, err := template.New("config").Parse(ConfigTemplate)

	if err != nil {
		return err
	}

	configInfo := &ConfigInfo{
		Author: info.Author,
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = configTemplate.Execute(file, configInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate config.toml
func GenerateConfigTomlFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "configs", "config.toml")

	//Check file is exist or not
	isExist, err := commonUtils.CheckPathExists(filePath)

	if err != nil {
		return err
	}

	if isExist && !override {
		filePathRel, err := filepath.Rel(workingPath, filePath)

		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will override it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	configTomlTemplate, err := template.New("config_toml").Parse(ConfigTomlTemplate)

	if err != nil {
		return err
	}

	configTomlInfo := &ConfigTomlInfo{
		Addr: info.Server.Addr,

		MonitorInitParams: info.MonitorPackage.InitParam,
		TracingInitParams: info.TracingPackage.InitParam,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = configTomlTemplate.Execute(file, configTomlInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate log.config
func GenerateConfigLogFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "configs", "log.config")

	//Check file is exist or not
	isExist, err := commonUtils.CheckPathExists(filePath)

	if err != nil {
		return err
	}

	if isExist && !override {
		filePathRel, err := filepath.Rel(workingPath, filePath)

		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will override it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	configLogTemplate, err := template.New("config_log").Parse(LogTemplate)

	if err != nil {
		return err
	}

	configLogInfo := &LogInfo{
		Path: fmt.Sprintf("../log/%v.log", info.Project),
		MaxSize: 100 * 1024 * 1024,
		MaxRolls: 10,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = configLogTemplate.Execute(file, configLogInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate application.go
func GenerateApplicationFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "internal", "app", "application.go")

	//Check file is exist or not
	isExist, err := commonUtils.CheckPathExists(filePath)

	if err != nil {
		return err
	}

	if isExist && !override {
		filePathRel, err := filepath.Rel(workingPath, filePath)

		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will override it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "src") + "/")
	workingPathRelative := filepath.Base(workingPath)

	//Parse template
	applicationTemplate, err := template.New("application").Parse(ApplicationTemplate)

	if err != nil {
		return err
	}

	applicationInfo := &ApplicationInfo{
		Author: info.Author,
		Time: time.Now().Format("2006-01-02 15:04:05"),
		Project: info.Project,
		WorkingPathRelative: workingPathRelative,
		MonitorPackage: info.MonitorPackage.Package,
		MonitorInitParam: info.MonitorPackage.InitParam,
        TracingPackage: info.TracingPackage.Package,
        TracingInitParam: info.TracingPackage.InitParam,
		MonitorInterceptorPackage: info.InterceptorPackage.MonitorPackage,
		TracingInterceptorPackage: info.InterceptorPackage.TracingPackage,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = applicationTemplate.Execute(file, applicationInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate main.go
func GenerateMainFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "cmd", "main.go")

	//Check file is exist or not
	isExist, err := commonUtils.CheckPathExists(filePath)

	if err != nil {
		return err
	}

	if isExist && !override {
		filePathRel, err := filepath.Rel(workingPath, filePath)

		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will override it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "src") + "/")
	workingPathRelative := filepath.Base(workingPath)

	//Parse template
	mainTemplate, err := template.New("main").Parse(MainTemplate)

	if err != nil {
		return err
	}

	mainInfo := &MainInfo{
		Author: info.Author,
		Time: time.Now().Format("2006-01-02 15:04:05"),
		Project: info.Project,
		WorkingPathRelative: workingPathRelative,
		Addr: info.Server.Addr,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = mainTemplate.Execute(file, mainInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate health check handler
func GenerateHealthCheckHandler(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	//Generate health check proto file
	err := GenerateHealthCheckProto(workingPath, info, override)

	if err != nil {
		return err
	}

	filePath := filepath.Join(workingPath, "api", "protobuf_spec", "health_check", "health_check.proto")

	err = GenerateService(workingPath, []string{"./api/protobuf_spec/health_check"}, info, filePath, override, false)

	if err != nil {
		return err
	}

	return nil
}

//Generate health check service proto
func GenerateHealthCheckProto(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "api", "protobuf_spec", "health_check", "health_check.proto")
	filePathRelative := filepath.Join(".", "api", "protobuf_spec", "health_check", "health_check.proto")

	//Check file is exist or not
	isExist, err := commonUtils.CheckPathExists(filePath)

	if err != nil {
		return err
	}

	if isExist && !override {
		filePathRel, err := filepath.Rel(workingPath, filePath)

		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will override it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	healthCheckServiceProtoTemplate, err := template.New("health_check").Parse(HealthCheckServiceProtoTemplate)

	if err != nil {
		return err
	}

	mainInfo := &HealthCheckServiceProtoInfo{
		Author: info.Author,
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = healthCheckServiceProtoTemplate.Execute(file, mainInfo)

	if err != nil {
		return err
	}

	//Execute protoc to generate .pb.go file
	shellCmd := exec.Command(ProtocPath, "--go_out=plugins=grpc:.", filePathRelative)

	err = shellCmd.Run()

	if err != nil {
		return errors.New(fmt.Sprintf("%v --go_out=plugins=grpc:. %v failed! error:%v",
			ProtocPath, filePathRelative, err.Error()))
	}

	return nil
}

//Generate service
func GenerateService(workingPath string, protoFileIncludePath []string, info *gofraTemplate.TemplateInfo, protoPath string, override bool, update bool) error {
	//Parse proto file
	raw, err := ioutil.ReadFile(protoPath)

	if err != nil {
		return errors.New(fmt.Sprintf("Unable to read proto file! error:%v", err.Error()))
	}

	r := strings.NewReader(string(raw[:]))

	// implement a dir based import module provider which reads
	// import modules from the same dir as the original proto file
	importProvider := pbparser_import_provider.GofraImportModuleProvider{ModuleSearchPaths:protoFileIncludePath}

	// invoke Parse() API to parse the file
	pf, err := pbparser.Parse(r, &importProvider)

	if err != nil {
		return errors.New(fmt.Sprintf("Unable to parse proto file! error:%v", err.Error()))
	}

	//workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "src") + "/")
	workingPathRelative := filepath.Base(workingPath)
	protoFileNamePrefix := strings.TrimSuffix(filepath.Base(protoPath), filepath.Ext(protoPath))

	//Generate service and handler
	for _, elem := range pf.Services {
		//Create path
		handlerPath := filepath.Join(workingPath, "internal", "grpc_handler", elem.Name)

		err := commonUtils.CreatePath(handlerPath, override)

		if err != nil {
			return err
		}

		service := &ServiceInfo{
			Author: info.Author,
			Time: time.Now().Format("2006-01-02 15:04:05"),
			ServiceName: elem.Name}

		//Create implementation file
		err = GenerateServiceImplementation(workingPath, info, service, override, update)

		if err != nil {
			return err
		}

		//Add service handler import to application
		err = AddServiceHandlerToApplicationFileImport(workingPath, info, service.ServiceName)

		if err != nil {
			return err
		}

		//Add service register to application
		err = AddServiceRegisterToApplicationFile(workingPath, info, protoPath, service.ServiceName)

		if err != nil {
			return err
		}

		//Create handlers
		for _, rpc := range elem.RPCs {
			rpc := &RpcInfo{
				Author: info.Author,
				Time: time.Now().Format("2006-01-02 15:04:05"),
				WorkingPathRelative: workingPathRelative,
				ServiceName: service.ServiceName,
				FileNamePrefix: protoFileNamePrefix,
				RpcName: rpc.Name,
				Request: rpc.RequestType.Name(),
				Response: rpc.ResponseType.Name(),
				MonitorPackage: info.MonitorPackage.Package,
				TracingPackage: info.TracingPackage.Package,}

				err = GenerateServiceHandler(workingPath, info, rpc, override, update)

				if err != nil {
					return err
				}
		}
	}

	//Add service proto import to application
	err = AddServiceProtoToApplicationFileImport(workingPath, info, protoPath)

	if err != nil {
		return err
	}

	return nil
}

//Add service proto import to application file
func AddServiceProtoToApplicationFileImport(workingPath string, info *gofraTemplate.TemplateInfo, protoPath string) error {
	applicationFilePath := filepath.Join(workingPath, "internal", "app", "application.go")

	applicationContent, err := ioutil.ReadFile(applicationFilePath)

	if err != nil {
		return err
	}

	//workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "api") + "/")
	workingPathRelative := filepath.Base(workingPath)
	protoFileNamePrefix := strings.TrimSuffix(filepath.Base(protoPath), filepath.Ext(protoPath))

	protoImport := fmt.Sprintf("%v \"%v/api/protobuf_spec/%v\"", protoFileNamePrefix, workingPathRelative, protoFileNamePrefix)
	protoImportStub := fmt.Sprintf("%v\r\n	/*@PROTO_STUB*/", protoImport)

	if strings.Contains(string(applicationContent), protoImport) {
		return nil
	}

	applicationContent = []byte(strings.Replace(string(applicationContent), "/*@PROTO_STUB*/", protoImportStub, 1))

	err = ioutil.WriteFile(applicationFilePath, applicationContent, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

//Add service handler import to application file
func AddServiceHandlerToApplicationFileImport(workingPath string, info *gofraTemplate.TemplateInfo, serviceName string) error {
	applicationFilePath := filepath.Join(workingPath, "internal", "app", "application.go")

	applicationContent, err := ioutil.ReadFile(applicationFilePath)

	if err != nil {
		return err
	}

	//workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "internal") + "/")
	workingPathRelative := filepath.Base(workingPath)

	protoImport := fmt.Sprintf("%vHandler \"%v/internal/grpc_handler/%v\"", serviceName, workingPathRelative, serviceName)
	protoImportStub := fmt.Sprintf("%v\r\n	/*@HANDLER_STUB*/", protoImport)

	if strings.Contains(string(applicationContent), protoImport) {
		return nil
	}

	applicationContent = []byte(strings.Replace(string(applicationContent), "/*@HANDLER_STUB*/", protoImportStub, 1))

	err = ioutil.WriteFile(applicationFilePath, applicationContent, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

//Add service handler import to application file
func AddServiceRegisterToApplicationFile(workingPath string, info *gofraTemplate.TemplateInfo, protoPath, serviceName string) error {
	applicationFilePath := filepath.Join(workingPath, "internal", "app", "application.go")

	applicationContent, err := ioutil.ReadFile(applicationFilePath)

	if err != nil {
		return err
	}

	protoFileNamePrefix := strings.TrimSuffix(filepath.Base(protoPath), filepath.Ext(protoPath))

	protoImport := fmt.Sprintf("%v.Register%vServer(s, %vHandler.%vImpl{})", protoFileNamePrefix, serviceName, serviceName, serviceName)
	protoImportStub := fmt.Sprintf("%v\r\n	/*@REGISTER_STUB*/", protoImport)

	if strings.Contains(string(applicationContent), protoImport) {
		return nil
	}

	applicationContent = []byte(strings.Replace(string(applicationContent), "/*@REGISTER_STUB*/", protoImportStub, 1))

	err = ioutil.WriteFile(applicationFilePath, applicationContent, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

//Generate service implementation
func GenerateServiceImplementation(workingPath string, info *gofraTemplate.TemplateInfo, service *ServiceInfo, override bool, update bool) error {
	filePath := filepath.Join(workingPath, "internal", "grpc_handler", service.ServiceName, service.ServiceName + ".go")

	//Check file is exist or not
	isExist, err := commonUtils.CheckPathExists(filePath)

	if err != nil {
		return err
	}

	if isExist && !override {
		if update {
			return nil
		}

		filePathRel, err := filepath.Rel(workingPath, filePath)

		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will override it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	serviceTemplate, err := template.New("implementation").Parse(ServiceTemplate)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = serviceTemplate.Execute(file, service)

	if err != nil {
		return err
	}

	return nil
}

//Generate service handler
func GenerateServiceHandler(workingPath string, info *gofraTemplate.TemplateInfo, rpc *RpcInfo, override bool, update bool) error {
	filePath := filepath.Join(workingPath, "internal", "grpc_handler", rpc.ServiceName, rpc.RpcName + ".go")

	//Check file is exist or not
	isExist, err := commonUtils.CheckPathExists(filePath)

	if err != nil {
		return err
	}

	if isExist && !override {
		if update {
			return nil
		}

		filePathRel, err := filepath.Rel(workingPath, filePath)

		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will override it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	serviceRpcTemplate, err := template.New("handler").Parse(ServiceRpcTemplate)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = serviceRpcTemplate.Execute(file, rpc)

	if err != nil {
		return err
	}

	return nil
}

