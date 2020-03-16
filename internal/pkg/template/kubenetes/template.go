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

// generate kubernetes deployment yaml file
func GenerateKubeDeploymentYAMLFile(workingPath, imagePath, namespace string, info *gofraTemplate.TemplateInfo, override bool) error {
	var filePath string

	if namespace == "" {
		filePath = filepath.Join(workingPath, "kubernetes", "deployment.yml")
	} else {
		filePath = filepath.Join(workingPath, "kubernetes", fmt.Sprintf("deployment-%v.yml", namespace))
	}

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
		err := os.Remove(filePath)

		if err != nil {
			return err
		}
	}

	// parse template
	kubeDeploymentTemplate, err := template.New("kube_deployment").Parse(KubeDeploymentTemplate)

	if err != nil {
		return err
	}

	_, port, err := net.SplitHostPort(info.Server.Addr)

	if err != nil {
		return err
	}

	kubeDeploymentInfo := &KubeDeploymentInfo {
		Namespace: namespace,
		Project: info.Project,
		Version: info.Version,
		ImagePath: imagePath,
		ContainerPort: port,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	// render template to file
	err = kubeDeploymentTemplate.Execute(file, kubeDeploymentInfo)

	if err != nil {
		return err
	}

	return nil
}

// generate kubernetes service yaml file
func GenerateKubeServiceYAMLFile(workingPath, namespace string, info *gofraTemplate.TemplateInfo, override bool) error {
	var filePath string

	if namespace == "" {
		filePath = filepath.Join(workingPath, "kubernetes", "service.yml")
	} else {
		filePath = filepath.Join(workingPath, "kubernetes", fmt.Sprintf("service-%v.yml", namespace))
	}

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
		err := os.Remove(filePath)

		if err != nil {
			return err
		}
	}

	// parse template
	kubeServiceTemplate , err := template.New("kube_service").Parse(kubeServiceTemplate)

	if err != nil {
		return err
	}

	_, port, err := net.SplitHostPort(info.Server.Addr)

	if err != nil {
		return err
	}

	kubeServiceInfo := &KubeServiceInfo {
		Namespace: namespace,
		Project: info.Project,
		Type: info.Type,
		Port: port,
		TargetPort: port,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	// render template to file
	err = kubeServiceTemplate.Execute(file, kubeServiceInfo)

	if err != nil {
		return err
	}

	return nil
}

// generate kubernetes configmap yaml file
func GenerateKubeConfigmapYAMLFile(workingPath, namespace string, info *gofraTemplate.TemplateInfo, override bool) error {
	var filePath string

	if namespace == "" {
		filePath = filepath.Join(workingPath, "kubernetes", "configmap.sh")
	} else {
		filePath = filepath.Join(workingPath, "kubernetes", fmt.Sprintf("configmap-%v.sh", namespace))
	}

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
		err := os.Remove(filePath)

		if err != nil {
			return err
		}
	}

	// parse template
	kubeConfigmapTemplate , err := template.New("kube_configmap").Parse(kubeConfigmapTemplate)

	if err != nil {
		return err
	}

	kubeConfigmapInfo := &KubeConfigmapInfo {
		Namespace: namespace,
		Project: info.Project,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	// render template to file
	err = kubeConfigmapTemplate.Execute(file, kubeConfigmapInfo)

	if err != nil {
		return err
	}

	return nil
}
