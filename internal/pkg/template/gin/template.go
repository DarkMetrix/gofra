package gin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	gofraTemplate "github.com/DarkMetrix/gofra/internal/pkg/template"
	commonUtils "github.com/DarkMetrix/gofra/pkg/utils"
)

// generate template.json
func GenerateTemplateJsonFile(workingPath string, override bool, jsonInfo gofraTemplate.JsonInfo) error {
	filePath := filepath.Join(workingPath, "template.json")

	// check file is exist or not
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

	// parse template
	jsonTemplate, err := template.New("template_json").Parse(gofraTemplate.JsonTemplate)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	// render template to file
	err = jsonTemplate.Execute(file, jsonInfo)

	if err != nil {
		return err
	}

	return nil
}

// generate common.go
func GenerateCommonFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "internal", "pkg", "common", "common.go")

	// check file is exist or not
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

	// parse template
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

	// render template to file
	err = commonTemplate.Execute(file, commonInfo)

	if err != nil {
		return err
	}

	return nil
}

// generate config.go
func GenerateConfigFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "internal", "pkg", "config", "config.go")

	// check file is exist or not
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

	// parse template
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

	// render template to file
	err = configTemplate.Execute(file, configInfo)

	if err != nil {
		return err
	}

	return nil
}

// generate config.toml
func GenerateConfigTomlFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "configs", "config.toml")

	// check file is exist or not
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

	// parse template
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

	// render template to file
	err = configTomlTemplate.Execute(file, configTomlInfo)

	if err != nil {
		return err
	}

	return nil
}

// generate log.config
func GenerateConfigLogFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "configs", "log.config")

	// check file is exist or not
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

	// parse template
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

	// render template to file
	err = configLogTemplate.Execute(file, configLogInfo)

	if err != nil {
		return err
	}

	return nil
}

// generate application.go
func GenerateApplicationFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "internal", "app", "application.go")

	// check file is exist or not
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

	workingPathRelative := filepath.Base(workingPath)

	// parse template
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

	// render template to file
	err = applicationTemplate.Execute(file, applicationInfo)

	if err != nil {
		return err
	}

	return nil
}

// generate main.go
func GenerateMainFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "cmd", "main.go")

	// check file is exist or not
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
func GenerateHealthCheckHttpHandler(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	//Generate health check proto file
	err := GenerateHttpHandler(workingPath, info, "/health", "GET", override)

	if err != nil {
		return err
	}

	return nil
}

//Generate http service
func GenerateHttpHandler(workingPath string, info *gofraTemplate.TemplateInfo, uri, method string, override bool) error {
	//Generate http service
	//Create path
	handlerPath := filepath.Join(workingPath, "internal", "http_handler")

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
	uriUpperPaths = append(uriUpperPaths, strings.ToUpper(method))

	for _, item := range uriPaths {
		if len(item) == 0 {
			continue
		}

		uriUpperPaths = append(uriUpperPaths, strings.ToUpper(item))
	}

	handlerName := strings.Join(uriUpperPaths, "_")
	handlerName = strings.Replace(handlerName, ":", "", -1)
	handlerName = strings.Replace(handlerName, "*", "", -1)

	//Mkdir
	fmt.Printf("\r\nMake dir ......")
	httpSpecPath := filepath.Join(workingPath, "api", "http_spec", strings.ToLower(handlerName))

	commonUtils.CreatePath(httpSpecPath, override)

	if err != nil {
		fmt.Printf(" failed! \r\nerror:%v\r\n", err.Error())
		return err
	} else {
		fmt.Printf(" success! \r\n")
	}

	handler := &HttpHandlerInfo{
		Author: info.Author,
		Time: time.Now().Format("2006-01-02 15:04:05"),
		HandlerName: handlerName,
		URI: uri,
		Method: method,
		MonitorPackage: info.MonitorPackage.Package,
		TracingPackage: info.TracingPackage.Package,
	}

	err = GenerateHttpHandlerFile(workingPath, info, handler, override)

	if err != nil {
		return err
	}

	//Add http handler to engine in application
	err = AddHttpHandlerToApplicationFile(workingPath, info, handler, uri, method)

	if err != nil {
		return err
	}

	return nil
}

//Add http handler import to application file
func AddHttpHandlerToApplicationFile(workingPath string, info *gofraTemplate.TemplateInfo, handler *HttpHandlerInfo, uri, method string) error {
	applicationFilePath := filepath.Join(workingPath, "internal", "app", "application.go")

	applicationContent, err := ioutil.ReadFile(applicationFilePath)

	if err != nil {
		return err
	}

	httpHandler := ""
	httpHandlerStub := ""

	switch method {
	case "GET":
		httpHandler = fmt.Sprintf(`group.GET("%v", httpHandler.%v)`, uri, handler.HandlerName)
	case "POST":
		httpHandler = fmt.Sprintf(`group.POST("%v", httpHandler.%v)`, uri, handler.HandlerName)
	case "PUT":
		httpHandler = fmt.Sprintf(`group.PUT("%v", httpHandler.%v)`, uri, handler.HandlerName)
	case "PATCH":
		httpHandler = fmt.Sprintf(`group.PATCH("%v", httpHandler.%v)`, uri, handler.HandlerName)
	case "DELETE":
		httpHandler = fmt.Sprintf(`group.DELETE("%v", httpHandler.%v)`, uri, handler.HandlerName)
	case "OPTIONS":
		httpHandler = fmt.Sprintf(`group.OPTIONS("%v", httpHandler.%v)`, uri, handler.HandlerName)
	default:
		return fmt.Errorf("Invalid http method! method:%v", method)
	}

	httpHandlerStub = fmt.Sprintf("%v\r\n	/*@REGISTER_HTTP_STUB*/", httpHandler)

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
func GenerateHttpHandlerFile(workingPath string, info *gofraTemplate.TemplateInfo, handler *HttpHandlerInfo, override bool) error {
	filePath := filepath.Join(workingPath, "internal", "http_handler", strings.ToLower(handler.HandlerName) + ".go")

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
	httpHandlerTemplate, err := template.New("http_handler").Parse(HttpHandlerTemplate)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = httpHandlerTemplate.Execute(file, handler)

	if err != nil {
		return err
	}

	return nil
}
