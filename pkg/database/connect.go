package database

import (
	"fmt"
	"sync"

	"github.com/kiwipanel/scaffolding/config"
	// "gorm.io/driver/sqlite" # needs to use cgo
	// pure go
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	once                      sync.Once // create sync.Once primitive)
	DB                        *gorm.DB
	database_path_development = "state/database/kiwipanel.sqlite"
	database_path_production  = "/home/scaffolding/state/database/kiwipanel.sqlite"
	database_path             string
)

func Connect(app *config.AppConfig) {
	once.Do(func() {
		if app.KIWIPANEL_MODE == "production" {
			database_path = database_path_production
		} else {
			database_path = database_path_development
		}
		fmt.Println("Inside connect database, mode: ", app.KIWIPANEL_MODE)

		db, err := gorm.Open(sqlite.Open(database_path), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		app.DB = db
		DB = db
	})
}
