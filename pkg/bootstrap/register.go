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

	app.UseCache = true
	app.InProduction = true

	app.KIWIPANEL_MODE = "development"

	fmt.Println("Show app", app)
	fmt.Println("env loaded in Register.com: ", env)

	view.LoadConfig(app)

	static.Register(r)
	r.Renderer = view.RenderTemplates
	routing.Register(&app)
	sessionstore.Register(r)
	database.Connect(&app)
}
