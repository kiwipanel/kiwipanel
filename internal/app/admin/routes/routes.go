package routes

import (
	"github.com/kiwipanel/kiwpanel/config"
	"github.com/kiwipanel/kiwpanel/internal/app/admin/controllers"
	"github.com/labstack/echo/v4"
)

func PublicRoutes(r *echo.Echo, appconfig *config.AppConfig) {
	controller := controllers.New(appconfig)
	r.GET("/admin", controller.AdminHompage)

}
