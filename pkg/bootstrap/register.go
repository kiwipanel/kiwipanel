package bootstrap

import (
	"fmt"

	"github.com/kiwipanel/kiwpanel/config"
	"github.com/kiwipanel/kiwpanel/pkg/database"
	"github.com/kiwipanel/kiwpanel/pkg/routing"
	"github.com/kiwipanel/kiwpanel/pkg/sessionstore"
	"github.com/kiwipanel/kiwpanel/pkg/ui/static"
	"github.com/kiwipanel/kiwpanel/pkg/ui/view"
)

var app config.AppConfig
var r = routing.GetRoute()

func Register(mode string) {
	env := config.NewENV()
	app.UseCache = true

	app.KIWIPANEL_MODE = mode
	fmt.Println("env loaded in Register.com: ", env)

	static.Register(r)
	r.Renderer = view.RenderTemplates
	routing.Register(&app)
	sessionstore.Register(r)
	database.Connect(&app)
}
