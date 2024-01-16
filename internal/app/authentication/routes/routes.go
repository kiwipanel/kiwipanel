package routes

import (
	"github.com/kiwipanel/scaffolding/internal/app/authentication/controllers"
	"github.com/labstack/echo/v4"
)

func PublicRoutes(r *echo.Echo) {
	controller := controllers.New()
	r.GET("/", controller.Homepage)
	r.GET("/:passcode", controller.HomeAccess)
	r.GET("/hello", controller.Hello)
}
