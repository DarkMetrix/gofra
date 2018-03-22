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
    "client":
    {
        "pool":
        {
            "init_conns":5,
            "max_conns":10,
            "idle_time":30
        }
    },
    "monitor_package":
    {
        "package":"github.com/DarkMetrix/gofra/grpc-utils/monitor/statsd",
        "init_param":"\"127.0.0.1:8125\""
    },
    "tracing_package":
    {
        "package":"github.com/DarkMetrix/gofra/grpc-utils/tracing/zipkin",
        "init_param":"\"http://127.0.0.1:9411/api/v1/spans\", \"false\", \"{{.Addr}}\", \"{{.Project}}\""
    },
    "interceptor_package":
    {
        "monitor_package":"github.com/DarkMetrix/gofra/grpc-utils/interceptor/statsd_interceptor",
        "tracing_package":"github.com/DarkMetrix/gofra/grpc-utils/interceptor/zipkin_interceptor"
    }
}
`
