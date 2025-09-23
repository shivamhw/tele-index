package api

import (

	"github.com/gin-gonic/gin"
	"github.com/shivamhw/tele-index/pkg/indexer"
)

type App struct {
	g    *gin.Engine
	indx *indexer.Indexer
}

type AppOpts struct {
	IndxrOpts *indexer.IndexerOpts
}

func NewApp(opts *AppOpts) (*App, error) {
	app := &App{}
	g := gin.Default()
	RegisterRoutes(g, app)
	indx, err := indexer.NewIndexer(opts.IndxrOpts)
	if err != nil {
		return nil, err
	}
	app.g = g
	app.indx = indx
	return app, nil
}

func (a *App) RunApp() {
	a.g.Run()
}