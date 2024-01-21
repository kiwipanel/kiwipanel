package bootstrap

import (
	"fmt"

	"github.com/kiwipanel/scaffolding/config"
	"github.com/kiwipanel/scaffolding/pkg/database"
	"github.com/kiwipanel/scaffolding/pkg/routing"
	"github.com/kiwipanel/scaffolding/pkg/sessionstore"
	"github.com/kiwipanel/scaffolding/pkg/ui/static"
	"github.com/kiwipanel/scaffolding/pkg/ui/view"
)

var app config.AppConfig
var r = routing.GetRoute()

func Register() {
	env := config.NewENV()
	view.LoadConfig(*env)

	fmt.Println(env.KIWIPANEL_MODE)

	// Setting the config for the app from the ENV
	app.KIWIPANEL_MODE = env.KIWIPANEL_MODE

	static.Register(r)
	r.Renderer = view.RenderTemplates
	routing.Register(&app)
	sessionstore.Register(r)
	database.Connect(&app)
	// config := config.NewConfigServer()
	// fmt.Println(config.Server.Port)

}
