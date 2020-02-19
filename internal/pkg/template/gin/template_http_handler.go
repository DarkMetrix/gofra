package gin

type HttpHandlerInfo struct {
	Author string
	Time string

	HandlerName string

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
	//Log package
	log "github.com/cihub/seelog"

	//Monitor package
	//monitor "{{.MonitorPackage}}"

	//Tracing package
	//tracing "{{.TracingPackage}}"

	"github.com/gin-gonic/gin"
)

func {{.HandlerName}}(ctx *gin.Context) {
	log.Tracef("====== {{.HandlerName}} start ======")

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
	err = check{{.HandlerName}}Params(&req)

	if err != nil {
		log.Warnf("check{{.HandlerName}}Params failed! error:%v", err.Error())
		ctx.AbortWithStatusJSON(520, gin.H{"ret":-1, "msg":fmt.Sprintf("Param invalid! error:%v", err.Error())})
		return
	}
	*/

	//Reply success
	ctx.JSON(200, gin.H{"ret":0, "msg":"success"})
}

/*
//TODO: Implement check{{.HandlerName}}Params function
func check{{.HandlerName}}Params(req *xxx) error {
	return nil
}
*/
`
