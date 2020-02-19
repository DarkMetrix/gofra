package template

//Server config
type ServerInfo struct {
	Addr string `json:"addr"`
}

//Monitor package
type MonitorPackageInfo struct {
	Package string `json:"package"`
	InitParam string `json:"init_param"`
}

//Tracing package
type TracingPackageInfo struct {
	Package string `json:"package"`
	InitParam string `json:"init_param"`
}

//Interceptor package
type InterceptorPackageInfo struct {
	MonitorPackage string `json:"monitor_package"`
	TracingPackage string `json:"tracing_package"`
}

//Template info
type TemplateInfo struct {
	Author string `json:"author"`
	Project string `json:"project"`
	Version string `json:"version"`
	Type string `json:"type"`

	Server ServerInfo `json:"server"`

	MonitorPackage MonitorPackageInfo `json:"monitor_package"`
	TracingPackage TracingPackageInfo `json:"tracing_package"`

	InterceptorPackage InterceptorPackageInfo `json:"interceptor_package"`
}

type JsonInfo struct {
	Author string
	Project string
	Type string

	Addr string
}

var JsonTemplate string = `
{
    "author":"{{.Author}}",
    "project":"{{.Project}}",
    "version":"0.0.1",
    "type":"{{.Type}}",
    "server":
    {
        "addr":"{{.Addr}}"
    },
    "monitor_package":
    {
        "package":"github.com/DarkMetrix/gofra/pkg/monitor/statsd",
        "init_param":"\"127.0.0.1:8125\", \"{{.Project}}\""
    },
    "tracing_package":
    {
        "package":"github.com/DarkMetrix/gofra/pkg/tracing/jaeger",
        "init_param":"\"127.0.0.1:6831\", \"{{.Project}}\""
    },
    "interceptor_package":
    {
        "monitor_package":"github.com/DarkMetrix/gofra/pkg/grpc-utils/interceptor/statsd_interceptor",
        "tracing_package":"github.com/DarkMetrix/gofra/pkg/grpc-utils/interceptor/opentracing_interceptor"
    }
}
`
