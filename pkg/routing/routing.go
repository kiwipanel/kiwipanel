package routing

import (
	"github.com/kiwipanel/scaffolding/internal/providers/routes"
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

func Register() {
	routes.ProvidersRoutes(r)
	middlewares()
}
