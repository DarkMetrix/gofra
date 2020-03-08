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

- [Creating a Service Template](#creating-a-service-template)

### Creating a Service Template

**Now gofra used go module, so you need to update your Go version up to at least v1.11**

First initialize default template file as below.

```bash
$ gofra template init
```

You need to type Author Name, Project Name, Project Address & Server Type.

 A **template.json** file will be generated in the current directory which looks like this.

```json
{
    "author":"Author Name",
    "project":"Project Name",
    "version":"v1",
    "type":"grpc",
    "server":
    {
        "addr":"localhost:58888"
    },
    "monitor_package":
    {
        "package":"github.com/DarkMetrix/gofra/pkg/grpc-utils/monitor/statsd",
        "init_param":"\"127.0.0.1:8125\", \"Project Name\""
    },
    "tracing_package":
    {
        "package":"github.com/DarkMetrix/gofra/pkg/grpc-utils/tracing/zipkin",
        "init_param":"\"http://127.0.0.1:9411/api/v1/spans\", \"false\", \"localhost:58888\", \"Project Name\""
    },
    "interceptor_package":
    {
        "monitor_package":"github.com/DarkMetrix/gofra/pkg/grpc-utils/interceptor/statsd_interceptor",
        "tracing_package":"github.com/DarkMetrix/gofra/pkg/grpc-utils/interceptor/zipkin_interceptor"
    }
}
```

| Json object                         | Value                                                        | Remark       |
| ----------------------------------- | ------------------------------------------------------------ | ------------ |
| author                              | Author's name                                                | User defined |
| project                             | Project's name                                               | User defined |
| version                             | Version                                                      | User defined |
| type                                | Type                                                         | grpc or http |
| server.addr                         | Servicehe 's ip and port                                     | User defined |
| monitor_package.package             | Monitor package import path used in the service(by default statsd is used as the monitor backend), you could write your own monitor package as long as you implement Init & Increment interfaces | Pre defined  |
| monitor_package.init_param          | Monitor params used by Init method(it's the statsd address and project name used in your environment) | Pre defined  |
| Tracing_package.package             | Tracing package import path used in the service(by default zipkin is used as the monitor backend), you could write your own tracing package as long as you implement Init interface | Pre defined  |
| Tracing_package.init_param          | Tracing params used by Init method(it's the zipkin address, debug flag, service address and project name used in your environment) | Pre defined  |
| Interceptor_package.monitor_package | Monitor gRPC interceptor(by default statsd is used as the monitor backend) | Pre defined  |
| Interceptor_package.tracing_package | Tracing gPRC interceptor(by default zipkin is used as the tracing system) | Pre defined  |



## gRPC
- [gRPC Service Generation](#grpc-service-generation)
- [Add or Update Services](#add-or-update-services)
- [Implement RPC Methods](#implement-rpc-methods)
- [Compile and Run](#compile-and-run)
- [Test using Health Check](#test-using-health-check)

### gRPC Service Generation

```bash
$gofra init --path=./template.json
```

All files needed are generated.

```bash
$tree
.
|-- api
|   `-- protobuf_spec
|       |-- health_check
|           |-- health_check.pb.go
|           `-- health_check.proto
|-- build
|-- cmd
|   `-- main.go
|-- configs
|   |-- config.toml
|   `-- log.config
|-- go.mod
|-- internal
|   |-- app
|   |   `-- application.go
|   |-- grpc_handler
|   |   |-- HealthCheckService
|   |       |-- HealthCheck.go
|   |       `-- HealthCheckService.go
|   `-- pkg
|       |-- common
|       |   `-- common.go
|       `-- config
|           `-- config.go
|-- log
|-- template.json
`-- test
    `-- main.go

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
|-- api
|   `-- protobuf_spec
|       |-- health_check
|       |   |-- health_check.pb.go
|       |   `-- health_check.proto
|       `-- user
|           |-- user.pb.go
|           `-- user.proto
|-- build
|-- cmd
|   `-- main.go
|-- configs
|   |-- config.toml
|   `-- log.config
|-- go.mod
|-- internal
|   |-- app
|   |   `-- application.go
|   |-- grpc_handler
|   |   |-- HealthCheckService
|   |   |   |-- HealthCheck.go
|   |   |   `-- HealthCheckService.go
|   |   `-- UserService
|   |       |-- AddUser.go
|   |       `-- UserService.go
|   `-- pkg
|       |-- common
|       |   `-- common.go
|       `-- config
|           `-- config.go
|-- log
|-- template.json
`-- test
    `-- main.go
```

#### Update Service

You could update the service if the pb file has updates, such as adding a new method to the UserService.

The only difference between add and update is update won't override the file already generated.

Using add command, a —override=true flag will help to override all the files generated about the pb file.



### Implement RPC Methods

After add or update service, all the thing left is to implement your own business logic.

All basic logging, metrics, tracing, panic recovery have been added already.



#### UserService/UserService.go

You do not need to do anything about this file, just leave it.

```bash
$cat src/handler/UserService/UserService.go

/**********************************
 * Author : foo
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
 * Author : foo
 * Time : 2018-04-01 02:26:20
 **********************************/

package UserService

import (
        "context"

        //Log package
        //log "github.com/cihub/seelog"

        //Monitor package
        //monitor "github.com/DarkMetrix/gofra/pkg/grpc-utils/monitor/statsd"

        //Tracing package
        //tracing "github.com/DarkMetrix/gofra/pkg/grpc-utils/tracing/zipkin"

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
[DEBUG][2018-04-11T10:55:24.656902][test_gofra_test][seelog_interceptor.go:30][GofraClientInterceptorFunc] => invoke success! req=message:"ping" , reply:
```



#### Service output

```bash
[DEBUG][2018-04-11T10:55:24.656279][test_gofra][seelog_interceptor.go:48][GofraServerInterceptorFunc] => handle success! req=message:"ping" , reply:
```

## Http
- [Http Service Generation](#http-service-generation)
- [Add Http Handlers](#add-http-handlers)
- [Implement Http Methods](#implement-http-methods)
- [Compile and Run](#compile-and-run)
- [Test using Health Check](#test-using-health-check)

### Http Service Generation

```bash
$gofra init --path=./template.json
```

All files needed are generated.

```bash
$tree

.
|-- api
|   `-- http_spec
|       `-- post_health
|-- build
|-- cmd
|   `-- main.go
|-- configs
|   |-- config.toml
|   `-- log.config
|-- go.mod
|-- init.sh
|-- internal
|   |-- app
|   |   `-- application.go
|   |-- http_handler
|   |   `-- post_health.go
|   `-- pkg
|       |-- common
|       |   `-- common.go
|       `-- config
|           `-- config.go
|-- log
|-- template.json
`-- test
    `-- main.go

13 directories, 12 files
```

### Add Http Handlers

By default a health handler is generated.

You could using test/main.go to send health http request to test if the service is working fine.

And you could add your own http handler, e.g.:

#### Add Handler

```bash
$gofra http add --uri=/my/test --method=POST
```

**gofra use gin, so '/contact/:name/\*action' like URI is also ok**

Then a new file named MY_TEST.go is generated.

```bash
$ tree

.
|-- api
|   `-- http_spec
|       |-- post_health
|       `-- post_my_test
|-- build
|-- cmd
|   `-- main.go
|-- configs
|   |-- config.toml
|   `-- log.config
|-- go.mod
|-- internal
|   |-- app
|   |   `-- application.go
|   |-- http_handler
|   |   |-- post_health.go
|   |   `-- post_my_test.go
|   `-- pkg
|       |-- common
|       |   `-- common.go
|       `-- config
|           `-- config.go
|-- log
|-- template.json
`-- test
    `-- main.go
```

### Implement Http Methods

After add http handler, all the thing left is to implement your own business logic.

All basic logging, metrics, panic recovery have been added already.

#### http_handler/MY_TEST.go

Implement your own business logic.

```bash
$cat src/http_handler/MY_TEST.go

/**********************************
 * Author : foo
 * Time : 2020-02-20 13:07:51
 **********************************/

package http_handler

import (
	//Log package
	log "github.com/cihub/seelog"

	//Monitor package
	//monitor "github.com/DarkMetrix/gofra/pkg/monitor/statsd"

	//Tracing package
	//tracing "github.com/DarkMetrix/gofra/pkg/tracing/jaeger"

	"github.com/gin-gonic/gin"
)

//URI(for gin use): [POST] -> "/my/test"
func POST_MY_TEST(ctx *gin.Context) {
	log.Tracef("====== POST_MY_TEST start ======")

	/*
	//Parse request
	//TODO: Bind json to request
	var req xxx

	err := ctx.BindJSON(&req)

	if err != nil {
		log.Warnf("ctx.BindJSON failed! error:%v", err.Error())
		ctx.AbortWithStatusJSON(520, gin.H{"ret":-1, "msg":"Bad json body!"})
		return
	}

	//Check params
	//TODO: Check params
	err = checkPOST_MY_TESTParams(&req)

	if err != nil {
		log.Warnf("checkPOST_MY_TESTParams failed! error:%v", err.Error())
		ctx.AbortWithStatusJSON(520, gin.H{"ret":-1, "msg":fmt.Sprintf("Param invalid! error:%v", err.Error())})
		return
	}
	*/

	//Reply success
	ctx.JSON(200, gin.H{"ret":0, "msg":"success"})
}

/*
//TODO: Implement checkPOST_MY_TESTParams function
func checkPOST_MY_TESTParams(req *xxx) error {
	return nil
}
*/
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
Listen on port [localhost:58888]
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /health                   --> test/src/http_handler.HEALTH (6 handlers)
[GIN-debug] POST   /my/test                  --> test/src/http_handler.MY_TEST (6 handlers)
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
[INFO][2018-09-27T19:08:18.473573][test_gofra_test][main.go:55][testHealth] => http request end! url:http://localhost:58888/health, resp body:{"msg":"success","ret":0}
```

#### Service output

```bash
[GIN] 2018/09/27 - 19:08:18 | 200 |      95.701µs |       127.0.0.1 | POST     /health
[DEBUG][2018-09-27T19:08:18.473165][test_gofra][seelog_middleware.go:28][func1] => Handle success! URI:/health, Host:localhost:58888, Remote address:127.0.0.1:47856, header:map[User-Agent:[Go-http-client/1.1] Content-Length:[2] Content-Type:[application/json] Accept-Encoding:[gzip]]
```


## License

### Gofra

[MIT license](https://github.com/DarkMetrix/gofra/blob/master/LICENSE)

### Dependencies
- google.golang.org/grpc [Apache 2.0 License](https://github.com/grpc/grpc-go/blob/master/LICENSE)
- github.com/gin-gonic/gin [MIT License](https://github.com/gin-gonic/gin/blob/master/LICENSE)
- github.com/cihub/seelog [BSD License](https://github.com/cihub/seelog/blob/master/LICENSE.txt)
- github.com/spf13/viper [MIT License](https://github.com/spf13/viper/blob/master/LICENSE)
- github.com/spf13/cobra [Apache 2.0 License](https://github.com/spf13/cobra/blob/master/LICENSE.txt)
- github.com/go-ozzo/ozzo-validation [MIT License](https://github.com/go-ozzo/ozzo-validation/blob/master/LICENSE)
- github.com/mitchellh/go-homedir" [MIT License](https://github.com/mitchellh/go-homedir/blob/master/LICENSE)
- github.com/grpc-ecosystem/go-grpc-middleware [Apache 2.0 License](https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/LICENSE)
- github.com/tallstoat/pbparser [MIT License](https://github.com/tallstoat/pbparser/blob/master/LICENSE)
- github.com/alexcesaro/statsd [MIT License](https://github.com/alexcesaro/statsd/blob/master/LICENSE)
- github.com/silenceper/pool [MIT License](https://github.com/silenceper/pool/blob/master/LICENSE)
- github.com/opentracing/opentracing-go [MIT License](https://github.com/opentracing/opentracing-go/blob/master/LICENSE)
- github.com/openzipkin/zipkin-go-opentracing [MIT License](https://github.com/openzipkin/zipkin-go-opentracing/blob/master/LICENSE)
- github.com/jaegertracing/jaeger-client-go [Apache 2.0 License](https://github.com/jaegertracing/jaeger-client-go/blob/master/LICENSE)

