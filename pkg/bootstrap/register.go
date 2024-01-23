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

func Register(mode string) {
	env := config.NewENV()
	app.UseCache = true

	app.KIWIPANEL_MODE = mode
	fmt.Println("env loaded in Register.com: ", env)

	view.Loadmode("hello")
	static.Register(r)
	r.Renderer = view.RenderTemplates
	routing.Register(&app)
	sessionstore.Register(r)
	database.Connect(&app)
}
