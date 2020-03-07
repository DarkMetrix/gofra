package kubenetes

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

//Generate kubernetes deployment yaml file
func GenerateKubeDeploymentYAMLFile(workingPath, imagePath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "kubernetes", "deployment.yml")

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
	kubeDeploymentTemplate, err := template.New("kube_deployment").Parse(KubeDeploymentTemplate)

	if err != nil {
		return err
	}

	_, port, err := net.SplitHostPort(info.Server.Addr)

	if err != nil {
		return err
	}

	kubeDeploymentInfo := &KubeDeploymentInfo {
		Project: info.Project,
		Version: info.Version,
		ImagePath: imagePath,
		ContainerPort: port,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = kubeDeploymentTemplate.Execute(file, kubeDeploymentInfo)

	if err != nil {
		return err
	}

	return nil
}

//Generate kubernetes service yaml file
func GenerateKubeServiceYAMLFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "kubernetes", "service.yml")

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
	kubeServiceTemplate , err := template.New("kube_service").Parse(kubeServiceTemplate)

	if err != nil {
		return err
	}

	_, port, err := net.SplitHostPort(info.Server.Addr)

	if err != nil {
		return err
	}

	kubeServiceInfo := &KubeServiceInfo {
		Project: info.Project,
		Type: info.Type,
		Port: port,
		TargetPort: port,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = kubeServiceTemplate.Execute(file, kubeServiceInfo)

	if err != nil {
		return err
	}

	return nil
}
