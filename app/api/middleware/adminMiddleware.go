package middleware

import "github.com/gin-gonic/gin"


func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Admin authentication logic goes here
		c.Next()
	}
}