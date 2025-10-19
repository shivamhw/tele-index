package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamhw/tele-index/pkg/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		// Authentication logic goes here
		token := g.GetHeader("Authorization")
		if token == "" {
			g.AbortWithStatusJSON(401, gin.H{"msg": "missing or invalid auth token"})
			return
		}
		vToken, err := auth.VerifyToken(token, "testToken")
		if err != nil {
			g.AbortWithStatusJSON(401, gin.H{"msg": "invalid auth token"})
			return
		}
		userID := vToken["user_id"].(string)
		if userID == "" {
			g.AbortWithStatusJSON(401, gin.H{"msg": "invalid userid"})
			return
		}
		username := vToken["username"].(string)
		if userID == "" {
			g.AbortWithStatusJSON(401, gin.H{"msg": "invalid username"})
			return
		}
		g.Set("user_id", userID)
		g.Set("username", username)
		g.Next()
	}
}