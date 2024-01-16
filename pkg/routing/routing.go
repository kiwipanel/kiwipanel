package routing

import (
	"github.com/kiwipanel/scaffolding/internal/providers/routes"
	"github.com/kiwipanel/scaffolding/pkg/ui/view"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var r = echo.New()

func registers() {
	r.Renderer = view.RenderTemplates
	r.Static("/assets", "assets")
	routes.ProvidersRoutes(r)
}

func middlewares() {
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Pre(middleware.RemoveTrailingSlash())
	r.Use(middleware.Secure())
}

func GetRoute() *echo.Echo {
	return r
}

func Run() {
	registers()
	middlewares()
	r.Logger.Fatal(r.Start(":7879"))
}
