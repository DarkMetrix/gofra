package istio

import (
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/templates"
	"golang.org/x/xerrors"
)

// IstioVirtaulServiceInfo represents istio virtual-service information
type IstioVirtaulServiceInfo struct {
	Opts      *option.Options
	Project   string
	Namespace string
	Version   string
	Port      string
}

// NewIstioVirtaulServiceInfo returns a new IstioVirtaulServiceInfo pointer
func NewIstioVirtaulServiceInfo(opts ...option.Option) *IstioVirtaulServiceInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &IstioVirtaulServiceInfo{
		Opts:      newOpts,
		Project:   newOpts.Project,
		Namespace: newOpts.Namespace,
		Version:   newOpts.Version,
		Port:      newOpts.Port,
	}
}

// RenderFile render template and output to file
func (info *IstioVirtaulServiceInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, info.Opts.Override, info.Opts.IgnoreExist,
		"template-istio-virtual-service", IstioVirtualServiceTemplate, info); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

var IstioVirtualServiceTemplate string = `
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{.Project}}
  {{if ne .Namespace ""}}
  namespace: {{.Namespace}}
  {{end}}

spec:
  hosts:
    - {{.Project}}            # hosts will be interpreted as ${hosts}.${k8s-namespace}.svc.cluster.local in k8s environment
  http:
    - route:
      # destination works with host and subset defined in destination rule
      # if using destination rule to manage traffic policy, then you could use subset to manage route rule
      - destination:
          host: {{.Project}} # host will be interpreted as ${hosts}.${k8s-namespace}.svc.cluster.local in k8s environment
          subset: {{.Version}}
      
      timeout: 2s
      
      # destination works with host and port
      #- destination:
      #    host: {{.Project}}  # host will be interpreted as ${hosts}.${k8s-namespace}.svc.cluster.local in k8s environment
      #	  port: 
      #	    number: {{.Port}} # port defined in kubernetes service
      
      # retry setting, use it as business needs
      #retries:
      #  attempts: 2
      #  perTryTimeout: 2s
      
      # There's a lot of features, such as traffic mirroring, weighted routing etc.

  ##########################################
  # More features and details, such as:
  #     CORS Policy
  #     HTTP Fault Injection
  #     HTTP Match, Redirect, Retry, Rewrite, Route, Headers
  #     TPC Route
  #     Route Destination
  #     TLS setting
  #     ...
  # Please visit 
  #     'https://istio.io/docs/reference/config/networking/virtual-service/'
  # to get more information
  ##########################################
`

// IstioDestinationRuleInfo represents istio destination rule information
type IstioDestinationRuleInfo struct {
	Opts      *option.Options
	Namespace string
	Project   string
	Version   string
}

// NewIstioDestinationRuleInfo returns a new IstioDestinationRuleInfo pointer
func NewIstioDestinationRuleInfo(opts ...option.Option) *IstioDestinationRuleInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &IstioDestinationRuleInfo{
		Opts:      newOpts,
		Project:   newOpts.Project,
		Namespace: newOpts.Namespace,
		Version:   newOpts.Version,
	}
}

// RenderFile render template and output to file
func (info *IstioDestinationRuleInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, info.Opts.Override, info.Opts.IgnoreExist,
		"template-istio-destination-rule", IstioDestinationRuleTemplate, info); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

var IstioDestinationRuleTemplate string = `
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: {{.Project}}
  {{if ne .Namespace ""}}
  namespace: {{.Namespace}}
  {{end}}

spec:
  host: {{.Project}} # host will be interpreted as ${hosts}.${k8s-namespace}.svc.cluster.local in k8s environment

  # traffic policy setting
  trafficPolicy:
    loadBalancer:
      simple: ROUND_ROBIN

  # subsets definition
  subsets:
    - name: {{.Version}}        # subset name, it could be used in virtual service to define routing rule
      labels:
        version: {{.Version}}   # all traffic from subset v1 will be routed to application with v1 label

  ##########################################
  # More features and details, such as:
  #     Connection Pool
  #     Load Balancer
  #     Outliner Detection
  #     TLS setting
  #     ...
  # Please visit 
  #     'https://istio.io/docs/reference/config/networking/destination-rule/' 
  # to get more information
  ##########################################
`
