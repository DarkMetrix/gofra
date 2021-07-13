package docker

import (
	"strings"

	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/templates"
	"golang.org/x/xerrors"
)

// DockerFileInfo represents docker file information
type DockerFileInfo struct {
	Opts    *option.Options
	Project string
	Labels  string
}

// NewDockerFileInfo returns a new DockerFileInfo pointer
func NewDockerFileInfo(opts ...option.Option) *DockerFileInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &DockerFileInfo{
		Opts:    newOpts,
		Project: newOpts.Project,
		Labels:  strings.Join(newOpts.Labels, " "),
	}
}

// RenderFile render template and output to file
func (info *DockerFileInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, info.Opts.Override, info.Opts.IgnoreExist,
		"template-docker-file", DockerFileTemplate, info); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

var DockerFileTemplate string = `
FROM centos:latest
LABEL {{.Labels}}

COPY ./build /application/bin
COPY ./configs /application/configs

WORKDIR /application/bin

ENTRYPOINT ["./{{.Project}}"]
`
