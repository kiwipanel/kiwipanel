package app

import (
	"fmt"

	"github.com/kiwipanel/kiwipanel/config"
	"github.com/kiwipanel/kiwipanel/pkg/database"
	"github.com/kiwipanel/kiwipanel/pkg/routing"
	"github.com/kiwipanel/kiwipanel/pkg/sessionstore"
)

var app config.AppConfig
var r = routing.GetRoute()

func Register(mode string) {
	env := config.NewENV()
	app.UseCache = true

	app.KIWIPANEL_MODE = mode
	fmt.Println("env loaded in Register.com: ", env)

	routing.Register(&app)
	sessionstore.Register(r)
	database.Connect(&app)
}
