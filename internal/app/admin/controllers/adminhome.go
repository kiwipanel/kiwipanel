package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *Controller) AdminHompage(c echo.Context) error {
	return c.String(http.StatusOK, "hello admin")

}
