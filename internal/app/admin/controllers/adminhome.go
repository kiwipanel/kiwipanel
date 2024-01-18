package controllers

import (
	"math/rand"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (app *Controller) AdminHompage(c echo.Context) error {

	type Product struct {
		gorm.Model
		Code  string
		Price int
	}

	randomNumber := rand.Intn(100) + 1

	app.config.DB.Create(&Product{Code: "D42", Price: randomNumber})

	return c.String(http.StatusOK, "hello admin")

}
