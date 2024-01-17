package routes

import (
	"github.com/kiwipanel/scaffolding/internal/app/admin/controllers"
	"github.com/labstack/echo/v4"
)

func PublicRoutes(r *echo.Echo) {
	controller := controllers.New()
	r.GET("/admin", controller.AdminHompage)

}
