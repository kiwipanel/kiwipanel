package controllers

import (
	"math/rand"
	"net/http"

	"github.com/kiwipanel/scaffolding/internal/app/admin/models"
	"github.com/labstack/echo/v4"
)

func (app *Controller) AdminHompage(c echo.Context) error {

	randomNumber := rand.Intn(100) + 1

	models.Create(app.config.DB, "hello", randomNumber)

	return c.String(http.StatusOK, "hello admin")

}
