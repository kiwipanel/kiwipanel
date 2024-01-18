package controllers

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (app *Controller) Homepage(c echo.Context) error {

	sess, err := session.Get("user_authenticated", c)

	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusNotAcceptable, "Cannot find the page. err - authentication failed")
	}
	username := sess.Values["foo"].(string)
	return c.String(http.StatusOK, "Hello, "+username)

	passcode := c.Param("passcode")

	//TODO: Check the session if logged in, then display the dashboard, if not check the passcode
	//
	if len(passcode) < 1 || len(passcode) != 7 {
		return c.String(http.StatusOK, "Cannot find the page. Using your passcode to access.")
	}

	fmt.Println("passcode: ", passcode)

	return c.String(http.StatusOK, "Hello there, kiwipanel.org!. It is good")
}

func (app *Controller) HomeAccess(c echo.Context) error {
	passcode := c.Param("passcode")
	fmt.Println(reflect.TypeOf(passcode))
	return c.String(http.StatusOK, passcode)
}

func (app *Controller) Hello(c echo.Context) error {

	sess, _ := session.Get("user_authenticated", c)
	sess.Values["foo"] = "bar update"
	sess.Save(c.Request(), c.Response())

	username := sess.Values["foo"].(string)
	return c.String(http.StatusOK, "Hello, "+username)

	return c.String(http.StatusOK, "hello sesion")
	return c.Render(http.StatusOK, "hello", "")
}
