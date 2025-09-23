package middleware

import (
"github.com/gin-gonic/gin"
"github.com/shivamhw/tele-index/pkg/log"
)


func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx = log.With(ctx)
		ctx.Next()
	}
}