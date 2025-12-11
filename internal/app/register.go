package app

import (
	"github.com/kiwipanel/kiwipanel/config"
	"github.com/kiwipanel/kiwipanel/pkg/database"
	"github.com/kiwipanel/kiwipanel/pkg/sessionstore"
)

var appConfig config.AppConfig

func Register(environment string) {

	appConfig.UseCache = true
	appConfig.KIWIPANEL_MODE = environment

	RegisterRoutes(&appConfig)
	sessionstore.Register(r)
	database.Connect(&appConfig)
}
