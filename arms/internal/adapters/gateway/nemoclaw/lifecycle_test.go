package nemoclaw

import (
	"context"
	"testing"
)

func TestEnsureSandboxRunning_skipsWhenUnset(t *testing.T) {
	ctx := context.Background()
	if err := EnsureSandboxRunning(ctx, "", "mybox"); err != nil {
		t.Fatal(err)
	}
	if err := EnsureSandboxRunning(ctx, "/bin/nemoclaw", ""); err != nil {
		t.Fatal(err)
	}
}
