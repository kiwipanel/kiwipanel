package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (app *Controller) Login(c echo.Context) error {
	passcode := c.Param("passcode")
	fmt.Println(passcode)
	sess, _ := session.Get("user_authenticated", c)
	sess.Values["user"] = "bar update"
	sess.Save(c.Request(), c.Response())

	return c.String(http.StatusOK, passcode)
}
