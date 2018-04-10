# DarkMetrix Gofra

>  Don't reinvent the wheel, just realign it.

It's famous in IT industry and I'm quiet a believer of that. 

Gofra is a framework and tool set helps to easily generate gRPC service in Go.

All the tools and libraries are open sourced, you could use it out of box.



## Installation

First [Google's Protocol Buffers](https://developers.google.com/protocol-buffers/) Version 3 is required to be installed.

Second [golang/protobuf](https://github.com/golang/protobuf) is also needed to be installed to generate gRPC service from *.proto file.

Then install Gofra

```bash
$ go get -u github.com/DarkMetrix/gofra/gofra
```



## Guide

- [Creating a service template](#creating-a-service-template)
- [Service Generation](#service-generation)
- [Add or Update Services](#add-or-update-services)
- [Implement RPC Methods](#implement-rpc-methods)
- [Compile and Run](#compile-and-run)
- [Test using Health Check](#test-using-health-check)



### Creating a Service Template

**All the commands should be used in the subdirectory of $GOPATH, that'll be forced by Gofra and it's really a good habit.**

First initialize default template file as below.

```bash
$ gofra template init
```

 A **template.json.default** file will be generated in the current directory which looks like this.

```json
{
    "author":"Author Name",
    "project":"Project Name",
    "version":"0.0.1",
    "server":
    {
        "addr":"localhost:58888"
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
        "init_param":"\"127.0.0.1:8125\", \"Project Name\""
    },
    "tracing_package":
    {
        "package":"github.com/DarkMetrix/gofra/grpc-utils/tracing/zipkin",
        "init_param":"\"http://127.0.0.1:9411/api/v1/spans\", \"false\", \"localhost:58888\", \"Project Name\""
    },
    "interceptor_package":
    {
        "monitor_package":"github.com/DarkMetrix/gofra/grpc-utils/interceptor/statsd_interceptor",
        "tracing_package":"github.com/DarkMetrix/gofra/grpc-utils/interceptor/zipkin_interceptor"
    }
}
```

| Json object                         | Value                                                        | Remark       |
| ----------------------------------- | ------------------------------------------------------------ | ------------ |
| author                              | Author's name                                                | User defined |
| project                             | Project's name                                               | User defined |
| version                             | Version                                                      | User defined |
| server.addr                         | Servicehe 's ip and port                                     | User defined |
| client.pool.init_conns              | Client's connection pool's initial connection number         | User defined |
| client.pool.max_conns               | Client's connection pool's max connection number             | User defined |
| client.pool.idle_time               | Client's connection pool's connection idle time in seconds   | User defined |
| monitor_package.package             | Monitor package import path used in the service(by default statsd is used as the monitor backend), you could write your own monitor package as long as you implement Init & Increment interfaces | Pre defined  |
| monitor_package.init_param          | Monitor params used by Init method(it's the statsd address and project name used in your environment) | Pre defined  |
| Tracing_package.package             | Tracing package import path used in the service(by default zipkin is used as the monitor backend), you could write your own tracing package as long as you implement Init interface | Pre defined  |
| Tracing_package.init_param          | Tracing params used by Init method(it's the zipkin address, debug flag, service address and project name used in your environment) | Pre defined  |
| Interceptor_package.monitor_package | Monitor gRPC interceptor(by default statsd is used as the monitor backend) | Pre defined  |
| Interceptor_package.tracing_package | Tracing gPRC interceptor(by default zipkin is used as the tracing system) | Pre defined  |

Then modify the author's name and project's name to 'tester' and 'test_gofra' and change the file name from template.json.default to template.json

### Service Generation

```bash
$gofra init --path=./template.json
```

All files needed are generated.

```bash
$tree

.
├── bin
├── conf
│   ├── config.json
│   ├── log.config
│   └── naming.json
├── log
├── src
│   ├── application
│   │   └── application.go
│   ├── common
│   │   └── common.go
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   └── HealthCheckService
│   │       ├── HealthCheck.go
│   │       └── HealthCheckService.go
│   ├── main.go
│   └── proto
│       ├── health_check
│       │   ├── health_check.pb.go
│       │   └── health_check.proto
│       └── user
├── template.json
└── test
    └── main.go

13 directories, 12 files
```



### Add or Update Services

By default a health check service is generated.

You could using test/main.go to send health check request to test if the service is working fine.

And you could add your own service, e.g.:

#### Write user.proto

```protobuf
syntax = "proto3";

package user;

// The greeting service definition.
service UserService {
  // Sends a greeting
  rpc AddUser (AddUserRequest) returns (AddUserResponse) {}
}

message AddUserRequest {
    string name = 1;
}

message AddUserResponse {
    string name = 1;
}
```

#### Add Service

```bash
$gofra service add --path=./user.proto
```

Then a directory UserService & two files AddUser.go and UserService.go are generated.

```bash
$ tree

.
├── bin
├── clean.sh
├── conf
│   ├── config.json
│   ├── log.config
│   └── naming.json
├── init.sh
├── log
├── src
│   ├── application
│   │   └── application.go
│   ├── common
│   │   └── common.go
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   ├── HealthCheckService
│   │   │   ├── HealthCheck.go
│   │   │   └── HealthCheckService.go
│   │   └── UserService
│   │       ├── AddUser.go
│   │       └── UserService.go
│   ├── main.go
│   └── proto
│       ├── health_check
│       │   ├── health_check.pb.go
│       │   └── health_check.proto
│       └── user
│           ├── user.pb.go
│           └── user.proto
├── template.json
├── test
│   └── main.go
└── user.proto
```

#### Update Service

You could update the service if the pb file has updates, such as adding a new method to the UserService.

The only difference between add and update is update won't override the file already generated.

Using add command, a —overried=true flag will help to override all the files generated about the pb file.



### Implement RPC Methods

After add or update service, all the thing left is to implement your own business logic.

All basic logging, metrics, tracing, panic recovery have been added already.



#### UserService/UserService.go

You do not need to do anything about this file, just leave it.

```bash
$cat src/handler/UserService/UserService.go

/**********************************
 * Author : tester
 * Time : 2018-03-26 23:04:23
 **********************************/

package UserService

// Implement UserService interface
type UserServiceImpl struct{}
```



#### UserService/AddUser.go

Implement your own business logic.

```bash
$cat src/handler/UserService/AddUser.go

/**********************************
 * Author : techieliu
 * Time : 2018-04-01 02:26:20
 **********************************/

package UserService

import (
        "context"

        //Log package
        //log "github.com/cihub/seelog"

        //Monitor package
        //monitor "github.com/DarkMetrix/gofra/grpc-utils/monitor/statsd"

        //Tracing package
        //tracing "github.com/DarkMetrix/gofra/grpc-utils/tracing/zipkin"

        pb "github.com/DarkMetrix/gofra/tmp/demo/src/proto/user"
)

func (this UserServiceImpl) AddUser (ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
        //Log Example:traceid must be logged
        //log.Infof("AddUser begin, traceid=%v, req=%v", tracing.GetTracingId(ctx), req)

        resp := new(pb.AddUserResponse)

        return resp, nil
}
```

#### 

### Compile and Run

#### Complie

```bash
$cd src
$go build -o test_gofra
```



#### Run

```bash
$./test_gofra

====== Server [test_gofra] Start ======
Listen on port [:58888]
```



### Test using Health Check

#### Complie test

By default the health check client request is already generated.

```bash
$cd test
$go build
```



#### Run test

Run to test health check.

```bash
$./test
```



#### Client output

```bash
[DEBUG][2018-04-01T02:26:48.890452][test_gofra_test][seelog_interceptor.go:26][GofraClientInterceptorFunc] => context=context.Background.WithValue(metadata.mdOutgoingKey{}, metadata.MD{"x-b3-traceid":[]string{"7d123708aa5946c52efc7ed01f8fde82"}, "x-b3-spanid":[]string{"5fd7cb9d44b4ed72"}, "x-b3-sampled":[]string{"1"}, "x-b3-flags":[]string{"0"}}).WithValue(opentracing.contextKey{}, &zipkintracer.spanImpl{tracer:(*zipkintracer.tracerImpl)(0xc42018c400), event:(func(zipkintracer.SpanEvent))(nil), observer:otobserver.SpanObserver(nil), Mutex:sync.Mutex{state:0, sema:0x0}, raw:zipkintracer.RawSpan{Context:zipkintracer.SpanContext{TraceID:types.TraceID{High:0x7d123708aa5946c5, Low:0x2efc7ed01f8fde82}, SpanID:0x5fd7cb9d44b4ed72, Sampled:true, Baggage:map[string]string(nil), ParentSpanID:(*uint64)(nil), Flags:0x8, Owner:true}, Operation:"/common.health.check.HealthCheckService/HealthCheck", Start:time.Time{wall:0xbea8129a34dfcbca, ext:6018394, loc:(*time.Location)(0xcdb0a0)}, Duration:3396107, Tags:opentracing.Tags{"component":"gRPC"}, Logs:[]opentracing.LogRecord(nil)}, numDroppedLogs:0, Endpoint:(*zipkincore.Endpoint)(nil)}), req=message:"ping" , invoke success!!! reply:
```



#### Service output

```
[DEBUG][2018-04-01T02:26:48.888806][test_gofra][seelog_interceptor.go:42][GofraServerInterceptorFunc] => context=context.Background.WithCancel.WithCancel.WithValue(peer.peerKey{}, &peer.Peer{Addr:(*net.TCPAddr)(0xc42016bbf0), AuthInfo:credentials.AuthInfo(nil)}).WithValue(transport.streamKey{}, <stream: 0xc42022e280, /common.health.check.HealthCheckService/HealthCheck>).WithValue(metadata.mdIncomingKey{}, metadata.MD{"user-agent":[]string{"grpc-go/1.8.2"}, "x-b3-traceid":[]string{"7d123708aa5946c52efc7ed01f8fde82"}, "x-b3-spanid":[]string{"5fd7cb9d44b4ed72"}, "x-b3-sampled":[]string{"1"}, "x-b3-flags":[]string{"0"}, ":authority":[]string{"172.16.101.128:58888"}}).WithValue(opentracing.contextKey{}, &zipkintracer.spanImpl{tracer:(*zipkintracer.tracerImpl)(0xc4200f0d00), event:(func(zipkintracer.SpanEvent))(nil), observer:otobserver.SpanObserver(nil), Mutex:sync.Mutex{state:0, sema:0x0}, raw:zipkintracer.RawSpan{Context:zipkintracer.SpanContext{TraceID:types.TraceID{High:0x7d123708aa5946c5, Low:0x2efc7ed01f8fde82}, SpanID:0x5fd7cb9d44b4ed72, Sampled:true, Baggage:map[string]string(nil), ParentSpanID:(*uint64)(nil), Flags:0x2, Owner:false}, Operation:"/common.health.check.HealthCheckService/HealthCheck", Start:time.Time{wall:0xbea8129a34f98f74, ext:3707775087, loc:(*time.Location)(0xdfb0a0)}, Duration:56535, Tags:opentracing.Tags{"component":"gRPC"}, Logs:[]opentracing.LogRecord(nil)}, numDroppedLogs:0, Endpoint:(*zipkincore.Endpoint)(nil)}), req=message:"ping" , handle success!!! reply:
```



## License

### Gofra

[MIT license](https://github.com/DarkMetrix/gofra/blob/master/LICENSE)

### Dependencies

- github.com/cihub/seelog [BSD License](https://github.com/cihub/seelog/blob/master/LICENSE.txt)
- github.com/spf13/viper [MIT License](https://github.com/spf13/viper/blob/master/LICENSE)
- github.com/spf13/cobra [Apache 2.0 License](https://github.com/spf13/cobra/blob/master/LICENSE.txt)
- github.com/grpc-ecosystem/go-grpc-middleware [Apache 2.0 License](https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/LICENSE)
- github.com/tallstoat/pbparser [MIT License](https://github.com/tallstoat/pbparser/blob/master/LICENSE)
- github.com/alexcesaro/statsd [MIT License](https://github.com/alexcesaro/statsd/blob/master/LICENSE)
- github.com/processout/grpc-go-pool [MIT License](https://github.com/processout/grpc-go-pool/blob/master/LICENSE)
- github.com/opentracing/opentracing-go [MIT License](https://github.com/opentracing/opentracing-go/blob/master/LICENSE)
- github.com/openzipkin/zipkin-go-opentracing [MIT License](https://github.com/openzipkin/zipkin-go-opentracing/blob/master/LICENSE)

