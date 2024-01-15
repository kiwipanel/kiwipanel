package routing

import (
	"github.com/kiwipanel/scaffolding/internal/providers/routes"
	"github.com/labstack/echo"
)

var r = echo.New()

func Router() {
	routes.ProvidersRoutes(r)

}

func GetRoute() *echo.Echo {
	return r
}

func Run() {
	Router()
	r.Logger.Fatal(r.Start(":7879"))
}
