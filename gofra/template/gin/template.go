package gin

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"errors"
	"strings"
	"text/template"

	commonUtils "github.com/DarkMetrix/gofra/common/utils"
	gofraTemplate"github.com/DarkMetrix/gofra/gofra/template"
)

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

		return errors.New(fmt.Sprintf("File:%v already exists! this operation will overide it!", filePathRel))
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
func GenerateCommonFile(workingPath, goPath string, info *gofraTemplate.TemplateInfo, override bool) error {
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
func GenerateConfigFile(workingPath, goPath string, info *gofraTemplate.TemplateInfo, override bool) error {
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

//Generate config.toml
func GenerateConfigTomlFile(workingPath, goPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "conf", "config.toml")

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
func GenerateConfigLogFile(workingPath, goPath string, info *gofraTemplate.TemplateInfo, override bool) error {
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
func GenerateApplicationFile(workingPath, goPath string, info *gofraTemplate.TemplateInfo, override bool) error {
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
func GenerateMainFile(workingPath, goPath string, info *gofraTemplate.TemplateInfo, override bool) error {
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
func GenerateHealthCheckHttpHandler(workingPath, goPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	//Generate health check proto file
	err := GenerateHttpHandler(workingPath, goPath, info, "/health", override)

	if err != nil {
		return err
	}

	return nil
}

//Generate http service
func GenerateHttpHandler(workingPath, goPath string, info *gofraTemplate.TemplateInfo, uri string, override bool) error {
	//Generate http service
	//Create path
	handlerPath := filepath.Join(workingPath, "src", "http_handler")

	err := commonUtils.CreatePath(handlerPath, false)

	if err != nil {
		return err
	}

	//Add http handler
	uriPaths := strings.Split(uri, "/")

	if len(uriPaths) == 0 {
		return errors.New(fmt.Sprintf("URI is not valid! uri:%v", uri))
	}

	uriUpperPaths := make([]string, 0)

	for _, item := range uriPaths {
		if len(item) == 0 {
			continue
		}

		uriUpperPaths = append(uriUpperPaths, strings.ToUpper(item))
	}

	handlerName := strings.Join(uriUpperPaths, "_")

	handler := &HttpHandlerInfo{
		Author: info.Author,
		Time: time.Now().Format("2006-01-02 15:04:05"),
		HandlerName: handlerName,
		MonitorPackage: info.MonitorPackage.Package,
		TracingPackage: info.TracingPackage.Package,
	}

	err = GenerateHttpHandlerFile(workingPath, goPath, info, handler, override)

	if err != nil {
		return err
	}

	//Add http handler to engine in application
	err = AddHttpHandlerToApplicationFile(workingPath, goPath, info, handler, uri)

	if err != nil {
		return err
	}

	return nil
}

//Add http handler import to application file
func AddHttpHandlerToApplicationFile(workingPath, goPath string, info *gofraTemplate.TemplateInfo, handler *HttpHandlerInfo, uri string) error {
	applicationFilePath := filepath.Join(workingPath, "src", "application", "application.go")

	applicationContent, err := ioutil.ReadFile(applicationFilePath)

	if err != nil {
		return err
	}

	httpHandler := fmt.Sprintf(`group.POST("%v", httpHandler.%v)`, uri, handler.HandlerName)
	httpHandlerStub := fmt.Sprintf("%v\r\n	/*@REGISTER_HTTP_STUB*/", httpHandler)

	if strings.Contains(string(applicationContent), httpHandler) {
		return nil
	}

	applicationContent = []byte(strings.Replace(string(applicationContent), "/*@REGISTER_HTTP_STUB*/", httpHandlerStub, 1))

	err = ioutil.WriteFile(applicationFilePath, applicationContent, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

//Generate http handler file
func GenerateHttpHandlerFile(workingPath, goPath string, info *gofraTemplate.TemplateInfo, handler *HttpHandlerInfo, override bool) error {
	filePath := filepath.Join(workingPath, "src", "http_handler", handler.HandlerName + ".go")

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
	serviceRpcTemplate, err := template.New("http_handler").Parse(HttpHandlerTemplate)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = serviceRpcTemplate.Execute(file, handler)

	if err != nil {
		return err
	}

	return nil
}

//Generate test client
func GenerateTestClient(workingPath, goPath string, info *gofraTemplate.TemplateInfo, override bool) error {
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
