package controllers

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/kiwipanel/kiwpanel/internal/app/admin/models"
	"github.com/labstack/echo/v4"
)

func (app *Controller) AdminHompage(c echo.Context) error {

	randomNumber := rand.Intn(100) + 1

	fmt.Println("MODE in adminhome controller: ", app.config.KIWIPANEL_MODE)

	models.Create(app.config.DB, "hello", randomNumber)

	return c.String(http.StatusOK, "hello admin")

}
