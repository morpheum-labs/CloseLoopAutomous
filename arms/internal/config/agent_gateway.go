package config

import (
	"strings"
	"time"
)

// AgentGatewayDriver selects the execution-gateway implementation for task dispatch
// ([ports.AgentGateway]). Unknown values resolve to stub at runtime.
type AgentGatewayDriver string

const (
	// AgentGatewayDriverAuto picks stub when no gateway URL is configured; otherwise openclaw_ws (legacy default).
	AgentGatewayDriverAuto AgentGatewayDriver = "auto"
	// AgentGatewayDriverStub never dials; returns synthetic external refs (local dev / tests).
	AgentGatewayDriverStub AgentGatewayDriver = "stub"
	// AgentGatewayDriverOpenClawWS uses the Mission Control–compatible WebSocket client (chat.send).
	AgentGatewayDriverOpenClawWS AgentGatewayDriver = "openclaw_ws"
	// AgentGatewayDriverNullClawWS uses the NullClaw-oriented adapter. Today NullClaw documents
	// OpenClaw-compatible gateway wire format, so this reuses the same client until a divergent
	// protocol needs a dedicated implementation (github.com/nullclaw/nullclaw).
	AgentGatewayDriverNullClawWS AgentGatewayDriver = "nullclaw_ws"
)

// AgentGatewaySettings is the normalized input for constructing [ports.AgentGateway].
type AgentGatewaySettings struct {
	Driver     AgentGatewayDriver
	URL        string
	Token      string
	DeviceID   string
	SessionKey string
	Timeout    time.Duration
}

// ResolveAgentGateway maps file/env config into a single struct for the gateway factory.
func (c Config) ResolveAgentGateway() AgentGatewaySettings {
	to := c.OpenClawDispatchTimeout
	if to <= 0 {
		to = 30 * time.Second
	}
	d := normalizeAgentGatewayDriver(c.AgentGatewayDriver)
	url := strings.TrimSpace(c.OpenClawGatewayURL)
	tok := strings.TrimSpace(c.OpenClawGatewayToken)
	sess := strings.TrimSpace(c.OpenClawSessionKey)
	dev := strings.TrimSpace(c.ArmsDeviceID)

	stub := AgentGatewaySettings{Driver: AgentGatewayDriverStub, Timeout: to}

	switch d {
	case "", AgentGatewayDriverAuto:
		if url == "" {
			return stub
		}
		return AgentGatewaySettings{
			Driver:     AgentGatewayDriverOpenClawWS,
			URL:        url,
			Token:      tok,
			DeviceID:   dev,
			SessionKey: sess,
			Timeout:    to,
		}
	case AgentGatewayDriverStub:
		return stub
	case AgentGatewayDriverOpenClawWS:
		if url == "" {
			return stub
		}
		return AgentGatewaySettings{
			Driver:     AgentGatewayDriverOpenClawWS,
			URL:        url,
			Token:      tok,
			DeviceID:   dev,
			SessionKey: sess,
			Timeout:    to,
		}
	case AgentGatewayDriverNullClawWS:
		if u := strings.TrimSpace(c.NullClawGatewayURL); u != "" {
			url = u
		}
		if t := strings.TrimSpace(c.NullClawGatewayToken); t != "" {
			tok = t
		}
		if s := strings.TrimSpace(c.NullClawSessionKey); s != "" {
			sess = s
		}
		if url == "" {
			return stub
		}
		return AgentGatewaySettings{
			Driver:     AgentGatewayDriverNullClawWS,
			URL:        url,
			Token:      tok,
			DeviceID:   dev,
			SessionKey: sess,
			Timeout:    to,
		}
	default:
		return stub
	}
}

func normalizeAgentGatewayDriver(s string) AgentGatewayDriver {
	v := strings.ToLower(strings.TrimSpace(s))
	switch v {
	case "", "auto", "default":
		return AgentGatewayDriverAuto
	case "stub", "none", "off", "disabled":
		return AgentGatewayDriverStub
	case "openclaw", "openclaw_ws", "openclaw-ws":
		return AgentGatewayDriverOpenClawWS
	case "nullclaw", "nullclaw_ws", "nullclaw-ws":
		return AgentGatewayDriverNullClawWS
	default:
		return AgentGatewayDriver("unknown:" + v)
	}
}
