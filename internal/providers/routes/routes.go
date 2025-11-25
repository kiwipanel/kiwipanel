package routes

import (
	"github.com/kiwipanel/kiwpanel/config"
	admin "github.com/kiwipanel/kiwpanel/internal/app/admin/routes"
	public "github.com/kiwipanel/kiwpanel/internal/app/home/routes"
	"github.com/labstack/echo/v4"
)

func ProvidersRoutes(r *echo.Echo, appconfig *config.AppConfig) {
	public.PublicRoutes(r)
	admin.PublicRoutes(r, appconfig)
}
