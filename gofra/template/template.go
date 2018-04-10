package template

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"errors"
	"strings"
	"text/template"

	"github.com/tallstoat/pbparser"

	commonUtils "github.com/DarkMetrix/gofra/common/utils"
)

//protoc command path
var ProtocPath string = "protoc"

//Server config
type ServerInfo struct {
	Addr string `json:"addr"`
}

//Client config
type ClientInfo struct {
	Pool PoolInfo `json:"pool"`
}

//Pool config
type PoolInfo struct {
	InitConns int `json:"init_conns"`
	MaxConns int `json:"max_conns"`
	IdleTime int `json:"idle_time"`
}

//Monitor package
type MonitorPackageInfo struct {
	Package string `json:"package"`
	InitParam string `json:"init_param"`
}

//Tracing package
type TracingPackageInfo struct {
	Package string `json:"package"`
	InitParam string `json:"init_param"`
}

//Interceptor package
type InterceptorPackageInfo struct {
	MonitorPackage string `json:"monitor_package"`
	TracingPackage string `json:"tracing_package"`
}

//Template info
type TemplateInfo struct {
	Author string `json:"author"`
	Project string `json:"project"`
	Version string `json:"version"`

	Server ServerInfo `json:"server"`
	Client ClientInfo `json:"client"`

	MonitorPackage MonitorPackageInfo `json:"monitor_package"`
	TracingPackage TracingPackageInfo `json:"tracing_package"`

	InterceptorPackage InterceptorPackageInfo `json:"interceptor_package"`
}

//Generate template.json
func GenerateTemplateJsonFile(workingPath string, override bool) error {
	filePath := filepath.Join(workingPath, "template.json.default")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	jsonTemplate, err := template.New("template_json").Parse(JsonTemplate)

	if err != nil {
		return err
	}

	jsonInfo := &JsonInfo{
		Author: "Author Name",
		Project: "Project Name",
		Addr: "localhost:58888",
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
func GenerateCommonFile(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "src", "common", "common.go")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
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
func GenerateConfigFile(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "src", "config", "config.go")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
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

//Generate config.json
func GenerateConfigJsonFile(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "conf", "config.json")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	configJsonTemplate, err := template.New("config_json").Parse(ConfigJsonTemplate)

	if err != nil {
		return err
	}

	configJsonInfo := &ConfigJsonInfo{
		Addr: info.Server.Addr,
		InitConns: info.Client.Pool.InitConns,
		MaxConns: info.Client.Pool.MaxConns,
		IdleTime: info.Client.Pool.IdleTime,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = configJsonTemplate.Execute(file, configJsonInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate naming.json
func GenerateNamingJsonFile(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "conf", "naming.json")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	namingJsonTemplate, err := template.New("naming_json").Parse(NamingJsonTemplate)

	if err != nil {
		return err
	}

	namingJsonInfo := &NamingJsonInfo{
		Project: info.Project,
		Addr: info.Server.Addr,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = namingJsonTemplate.Execute(file, namingJsonInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate log.config
func GenerateConfigLogFile(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "conf", "log.config")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
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
		MaxSize: 524288000,
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
func GenerateApplicationFile(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "src", "application", "application.go")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "src") + "/")

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
func GenerateMainFile(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "src", "main.go")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "src") + "/")

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
func GenerateHealthCheckHandler(workingPath, goPath string, info *TemplateInfo, override bool) error {
	//Generate health check proto file
	err := GenerateHealthCheckProto(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	filePath := filepath.Join(workingPath, "src", "proto", "health_check", "health_check.proto")

	err = GenerateService(workingPath, goPath, info, filePath, override, false)

	if err != nil {
		return err
	}

	return nil
}

//Generate health check service proto
func GenerateHealthCheckProto(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "src", "proto", "health_check", "health_check.proto")
	filePathRelative := filepath.Join(".", "src", "proto", "health_check", "health_check.proto")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
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
func GenerateService(workingPath, goPath string, info *TemplateInfo, protoPath string, override bool, update bool) error {
	//Parse proto file
	pf, err := pbparser.ParseFile(protoPath)

	if err != nil {
		return errors.New(fmt.Sprintf("Unable to parse proto file! error:%v", err.Error()))
	}

	workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "src") + "/")
	protoFileNamePrefix := strings.TrimSuffix(filepath.Base(protoPath), filepath.Ext(protoPath))

	//Generate service and handler
	for _, elem := range pf.Services {
		//Create path
		handlerPath := filepath.Join(workingPath, "src", "handler", elem.Name)

		err := commonUtils.CreatePath(handlerPath, override)

		if err != nil {
			return err
		}

		service := &ServiceInfo{
			Author: info.Author,
			Time: time.Now().Format("2006-01-02 15:04:05"),
			ServiceName: elem.Name}

		//Create implementation file
		err = GenerateServiceImplementation(workingPath, goPath, info, service, override, update)

		if err != nil {
			return err
		}

		//Add service handler import to application
		err = AddServiceHandlerToApplicationFileImport(workingPath, goPath, info, service.ServiceName)

		if err != nil {
			return err
		}

		//Add service register to application
		err = AddServiceRegisterToApplicationFile(workingPath, goPath, info, protoPath, service.ServiceName)

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

				err = GenerateServiceHandler(workingPath, goPath, info, rpc, override, update)

				if err != nil {
					return err
				}
		}
	}

	//Add service proto import to application
	err = AddServiceProtoToApplicationFileImport(workingPath, goPath, info, protoPath)

	if err != nil {
		return err
	}

	return nil
}

//Add service proto import to application file
func AddServiceProtoToApplicationFileImport(workingPath, goPath string, info *TemplateInfo, protoPath string) error {
	applicationFilePath := filepath.Join(workingPath, "src", "application", "application.go")

	applicationContent, err := ioutil.ReadFile(applicationFilePath)

	if err != nil {
		return err
	}

	workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "src") + "/")
	protoFileNamePrefix := strings.TrimSuffix(filepath.Base(protoPath), filepath.Ext(protoPath))

	protoImport := fmt.Sprintf("%v \"%v/src/proto/%v\"", protoFileNamePrefix, workingPathRelative, protoFileNamePrefix)
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
func AddServiceHandlerToApplicationFileImport(workingPath, goPath string, info *TemplateInfo, serviceName string) error {
	applicationFilePath := filepath.Join(workingPath, "src", "application", "application.go")

	applicationContent, err := ioutil.ReadFile(applicationFilePath)

	if err != nil {
		return err
	}

	workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "src") + "/")

	protoImport := fmt.Sprintf("%vHandler \"%v/src/handler/%v\"", serviceName, workingPathRelative, serviceName)
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
func AddServiceRegisterToApplicationFile(workingPath, goPath string, info *TemplateInfo, protoPath, serviceName string) error {
	applicationFilePath := filepath.Join(workingPath, "src", "application", "application.go")

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
func GenerateServiceImplementation(workingPath, goPath string, info *TemplateInfo, service *ServiceInfo, override bool, update bool) error {
	filePath := filepath.Join(workingPath, "src", "handler", service.ServiceName, service.ServiceName + ".go")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
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
func GenerateServiceHandler(workingPath, goPath string, info *TemplateInfo, rpc *RpcInfo, override bool, update bool) error {
	filePath := filepath.Join(workingPath, "src", "handler", rpc.ServiceName, rpc.RpcName + ".go")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
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

//Generate test client
func GenerateTestClient(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "test", "main.go")

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
	}

	if isExist && override {
		err := os.RemoveAll(filePath)

		if err != nil {
			return err
		}
	}

	workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "src") + "/")

	//Parse template
	testClientTemplate, err := template.New("test_client").Parse(TestClientTemplate)

	if err != nil {
		return err
	}

	testClientInfo := &TestClientInfo{
		Author: info.Author,
		Time: time.Now().Format("2006-01-02 15:04:05"),
		Project: info.Project,
		Addr: info.Server.Addr,
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
	err = testClientTemplate.Execute(file, testClientInfo)

	if err != nil {
		return err
	}

	return nil
}
