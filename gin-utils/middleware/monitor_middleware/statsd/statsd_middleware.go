package statsd

import (
	"fmt"

	monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"

	"github.com/gin-gonic/gin"
)

func GetMiddleware() gin.HandlerFunc{
	return func (ctx *gin.Context) {
		client := monitor.GetStatsd()

		if client != nil {
			//Before request
			defer monitor.GetStatsd().NewTiming().Send(fmt.Sprintf("%v|Time,type=Server.Time", ctx.Request.RequestURI))

			monitor.Increment(fmt.Sprintf("%v,type=Server.Total", ctx.Request.RequestURI))

			//Switch to another middleware handler
			ctx.Next()

			//After request
			status := ctx.Writer.Status()

			if status != 200 {
				monitor.Increment(fmt.Sprintf("%v,type=Server.Fail", ctx.Request.RequestURI))
			} else {
				monitor.Increment(fmt.Sprintf("%v,type=Server.Success", ctx.Request.RequestURI))
			}
		} else {
			monitor.Increment(fmt.Sprintf("%v,type=Server.Total", ctx.Request.RequestURI))

			//Switch to another middleware handler
			ctx.Next()

			//After request
			status := ctx.Writer.Status()

			if status != 200 {
				monitor.Increment(fmt.Sprintf("%v,type=Server.Fail", ctx.Request.RequestURI))
			} else {
				monitor.Increment(fmt.Sprintf("%v,type=Server.Success", ctx.Request.RequestURI))
			}
		}
	}
}

