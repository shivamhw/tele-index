package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamhw/tele-index/app/api/middleware"
)


func RegisterRoutes(g *gin.Engine, a *App) {
	g.GET("/collections/:id", a.getCollection)
	g.Use(middleware.LoggerMiddleware())

	itemsG := g.Group("/idx")// TODO: add admin middleware
	itemsG.POST("/create", a.saveItem)
	itemsG.GET("/get", a.searchItem)
	itemsG.GET("/count", a.getCount)
	itemsG.GET("/all", a.getAll)

	authG := g.Group("/auth")
	authG.POST("/login", a.loginUser)
	authG.POST("/register", a.createUser)

	userG := g.Group("/user")
	userG.Use(middleware.AuthMiddleware())
	userG.GET("/profile", a.getUser)
	userG.POST("/collections", a.createCollection)
	userG.GET("/collections", a.getUserCollections)
	userG.POST("/collections/:id/item", a.addItemToCollection)
	userG.GET("addTelegram", a.updateTelegram)

	teleG := g.Group("/tele")
	teleG.Use(middleware.AuthMiddleware())
	teleG.GET("/sendOtp", a.sendOtp)
	teleG.GET("/verifyOtp", a.verifyOtp)
}