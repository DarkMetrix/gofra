package seelog

import (
	"runtime/debug"

	log "github.com/cihub/seelog"

	"github.com/gin-gonic/gin"
)

func GetMiddleware() gin.HandlerFunc{
	return func (ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Warnf("Got panic! error:%v, stack:%v", err, string(debug.Stack()))
				ctx.AbortWithStatus(500)
			}
		}()

		//Switch to another middleware handler
		ctx.Next()
	}
}

