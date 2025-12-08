package routes

import (
	"github.com/kiwipanel/kiwipanel/config"
	admin "github.com/kiwipanel/kiwipanel/internal/modules/app/admin/routes"
	public "github.com/kiwipanel/kiwipanel/internal/modules/app/home/routes"
	"github.com/labstack/echo/v4"
)

func ProvidersRoutes(r *echo.Echo, appconfig *config.AppConfig) {
	public.PublicRoutes(r)
	admin.PublicRoutes(r, appconfig)
}
