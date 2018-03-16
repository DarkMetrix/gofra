package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"errors"
	"strings"
	"text/template"

	"github.com/tallstoat/pbparser"

	gofraUtils "github.com/DarkMetrix/gofra/gofra/utils"
)

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

//Interceptor package
type InterceptorPackageInfo struct {
	Package string `json:"package"`
}

//Template info
type TemplateInfo struct {
	Author string `json:"author"`
	Project string `json:"project"`
	Version string `json:"version"`
	Server ServerInfo `json:"server"`
	Client ClientInfo `json:"client"`
	MonitorPackage MonitorPackageInfo `json:"monitor_package"`
	InterceptorPackage InterceptorPackageInfo `json:"interceptor_package"`
}

//Generate common.go
func GenerateCommonFile(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "src", "common", "common.go")

	//Check file is exist or not
	isExist, err := gofraUtils.CheckPathExists(filePath)

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
	isExist, err := gofraUtils.CheckPathExists(filePath)

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
	isExist, err := gofraUtils.CheckPathExists(filePath)

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

//Generate log.config
func GenerateConfigLogFile(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "conf", "log.config")

	//Check file is exist or not
	isExist, err := gofraUtils.CheckPathExists(filePath)

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
	isExist, err := gofraUtils.CheckPathExists(filePath)

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

	workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "src") + "/")

	//Parse template
	applicationTemplate, err := template.New("application").Parse(ApplicationTemplate)

	if err != nil {
		return err
	}

	allServices := make([]string, 0)
	allServices = append(allServices, "HealthCheckService")

	applicationInfo := &ApplicationInfo{
		Author: info.Author,
		Time: time.Now().Format("2006-01-02 15:04:05"),
		WorkingPathRelative: workingPathRelative,
		MonitorPackage: info.MonitorPackage.Package,
		MonitorInitParam: info.MonitorPackage.InitParam,
		InterceptorPackage: info.InterceptorPackage.Package,
		Services: allServices,
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
	isExist, err := gofraUtils.CheckPathExists(filePath)

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
	err := GenerateHealthCheckProto(workingPath, goPath, info, override)

	if err != nil {
		return err
	}

	// invoke ParseFile() API to parse the file
	filePath := filepath.Join(workingPath, "src", "proto", "health_check", "health_check.proto")

	pf, err := pbparser.ParseFile(filePath)

	if err != nil {
		return errors.New(fmt.Sprintf("Unable to parse proto file! error:%v", err.Error()))
	}

	workingPathRelative := strings.TrimPrefix(workingPath, filepath.Join(goPath, "src") + "/")

	// print attributes of the returned datastructure
	for _, elem := range pf.Services {
		//Create path
		handlerPath := filepath.Join(workingPath, "src", "handler", elem.Name)

		err := gofraUtils.CreatePath(handlerPath)

		if err != nil {
			return err
		}

		service := &ServiceInfo{
			Author: info.Author,
			Time: time.Now().Format("2006-01-02 15:04:05"),
			ServiceName: elem.Name}

		//Create implementation file
		err = GenerateServiceImplementation(workingPath, goPath, info, service, override)

		if err != nil {
			return err
		}

		//Create handlers
		for _, rpc := range elem.RPCs {
			fmt.Printf("  Name:%v\r\n  Doc:%v\r\n  Opt:%v\r\n  ReqType:%v\r\n  RspType:%v\r\n",
				rpc.Name, rpc.Documentation, rpc.Options, rpc.RequestType, rpc.ResponseType)


			rpc := &RpcInfo{
				Author: info.Author,
				Time: time.Now().Format("2006-01-02 15:04:05"),
				WorkingPathRelative: workingPathRelative,
				ServiceName: service.ServiceName,
				RpcName: rpc.Name,
				Request: rpc.RequestType.Name(),
				Response: rpc.ResponseType.Name(),}

				err = GenerateServiceHandler(workingPath, goPath, info, rpc, override)

				if err != nil {
					return err
				}
		}
	}

	return nil
}

//Generate health check service proto
func GenerateHealthCheckProto(workingPath, goPath string, info *TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "src", "proto", "health_check", "health_check.proto")
	filePathRelative := filepath.Join(".", "src", "proto", "health_check", "health_check.proto")

	//Check file is exist or not
	isExist, err := gofraUtils.CheckPathExists(filePath)

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
	cmd := exec.Command("protoc", "--go_out=plugins=grpc:.", filePathRelative)
	fmt.Println(filePath)

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

//Generate service implementation
func GenerateServiceImplementation(workingPath, goPath string, info *TemplateInfo, service *ServiceInfo, override bool) error {
	filePath := filepath.Join(workingPath, "src", "handler", service.ServiceName, service.ServiceName + ".go")

	//Check file is exist or not
	isExist, err := gofraUtils.CheckPathExists(filePath)

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
func GenerateServiceHandler(workingPath, goPath string, info *TemplateInfo, rpc *RpcInfo, override bool) error {
	filePath := filepath.Join(workingPath, "src", "handler", rpc.ServiceName, rpc.RpcName + ".go")

	//Check file is exist or not
	isExist, err := gofraUtils.CheckPathExists(filePath)

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
