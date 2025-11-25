package bootstrap

import (
	"github.com/kiwipanel/kiwpanel/pkg/database"
	"github.com/kiwipanel/kiwpanel/pkg/database/schema"
)

func Migrate() {
	database.DB.AutoMigrate(
		&schema.Panel{},
		&schema.Product{},
		&schema.User{},
	)
}
