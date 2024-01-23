package helpers

import (
	"fmt"
	"os"
)

func ExcutablePath() string {
	root_path, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	return root_path
}
