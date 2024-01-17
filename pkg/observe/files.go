package observe

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bitfield/script"
)

func GetFile() {
	script.FindFiles("/usr/bin").Stdout()
}

func GetCurrentFolder() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return dir
}
