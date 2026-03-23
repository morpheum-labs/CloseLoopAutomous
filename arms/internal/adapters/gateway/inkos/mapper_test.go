package inkos

import (
	"strings"
	"testing"
)

func TestBuildWriteNextArgs(t *testing.T) {
	got := BuildWriteNextArgs("my-book", "")
	want := []string{"write", "next", "my-book", "--count", "1", "--json"}
	if len(got) != len(want) {
		t.Fatalf("len got %d want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("idx %d: got %q want %q (full %v)", i, got[i], want[i], got)
		}
	}

	withCtx := BuildWriteNextArgs("b", "hello world")
	if len(withCtx) != 8 {
		t.Fatalf("expected 8 args, got %d: %v", len(withCtx), withCtx)
	}
	if withCtx[6] != "--context" || withCtx[7] != "hello world" {
		t.Fatalf("unexpected context tail: %v", withCtx[6:])
	}
}

func TestTruncateContext(t *testing.T) {
	long := strings.Repeat("α", 5)
	got := truncateContext(long, 3)
	if got != "ααα" {
		t.Fatalf("truncate runes: got %q", got)
	}
}

func TestExternalRefFromInkOSStdout(t *testing.T) {
	raw := `{"requestId":"req-1","chapter":2}`
	if got := ExternalRefFromInkOSStdout([]byte(raw)); got != "req-1" {
		t.Fatalf("got %q want req-1", got)
	}
	nested := `{"data":{"jobId":"j9"}}`
	if got := ExternalRefFromInkOSStdout([]byte(nested)); got != "j9" {
		t.Fatalf("nested got %q want j9", got)
	}
	if ExternalRefFromInkOSStdout([]byte("not json")) != "" {
		t.Fatal("expected empty for invalid json")
	}
}
