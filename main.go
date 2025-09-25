package main

import (
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source/file"
	"github.com/shivamhw/tele-index/app/api"
	"github.com/shivamhw/tele-index/pkg/db"
	"github.com/shivamhw/tele-index/pkg/indexer"
)

func main(){
	migrate := flag.String("migrate","", "db folder path")
	flag.Parse()
	if *migrate != "" {
		migrate_db(*migrate)
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

func migrate_db(path string) {
	db, err := db.NewSqlDb(&db.SqliteOpts{
		Path: path,
	})
	if err != nil {
		panic(err)
	}
	instance, err := sqlite3.WithInstance(db.Db, &sqlite3.Config{})
	if err != nil {
		panic(err)
	}
	fScr, err := (&file.File{}).Open("./migrations")
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithInstance("file", fScr, "sqlite3", instance)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil {
		panic(err)
	}
	fmt.Print("Migration Successfull \n")
}