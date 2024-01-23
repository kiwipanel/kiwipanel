package bootstrap

import (
	"github.com/kiwipanel/scaffolding/pkg/ui/view"
)

func Root(flag string) {
	view.Loadmode(flag)
	Register(flag)
	Migrate()
	Setup()
	r.Logger.Fatal(r.Start(":8443"))
}
