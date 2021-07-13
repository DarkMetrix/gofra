package generate

import (
	"path/filepath"

	"github.com/DarkMetrix/gofra/internal/pkg/directory"
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/templates/docker"
	"golang.org/x/xerrors"
)

// InitDockerFile initializes the Dockerfile used to generate docker image
func InitDockerFile(layout directory.GRPCServiceLayout, opts ...option.Option) error {
	if err := docker.NewDockerFileInfo(opts...).RenderFile(
		filepath.Join(layout.GetOutputPath(), "Dockerfile")); err != nil {
		return xerrors.Errorf("dockerFileInfo.RenderFile failed: %v", err)
	}
	return nil
}
