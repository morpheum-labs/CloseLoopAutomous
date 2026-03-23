// Package nullclaw wires arms task dispatch to a NullClaw gateway when it exposes an
// OpenClaw-compatible WebSocket RPC surface (Mission Control handshake + chat.send with sessionKey).
// Upstream: https://github.com/nullclaw/nullclaw
//
// If NullClaw’s wire protocol diverges, add a dedicated ports.AgentGateway implementation here
// and branch gateway.NewFromConfig for config.AgentGatewayDriverNullClawWS.
package nullclaw

import (
	"github.com/closeloopautomous/arms/internal/adapters/gateway/openclaw"
)

// NewOpenClawCompatible returns the shared WebSocket client used for OpenClaw-class gateways.
// NullClaw documents OpenClaw-compatible configuration and gateway behavior for this path.
func NewOpenClawCompatible(opts openclaw.Options) *openclaw.Client {
	return openclaw.New(opts)
}
