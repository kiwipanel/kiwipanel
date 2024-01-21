package bootstrap

import (
	"fmt"
	"os"

	"github.com/kiwipanel/scaffolding/pkg/ui/view"
)

func Root(flag string) {
	Register()
	Migrate()
	Setup()
	root_path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(root_path)

	fmt.Println("flag: ", flag)

	view.Loadmode(flag)

	r.Logger.Fatal(r.Start(":8443"))
}
