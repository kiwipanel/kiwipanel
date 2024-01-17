package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() {
	db, err := gorm.Open(sqlite.Open("./state/database/kiwipanel.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	type Product struct {
		gorm.Model
		Code  string
		Price uint
	}
	db.AutoMigrate(&Product{})
}
