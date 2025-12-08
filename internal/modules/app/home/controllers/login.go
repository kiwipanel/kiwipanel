package controllers

import (
	"fmt"
	"net/http"

	"github.com/kiwipanel/kiwipanel/pkg/helpers"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (app *Controller) Login(c echo.Context) error {
	passcode := c.Param("passcode")
	fmt.Println(passcode)
	sess, _ := session.Get("user_authenticated", c)
	sess.Values["user"] = "bar update"
	sess.Save(c.Request(), c.Response())

	sess.Values["foo"] = "bar update"
	sess.Save(c.Request(), c.Response())

	username := sess.Values["foo"].(string)

	ip, err := helpers.GetLocalIP()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	fmt.Println("Local IP Address:", ip)

	return c.String(http.StatusOK, "Hello, cập nhật "+username+"ip: "+ip)

	return c.String(http.StatusOK, passcode)
}
