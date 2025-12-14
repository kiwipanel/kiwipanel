package controllers

import (
	"fmt"
	"net/http"
)

func (app *Controller) AdminHompage(w http.ResponseWriter, r *http.Request) {

	fmt.Println("MODE in adminhome controller: ", app.config.KIWIPANEL_MODE)

	w.Write([]byte("hi admin"))

}
