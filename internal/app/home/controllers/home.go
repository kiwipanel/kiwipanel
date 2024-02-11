package controllers

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/kiwipanel/scaffolding/internal/app/panel/models"
	"github.com/kiwipanel/scaffolding/pkg/helpers"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (app *Controller) Homepage(c echo.Context) error {

	// sess, err := session.Get("user_authenticated", c)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return c.String(http.StatusNotAcceptable, "Cannot find the page. err - authentication failed")
	// }
	//username := sess.Values["foo"].(string)

	//	return c.String(http.StatusOK, "Hello, "+username)

	// If
	passcode := c.Param("passcode")
	if len(passcode) < 1 || len(passcode) != 9 {
		return c.String(http.StatusOK, "Cannot find the page. Using your passcode to access.")
	}
	secure, err := models.ReadPanel()

	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusNotAcceptable, "Cannot find the page. err - authentication failed")
	}
	if secure.Passcode != passcode {
		fmt.Println(err)
		fmt.Println("secure", len(secure.Passcode))
		fmt.Println("passcode", len(passcode))
		return c.String(http.StatusNotAcceptable, "Passcode is not correct. err - authentication failed")
	}

	return c.Render(http.StatusOK, "signup", "")
	//return c.String(http.StatusOK, "Hello there, kiwipanel.org!. It is good, update")

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

	ip, err := helpers.GetLocalIP()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	fmt.Println("Local IP Address:", ip)

	return c.String(http.StatusOK, "Hello, cập nhật "+username+"ip: "+ip)

	return c.String(http.StatusOK, "hello sesion")

}
