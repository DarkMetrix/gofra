package docker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	gofraTemplate "github.com/DarkMetrix/gofra/internal/pkg/template"
	commonUtils "github.com/DarkMetrix/gofra/pkg/utils"
)

//Generate docker file
func GenerateDockerFile(workingPath string, info *gofraTemplate.TemplateInfo, override bool) error {
	filePath := filepath.Join(workingPath, "Dockerfile")

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
	dockerFileTemplate, err := template.New("docker_file").Parse(DockerFileTemplate)

	if err != nil {
		return err
	}

	dockerFileInfo := &DockerFileInfo{
		Author: info.Author,
		Project: info.Project,
	}

	file, err := os.OpenFile(filePath, os.O_RDWR | os.O_CREATE, 0755)

	if err != nil {
		return err
	}

	//Render template to file
	err = dockerFileTemplate.Execute(file, dockerFileInfo)

	if err != nil {
		return err
	}

	return nil
}

