package kubernetes

import (
	"github.com/DarkMetrix/gofra/internal/pkg/option"
	"github.com/DarkMetrix/gofra/internal/pkg/templates"
	"golang.org/x/xerrors"
)

// KubeDeploymentInfo represents kubernetes deployment information
type KubeDeploymentInfo struct {
	Opts          *option.Options
	Namespace     string
	Project       string
	Version       string
	ImagePath     string
	ContainerPort string
}

// NewKubeDeploymentInfo returns a new KubeDeploymentInfo pointer
func NewKubeDeploymentInfo(opts ...option.Option) *KubeDeploymentInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &KubeDeploymentInfo{
		Opts:          newOpts,
		Project:       newOpts.Project,
		Namespace:     newOpts.Namespace,
		Version:       newOpts.Version,
		ImagePath:     newOpts.ImagePath,
		ContainerPort: newOpts.Port,
	}
}

// RenderFile render template and output to file
func (info *KubeDeploymentInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, info.Opts.Override, info.Opts.IgnoreExist,
		"template-k8s-deployment", KubeDeploymentTemplate, info); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

var KubeDeploymentTemplate string = `
# API version
apiVersion: apps/v1

# resource type
kind: Deployment

# metadata of the deployment
metadata:
  name: {{.Project}}
  namespace: {{.Namespace}}
  labels:
    app: {{.Project}}
    version: {{.Version}}

# specification
spec:
  # replica number to run
  replicas: 1
  selector:
    matchLabels:
      app: {{.Project}}
      version: {{.Version}}

  # using this template to create pod
  template:
    metadata:
      labels:
        app: {{.Project}}
        version: {{.Version}}
    spec:
      restartPolicy: Always
      containers:
        - name: {{.Project}}
          image: {{.ImagePath}}
          ports:
            - containerPort: {{.ContainerPort}}

      ##########################################
      # if you would like to use config map to keep the config.toml and log.config
      # below configuration could mount config files to /app/{{.Project}}/configs directory
      # config map YAML file could be generated using 'gofra kube configmap' command
      ##########################################
      #    volumeMounts:
      #      - name: configs
      #        mountPath: /app/{{.Project}}/configs/config.toml
      #        subPath: config.toml
      #      - name: configs
      #        mountPath: /app/{{.Project}}/configs/log.config
      #        subPath: log.config
      # volumes:
      #   - name: configs
      #     configMap:
      #       name: {{.Project}}

  ##########################################
  # More features and details, please visit 
  #     'https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/'
  # to get more information, this URL is for kubernetes v1.17 only
  ##########################################
`

// KubeServiceInfo represents kubernetes service information
type KubeServiceInfo struct {
	Opts       *option.Options
	Namespace  string
	Project    string
	Type       string // grpc or http
	Port       string
	TargetPort string
}

// NewKubeServiceInfo returns a new KubeServiceInfo pointer
func NewKubeServiceInfo(opts ...option.Option) *KubeServiceInfo {
	// init options
	newOpts := option.NewOptions(opts...)
	return &KubeServiceInfo{
		Opts:       newOpts,
		Project:    newOpts.Project,
		Namespace:  newOpts.Namespace,
		Type:       "grpc",
		Port:       newOpts.Port,
		TargetPort: newOpts.TargetPort,
	}
}

// RenderFile render template and output to file
func (info *KubeServiceInfo) RenderFile(outputPath string) error {
	if err := templates.RenderToFile(outputPath, info.Opts.Override, info.Opts.IgnoreExist,
		"template-k8s-service", KubeServiceTemplate, info); err != nil {
		return xerrors.Errorf("RenderToFile failed! error:%w", err)
	}
	return nil
}

var KubeServiceTemplate string = `
# API version
apiVersion: v1

# resource type
kind: Service

# metadata of the service
metadata:
  name: {{.Project}} 
  namespace: {{.Namespace}}

# specification
spec:
  type: ClusterIP
  ports:
    - name: {{.Type}}
      port: {{.Port}}
      protocol: TCP
      targetPort: {{.TargetPort}}
  selector:
    app: {{.Project}}

  ##########################################
  # More features and details, please visit 
  #     'https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/'
  # to get more information, this URL is for kubernetes v1.17 only
  ##########################################
`
