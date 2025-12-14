package app

import (
	"net/http"

	"github.com/kiwipanel/kiwipanel/config"
	"github.com/kiwipanel/kiwipanel/pkg/database"
)

var appConfig config.AppConfig

func Register(environment string) http.Handler {

	appConfig.UseCache = true
	appConfig.KIWIPANEL_MODE = environment

	r := NewRoutes(&appConfig)
	database.Connect(&appConfig)
	return r
}
