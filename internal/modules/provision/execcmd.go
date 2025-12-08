package provision

import (
	"bytes"
	"fmt"
	"os/exec"
)

// ExecCommand wraps exec.Command. Use .Run() or .CombinedOutput() on result.
func ExecCommand(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// ExecCombinedOutput executes command and returns stdout+stderr as string
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
