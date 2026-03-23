package inkos

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath_home(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home dir")
	}
	got := ExpandPath("~/inkos/ws")
	want := filepath.Join(home, "inkos/ws")
	if got != want {
		t.Fatalf("ExpandPath(~/inkos/ws) = %q want %q", got, want)
	}
}
