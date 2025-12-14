package controllers

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/kiwipanel/kiwipanel/internal/modules/app/admin/models"
)

func (app *Controller) AdminHompage(w http.ResponseWriter, r *http.Request) {

	randomNumber := rand.Intn(100) + 1

	fmt.Println("MODE in adminhome controller: ", app.config.KIWIPANEL_MODE)

	models.Create(app.config.DB, "hello", randomNumber)

	w.Write([]byte("hi admin"))

}
