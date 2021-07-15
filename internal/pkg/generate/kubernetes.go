package generate

import (
	"path/filepath"

	"github.com/DarkMetrix/gofra/internal/pkg/directory"
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/templates/kubernetes"
	"golang.org/x/xerrors"
)

// InitKubeDeployment initializes the kubernetes deployment.yaml file
func InitKubeDeployment(layout directory.GRPCServiceLayout, opts ...option.Option) error {
	if err := kubernetes.NewKubeDeploymentInfo(opts...).RenderFile(
		filepath.Join(layout.GetDeployBasePath(), "deployment.yaml")); err != nil {
		return xerrors.Errorf("kubeDeploymentInfo.RenderFile failed: %v", err)
	}
	return nil
}

// InitKubeService initializes the kubernetes service.yaml file
func InitKubeService(layout directory.GRPCServiceLayout, opts ...option.Option) error {
	if err := kubernetes.NewKubeServiceInfo(opts...).RenderFile(
		filepath.Join(layout.GetDeployBasePath(), "service.yaml")); err != nil {
		return xerrors.Errorf("kubeServiceInfo.RenderFile failed: %v", err)
	}
	return nil
}
