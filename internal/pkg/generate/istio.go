package generate

import (
	"path/filepath"

	"github.com/DarkMetrix/gofra/internal/pkg/directory"
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/templates/istio"
	"golang.org/x/xerrors"
)

// InitIstioVirtualService initializes the istio virtual_service.yaml file
func InitIstioVirtualService(layout directory.GRPCServiceLayout, opts ...option.Option) error {
	if err := istio.NewIstioVirtaulServiceInfo(opts...).RenderFile(
		filepath.Join(layout.GetDeployBasePath(), "virtual_service.yaml")); err != nil {
		return xerrors.Errorf("istioVirtualServiceInfo.RenderFile failed: %v", err)
	}
	return nil
}

// InitIstioDestinationRule initializes the istio destination_rule.yaml file
func InitIstioDestinationRule(layout directory.GRPCServiceLayout, opts ...option.Option) error {
	if err := istio.NewIstioDestinationRuleInfo(opts...).RenderFile(
		filepath.Join(layout.GetDeployBasePath(), "destination_rule.yaml")); err != nil {
		return xerrors.Errorf("istioDestinationRuleInfo.RenderFile failed: %v", err)
	}
	return nil
}
