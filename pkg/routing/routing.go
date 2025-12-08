package routing

import (
	"github.com/kiwipanel/kiwipanel/config"
	"github.com/kiwipanel/kiwipanel/internal/modules/providers/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var r = echo.New()

func GetRoute() *echo.Echo {
	return r
}

func middlewares() {
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Pre(middleware.RemoveTrailingSlash())
	r.Use(middleware.Secure())
}

func Register(appconfig *config.AppConfig) {
	routes.ProvidersRoutes(r, appconfig)
	middlewares()
}
