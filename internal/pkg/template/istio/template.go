package istio

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"text/template"

	gofraTemplate "github.com/DarkMetrix/gofra/internal/pkg/template"
	commonUtils "github.com/DarkMetrix/gofra/pkg/utils"
)

//Generate istio virtual service yaml file
func GenerateIstioVirtaulServiceYAMLFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "istio", "virtual-service.yml")

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
		err := os.Remove(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	istioVirtualServiceTemplate, err := template.New("istio_virtual_service").Parse(IstioVirtualServiceTemplate)

	if err != nil {
		return err
	}

	_, port, err := net.SplitHostPort(info.Server.Addr)

	if err != nil {
		return err
	}

	istioVirtualServiceInfo := &IstioVirtaulServiceInfo {
		Project: info.Project,
		Version: info.Version,
		Port: port,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = istioVirtualServiceTemplate.Execute(file, istioVirtualServiceInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate istio destination rule yaml file
func GenerateIstioDestinationRuleYAMLFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "istio", "destination-rule.yml")

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
		err := os.Remove(filePath)

		if err != nil {
			return err
		}
	}

	//Parse template
	istioDestinationRuleTemplate, err := template.New("istio_destination_rule").Parse(IstioDestinationRuleTemplate)

	if err != nil {
		return err
	}

	istioDestinationRuleInfo := &IstioDestinationRuleInfo {
		Project: info.Project,
		Version: info.Version,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = istioDestinationRuleTemplate.Execute(file, istioDestinationRuleInfo)

	if err != nil {
		return err
	}

	return nil
}
