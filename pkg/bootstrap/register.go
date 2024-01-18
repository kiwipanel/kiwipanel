package bootstrap

import (
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
	static.Register(r)
	r.Renderer = view.RenderTemplates
	routing.Register(&app)
	sessionstore.Register(r)
	database.Connect(&app)
	// config := config.NewConfigServer()
	// fmt.Println(config.Server.Port)

}
