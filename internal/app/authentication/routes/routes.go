package routes

import (
	"github.com/kiwipanel/scaffolding/internal/app/authentication/controllers"
	"github.com/labstack/echo"
)

func PublicRoutes(r *echo.Echo) {
	controller := controllers.New()
	r.GET("/", controller.Homepage)
}
