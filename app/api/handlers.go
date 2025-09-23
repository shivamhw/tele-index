package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shivamhw/tele-index/internal/models"
	"github.com/shivamhw/tele-index/pkg/log"
)


func (a *App) searchItem(c *gin.Context) {
	var logger *slog.Logger
	if val, ok := c.Get(log.LoggerKeyType{}); ok {
		logger = val.(*slog.Logger)
	}
	var data models.ItemSearchRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error("failed parsing payload", "err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid payload"})
		return
	}
	logger.Info("searching item", "query: ", data)
	if data.Size == 0 {
		data.Size = 10
	}
	res, err := a.indx.Search(data.Q, data.Size, data.From, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	logger.Info("found results for search", "nof", len(res))
	c.JSON(200, gin.H{
		"count": len(res),
		"data": res,
	})
}

func (a *App) saveItem(c *gin.Context) {
	var logger *slog.Logger
	if val, ok := c.Get(log.LoggerKeyType{}); ok {
		logger = val.(*slog.Logger)
	}
	var data models.ItemAddRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		logger.Error("failed parsing payload", "err:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid payload"})
		return
	}
	logger.Info("saving item", "item: ", data)
	if err := a.indx.Index(data.GetItem()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	c.JSON(201, gin.H{"msg": "success"})
}

func (a *App) getCount(c *gin.Context) {
	count, err := a.indx.GetTotalDocs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err})
		return
	}
	c.JSON(200, gin.H{"count": count})
}

func (a *App) getAll(c *gin.Context) {
	a.indx.ListAll()
}