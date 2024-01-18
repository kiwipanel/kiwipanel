package bootstrap

import (
	"fmt"
	"os"
)

func Root() {
	Register()
	Migrate()
	Setup()
	root_path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(root_path)
	r.Logger.Fatal(r.Start(":8443"))
}
