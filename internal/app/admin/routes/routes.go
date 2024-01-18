package routes

import (
	"github.com/kiwipanel/scaffolding/config"
	"github.com/kiwipanel/scaffolding/internal/app/admin/controllers"
	"github.com/labstack/echo/v4"
)

func PublicRoutes(r *echo.Echo, appconfig *config.AppConfig) {
	controller := controllers.New(appconfig)
	r.GET("/admin", controller.AdminHompage)

}
