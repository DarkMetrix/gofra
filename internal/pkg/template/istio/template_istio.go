package kubenetes

// istio virtual service config
type IstioVirtaulServiceInfo struct {
	Project string
	Version string

	Port string
}

var IstioVirtualServiceTemplate string = `
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{.Project}}

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

// istio destination rule config
type IstioDestinationRuleInfo struct {
	Project string
	Version string
}

var IstioDestinationRuleTemplate string = `
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: {{.Project}}
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

