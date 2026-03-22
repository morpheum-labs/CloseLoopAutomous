package platform

// Build carries linker-injected version metadata (see cmd/arms and the repo Makefile -ldflags).
type Build struct {
	Version string
	Commit  string
}
