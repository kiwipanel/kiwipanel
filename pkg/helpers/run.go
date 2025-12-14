package helpers

import (
	"bytes"
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

func ExecCommand(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

func ExecCombinedOutput(name string, args ...string) (string, error) {
	cmd := ExecCommand(name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return buf.String(), fmt.Errorf("%w", err)
	}
	return buf.String(), nil
}
