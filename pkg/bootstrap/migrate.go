package bootstrap

import (
	"github.com/kiwipanel/scaffolding/pkg/database"
	"github.com/kiwipanel/scaffolding/pkg/database/schema"
)

func Migrate() {
	database.DB.AutoMigrate(&schema.Panel{}, &schema.Product{}, &schema.User{})
}
