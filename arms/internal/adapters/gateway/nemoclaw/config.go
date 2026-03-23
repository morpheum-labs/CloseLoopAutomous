package nemoclaw

// PoolSettings holds process-wide defaults for the NemoClaw adapter (from config / env).
type PoolSettings struct {
	BinaryPath       string
	AutoStart        bool
	DefaultBlueprint string // reserved for future policy / blueprint hooks
}
