package bootstrap

import (
	"fmt"
	"os"
)

func Root(flag string) {
	view.Loadmode(flag)
	Register(flag)
	Migrate()
	Setup()
	root_path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(root_path)

	fmt.Println("flag: ", flag)

	r.Logger.Fatal(r.Start(":8443"))
}
