package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Run(name string, args ...string) {
	fmt.Printf("→ %s %s\n", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  ignored error: %v\n", err)
	}
}
