package routes

import (
	"github.com/kiwipanel/kiwpanel/internal/app/home/controllers"
	"github.com/labstack/echo/v4"
)

func PublicRoutes(r *echo.Echo) {
	controller := controllers.New()
	r.GET("/", controller.Homepage)
	r.GET("/:passcode", controller.Homepage)
	r.GET("/hello", controller.Hello)
	r.POST("/home/login", controller.Login)

}
