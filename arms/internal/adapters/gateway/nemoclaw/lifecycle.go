package nemoclaw

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// EnsureSandboxRunning runs `nemoclaw <sandbox> start` when binary and sandbox are set.
// It is a best-effort idempotent start used before WebSocket dispatch; skip when binary or sandbox is empty.
func EnsureSandboxRunning(ctx context.Context, nemoclawBin, sandbox string) error {
	bin := strings.TrimSpace(nemoclawBin)
	name := strings.TrimSpace(sandbox)
	if bin == "" || name == "" {
		return nil
	}
	cmd := exec.CommandContext(ctx, bin, name, "start")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return fmt.Errorf("nemoclaw %q start: %w: %s", name, err, msg)
		}
		return fmt.Errorf("nemoclaw %q start: %w", name, err)
	}
	return nil
}
