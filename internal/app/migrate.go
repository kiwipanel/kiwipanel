package app

import (
	"github.com/kiwipanel/kiwipanel/pkg/database"
	"github.com/kiwipanel/kiwipanel/pkg/database/schema"
)

func Migrate() {
	database.DB.AutoMigrate(
		&schema.Panel{},
		&schema.Product{},
		&schema.User{},
	)
}
