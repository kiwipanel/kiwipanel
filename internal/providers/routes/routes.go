package routes

import (
	"github.com/kiwipanel/scaffolding/internal/app/authentication/routes"
	"github.com/labstack/echo"
)

func ProvidersRoutes(r *echo.Echo) {
	routes.PublicRoutes(r)
}
