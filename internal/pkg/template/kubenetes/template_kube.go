package kubenetes

// deployment config
type KubeDeploymentInfo struct {
	Namespace string
	Project string
	Version string

	ImagePath string
	ContainerPort string
}

var KubeDeploymentTemplate string = `
# API version
apiVersion: apps/v1

# resource type
kind: Deployment

# metadata of the deployment
metadata:
  name: {{.Project}}
  {{if ne .Namespace ""}}
  namespace: {{.Namespace}}
  {{end}}
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

// service config
type KubeServiceInfo struct {
	Namespace string
	Project string
	Type string				// grpc or http

	Port string
	TargetPort string
}

var kubeServiceTemplate string = `
# API version
apiVersion: v1

# resource type
kind: Service

# metadata of the service
metadata:
  name: {{.Project}} 
  {{if ne .Namespace ""}}
  namespace: {{.Namespace}}
  {{end}}

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

// config config
type KubeConfigmapInfo struct {
	Namespace string
	Project string
}

var kubeConfigmapTemplate string = `#!/bin/bash

case $1 in

"create")
  echo "configmap '{{.Project}}' creating..."
  kubectl create configmap {{.Project}} --namespace='{{.Namespace}}' --from-file=../configs
  ;;

"update")
  echo "configmap '{{.Project}}' deleting..."
  kubectl delete configmap {{.Project}} --namespace='{{.Namespace}}'

  echo "configmap '{{.Project}}' creating..."
  kubectl create configmap {{.Project}} --namespace='{{.Namespace}}' --from-file=../configs
  ;;

"delete")
  echo "configmap '{{.Project}}' deleting..."
  kubectl delete configmap {{.Project}} --namespace='{{.Namespace}}'
  ;;

"get")
  echo "configmap '{{.Project}}' getting..."
  kubectl describe configmap {{.Project}} --namespace='{{.Namespace}}'
  ;;

"")
  echo "no command found! command list [create, update, delete, get]"
  ;;

esac
`

