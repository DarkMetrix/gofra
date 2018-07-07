package template

type JsonInfo struct {
	Author string
	Project string

	Addr string
}

var JsonTemplate string = `
{
    "author":"{{.Author}}",
    "project":"{{.Project}}",
    "version":"0.0.1",
    "server":
    {
        "addr":"{{.Addr}}"
    },
    "monitor_package":
    {
        "package":"github.com/DarkMetrix/gofra/common/monitor/statsd",
        "init_param":"\"127.0.0.1:8125\", \"{{.Project}}\""
    },
    "tracing_package":
    {
        "package":"github.com/DarkMetrix/gofra/common/tracing/jaeger",
        "init_param":"\"127.0.0.1:6831\", \"{{.Project}}\""
    },
    "interceptor_package":
    {
        "monitor_package":"github.com/DarkMetrix/gofra/grpc-utils/interceptor/statsd_interceptor",
        "tracing_package":"github.com/DarkMetrix/gofra/grpc-utils/interceptor/opentracing_interceptor"
    }
}
`
