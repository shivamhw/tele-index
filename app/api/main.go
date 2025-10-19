package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shivamhw/tele-index/pkg/db"
	"github.com/shivamhw/tele-index/pkg/indexer"
	"github.com/shivamhw/tele-index/pkg/log"
)

type App struct {
	g    *gin.Engine
	indx *indexer.Indexer
	db  db.DB
}

type AppOpts struct {
	IndxrOpts *indexer.IndexerOpts
	DBOpts   *db.SqliteOpts
}

var logger = log.Logger

func NewApp(opts *AppOpts) (*App, error) {
	app := &App{}
	g := gin.Default()
	db, err := db.NewSqlDb(&db.SqliteOpts{
		Path: "./prodDB.sqlite",
	})
	if err != nil {
		return nil, err
	}
	RegisterRoutes(g, app)
	indx, err := indexer.NewIndexer(opts.IndxrOpts)
	if err != nil {
		return nil, err
	}
	app.g = g
	app.indx = indx
	app.db = db
	return app, nil
}

func (a *App) RunApp() {
	a.g.Run()
}