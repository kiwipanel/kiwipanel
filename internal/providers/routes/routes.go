package routes

import (
	admin "github.com/kiwipanel/scaffolding/internal/app/admin/routes"
	public "github.com/kiwipanel/scaffolding/internal/app/authentication/routes"
	"github.com/labstack/echo/v4"
)

func ProvidersRoutes(r *echo.Echo) {
	public.PublicRoutes(r)
	admin.PublicRoutes(r)
}
