package routing

import (
	"net/http"

	"github.com/labstack/echo"
)

func Router() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, kiwipanel.org!")
	})
	e.Logger.Fatal(e.Start(":7879"))
}
