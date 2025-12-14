package controllers

import (
	"net/http"
)

func (app *Controller) Homepage(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("hi"))

	//return c.String(http.StatusOK, "Hello there, kiwipanel.org!. It is good, update")

}
