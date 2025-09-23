package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamhw/tele-index/app/api/middleware"
)


func RegisterRoutes(g *gin.Engine, a *App) {
	itemsG := g.Group("/items")
	g.Use(middleware.LoggerMiddleware())
	itemsG.POST("/create", a.saveItem)
	itemsG.GET("/get", a.searchItem)
	itemsG.GET("/count", a.getCount)
	itemsG.GET("/all", a.getAll)
}