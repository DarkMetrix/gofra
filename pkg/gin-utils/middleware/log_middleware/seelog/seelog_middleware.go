package seelog

import (
	log "github.com/cihub/seelog"

	"github.com/gin-gonic/gin"
)

func GetMiddleware() gin.HandlerFunc{
	return func (ctx *gin.Context) {
		//Before request
		header := ctx.Request.Header
		host := ctx.Request.Host
		remoteAddr := ctx.Request.RemoteAddr
		uri := ctx.Request.RequestURI

		log.Tracef("Handle begin! URI:%v, Host:%v, Remote address:%v, header:%v", uri, host, remoteAddr, header)

		//Switch to another middleware handler
		ctx.Next()

		//After request
		status := ctx.Writer.Status()

		if status != 200 {
			log.Warnf("Handle failed! URI:%v, Host:%v, Remote address:%v, header:%v", uri, host, remoteAddr, header)
		} else {
			log.Debugf("Handle success! URI:%v, Host:%v, Remote address:%v, header:%v", uri, host, remoteAddr, header)
		}
	}
}

