package models

import (
	"github.com/kiwipanel/scaffolding/pkg/database/schema"
	"gorm.io/gorm"
)

func Create(DB *gorm.DB, code string, price int) {
	DB.Create(&schema.Product{Code: code, Price: price})
}
