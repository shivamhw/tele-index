package main

import (
	"github.com/shivamhw/tele-index/app/api"
	"github.com/shivamhw/tele-index/pkg/indexer"
)

func main(){
	app, err := api.NewApp(&api.AppOpts{
		IndxrOpts: &indexer.IndexerOpts{
			IndexPath: "./prodIdx",
			IndexSchema: "./pkg/indexer/index.json",
		},
	})
	if err != nil {
		panic(err)
	}
	app.RunApp()
}