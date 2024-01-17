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

func Register() {
	r := routing.GetRoute()
	static.Register(r)
	r.Renderer = view.RenderTemplates
	routing.Register()
	sessionstore.Register(r)
	database.Connect()

	config := config.NewConfigServer()
	fmt.Println(config.Server.Port)

	r.Logger.Fatal(r.Start(":8443"))

}
