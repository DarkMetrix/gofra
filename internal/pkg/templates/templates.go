package templates

import (
	"os"
	"text/template"

	"github.com/DarkMetrix/gofra/pkg/utils"
	"golang.org/x/xerrors"
)

// FileTemplate defines the interface every template should implement
type FileTemplate interface {
	RenderFile(filePath string) error
}

// RenderToFile render templates to file
func RenderToFile(filePath string, override, ignoreExist bool,
	templateName, templateContent string, content interface{}) error {
	// check file is exist or not
	isExist, err := utils.CheckPathExists(filePath)
	if err != nil {
		return xerrors.Errorf("utils.CheckPathExists failed! error:%w", err)
	}

	if isExist && !override {
		if ignoreExist {
			return nil
		}
		return xerrors.Errorf("File already exists! this operation will override it! file path:%v", filePath)
	}

	if isExist && override {
		if err := os.RemoveAll(filePath); err != nil {
			return xerrors.Errorf("os.RemoveAll failed! file path:%v, error:%w", filePath, err)
		}
	}

	// parse templates
	templateToRender, err := template.New(templateName).Parse(templateContent)
	if err != nil {
		return xerrors.Errorf("Unable to parse templates! error:%w", err)
	}

	// execute templates
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return xerrors.Errorf("os.OpenFile failed! error:%w", err)
	}

	if err := templateToRender.Execute(file, content); err != nil {
		return xerrors.Errorf("Unable to execute templates! error:%w", err)
	}
	return nil
}
