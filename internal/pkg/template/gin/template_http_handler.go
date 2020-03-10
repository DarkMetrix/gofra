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

	// reply success
	ctx.JSON(200, gin.H{"ret":0, "msg":"success"})
}

