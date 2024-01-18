package database

import (
	"sync"

	"github.com/kiwipanel/scaffolding/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	once          sync.Once // create sync.Once primitive)
	DB            *gorm.DB
	database_path = "./state/database/kiwipanel.sqlite"
)

func Connect(app *config.AppConfig) {
	once.Do(func() {
		db, err := gorm.Open(sqlite.Open(database_path), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		app.DB = db
		DB = db
	})
}
