package controllers

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

func (app *Controller) Homepage(c echo.Context) error {
	passcode := c.Param("passcode")

	//TODO: Check the session if logged in
	if len(passcode) < 1 {
		return c.String(http.StatusOK, "Cannot find the page. Using your passcode to access.")
	}

	fmt.Println(reflect.TypeOf(passcode))
	return c.String(http.StatusOK, "Hello there, kiwipanel.org!. It is good")
}

func (app *Controller) HomeAccess(c echo.Context) error {
	passcode := c.Param("passcode")
	fmt.Println(reflect.TypeOf(passcode))
	return c.String(http.StatusOK, passcode)
}

func (app *Controller) Hello(c echo.Context) error {
	passcode := c.Param("passcode")
	fmt.Println(reflect.TypeOf(passcode))
	return c.Render(http.StatusOK, "base", "")
}
