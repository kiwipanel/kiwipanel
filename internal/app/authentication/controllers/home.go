package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

func (app *Controller) Homepage(c echo.Context) error {
	return c.String(http.StatusOK, "Hello there, kiwipanel.org!. It is good")

}
