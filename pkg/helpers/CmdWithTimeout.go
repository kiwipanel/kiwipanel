package helpers

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func CmdWithTimeout(timeout time.Duration, name string, args ...string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return false, fmt.Sprintf("command timeout (%v)", timeout)
	}

	if err != nil {
		return false, string(out)
	}

	return true, string(out)
}
