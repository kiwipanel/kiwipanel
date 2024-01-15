package routes

import (
	"net/http"

	"github.com/labstack/echo"
)

func PublicRoutes(r *echo.Echo) {
	r.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello there, kiwipanel.org!")
	})
}
