package kubenetes

//Deployment config
type KubeDeploymentInfo struct {
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
  labels:
    app: {{.Version}}
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
      containers:
        - name: {{.Version}}
          image: {{.ImagePath}}
          ports:
            - containerPort: {{.ContainerPort}}
      restartPolicy: Always
`

//Service config
type KubeServiceInfo struct {
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
`

