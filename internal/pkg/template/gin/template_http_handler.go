package gin

type HttpHandlerInfo struct {
	Author string
	Time string

	HandlerName string
	URI string
	Method string

	MonitorPackage string
	TracingPackage string
}

var HttpHandlerTemplate string = `
/**********************************
 * Author : {{.Author}}
 * Time : {{.Time}}
 **********************************/

package http_handler

import (
	// log package
	log "github.com/cihub/seelog"

	// monitor package
	// monitor "{{.MonitorPackage}}"

	// tracing package
	// opentracing "github.com/opentracing/opentracing-go"

	"github.com/gin-gonic/gin"
)

// URI(for gin use): [{{.Method}}] -> "{{.URI}}"
func {{.HandlerName}}(ctx *gin.Context) {
	log.Tracef("====== {{.HandlerName}} start ======")

	/*
	// parse request
	// TODO: Bind json to request
	var req xxx

	err := ctx.BindJSON(&req)

	if err != nil {
		log.Warnf("ctx.BindJSON failed! error:%v", err.Error())
		ctx.AbortWithStatusJSON(520, gin.H{"ret":-1, "msg":"Bad json body!"})
		return
	}

	// check params
	// TODO: Check params
	err = check{{.HandlerName}}Params(&req)

	if err != nil {
		log.Warnf("check{{.HandlerName}}Params failed! error:%v", err.Error())
		ctx.AbortWithStatusJSON(520, gin.H{"ret":-1, "msg":fmt.Sprintf("Param invalid! error:%v", err.Error())})
		return
	}
	*/

	// reply success
	ctx.JSON(200, gin.H{"ret":0, "msg":"success"})
}

/*
// TODO: Implement check{{.HandlerName}}Params function
func check{{.HandlerName}}Params(req *xxx) error {
	return nil
}
*/
`
