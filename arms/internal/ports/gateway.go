package ports

import (
	"context"

	"github.com/closeloopautomous/arms/internal/domain"
)

// AgentGateway is the execution plane: arms pushes work to an external agent runtime and receives
// completion via webhooks or operator APIs. Implementations include in-process stubs, OpenClaw-class
// WebSocket gateways, and adapters for other runtimes (e.g. NullClaw when wire-compatible).
type AgentGateway interface {
	DispatchTask(ctx context.Context, task domain.Task) (externalRef string, err error)
	DispatchSubtask(ctx context.Context, parent domain.Task, sub domain.Subtask) (externalRef string, err error)
}
