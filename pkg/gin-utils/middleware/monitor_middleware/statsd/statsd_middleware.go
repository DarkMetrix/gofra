package statsd

import (
	"fmt"

	monitor "github.com/DarkMetrix/gofra/pkg/monitor/statsd"

	"github.com/gin-gonic/gin"
)

func GetMiddleware() gin.HandlerFunc{
	return func (ctx *gin.Context) {
		client := monitor.GetStatsd()

		if client != nil {
			// before request
			defer monitor.GetStatsd().NewTiming().Send(fmt.Sprintf("%v|Time,type=Server.Time", ctx.Request.RequestURI))

			monitor.Increment(fmt.Sprintf("%v,type=Server.Total", ctx.Request.RequestURI))

			// switch to another middleware handler
			ctx.Next()

			// after request
			status := ctx.Writer.Status()

			if status != 200 {
				monitor.Increment(fmt.Sprintf("%v,type=Server.Fail", ctx.Request.RequestURI))
			} else {
				monitor.Increment(fmt.Sprintf("%v,type=Server.Success", ctx.Request.RequestURI))
			}
		} else {
			monitor.Increment(fmt.Sprintf("%v,type=Server.Total", ctx.Request.RequestURI))

			// switch to another middleware handler
			ctx.Next()

			// after request
			status := ctx.Writer.Status()

			if status != 200 {
				monitor.Increment(fmt.Sprintf("%v,type=Server.Fail", ctx.Request.RequestURI))
			} else {
				monitor.Increment(fmt.Sprintf("%v,type=Server.Success", ctx.Request.RequestURI))
			}
		}
	}
}

