package static

import "github.com/labstack/echo/v4"

func Register(r *echo.Echo) {
	r.Static("/assets", "assets")
}
