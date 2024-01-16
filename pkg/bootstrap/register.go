package bootstrap

import (
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
	r.Logger.Fatal(r.Start(":7879"))
}
