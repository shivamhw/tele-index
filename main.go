package main

import (
	"flag"
	"github.com/shivamhw/tele-index/app/api"
	"github.com/shivamhw/tele-index/migrations"
	"github.com/shivamhw/tele-index/pkg/indexer"
)

func main(){
	migrate := flag.String("migrate","", "db folder path")
	flag.Parse()
	if *migrate != "" {
		migrations.Migrate_db(*migrate)
		return
	}
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
