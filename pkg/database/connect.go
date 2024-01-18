package database

import (
	"sync"

	"github.com/kiwipanel/scaffolding/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var once sync.Once // create sync.Once primitive

func Connect(app *config.AppConfig) {
	once.Do(func() {
		db, err := gorm.Open(sqlite.Open("./state/database/kiwipanel.sqlite"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		type Product struct {
			gorm.Model
			Code  string
			Price uint
		}
		app.DB = db
		db.AutoMigrate(&Product{})
	})
}
