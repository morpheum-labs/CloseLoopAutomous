package openclaw

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/morpheumstreet/CloseLoopAutomous/arms/internal/domain"
)

// mapAgentsListPayload maps OpenClaw agents.list JSON payload to domain.AgentIdentity.
// Accepts { "agents": [...] } or a bare JSON array (legacy / alternate gateways).
func mapAgentsListPayload(raw json.RawMessage, gatewayURL string, now time.Time) ([]domain.AgentIdentity, error) {
	raw = unwrapAgentsListEnvelope(bytesTrimSpaceJSON(raw))
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var asRoot []json.RawMessage
	if json.Unmarshal(raw, &asRoot) == nil {
		return mapAgentListEntries(asRoot, gatewayURL, now)
	}
	var env struct {
		Agents []json.RawMessage `json:"agents"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("openclaw agents.list payload: %w", err)
	}
	return mapAgentListEntries(env.Agents, gatewayURL, now)
}

// unwrapAgentsListEnvelope peels a single-key payload/result wrapper (some gateways double-wrap).
func unwrapAgentsListEnvelope(raw json.RawMessage) json.RawMessage {
	for len(raw) > 0 && string(raw) != "null" {
		var probe map[string]json.RawMessage
		if json.Unmarshal(raw, &probe) != nil || len(probe) != 1 {
			return raw
		}
		var next json.RawMessage
		for _, k := range []string{"payload", "result"} {
			if inner, ok := probe[k]; ok {
				next = inner
				break
			}
		}
		if len(bytesTrimSpaceJSON(next)) == 0 {
			return raw
		}
		raw = bytesTrimSpaceJSON(next)
	}
	return raw
}

func mapAgentListEntries(entries []json.RawMessage, gatewayURL string, now time.Time) ([]domain.AgentIdentity, error) {
	out := make([]domain.AgentIdentity, 0, len(entries))
	for i, e := range entries {
		var row struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			Identity    string `json:"identity"`
			Model       string `json:"model"`
			SessionKey  string `json:"sessionKey"`
		}
		_ = json.Unmarshal(e, &row)
		id := strings.TrimSpace(row.ID)
		if id == "" {
			id = strings.TrimSpace(row.SessionKey)
		}
		if id == "" {
			id = fmt.Sprintf("openclaw-agent-%d", i)
		}
		name := strings.TrimSpace(row.Name)
		if name == "" {
			name = strings.TrimSpace(row.DisplayName)
		}
		if name == "" {
			name = id
		}
		custom := map[string]any{"source": "agents.list"}
		if row.Identity != "" {
			custom["openclaw_identity"] = row.Identity
		}
		if row.Model != "" {
			custom["model"] = row.Model
		}
		if row.SessionKey != "" {
			custom["session_key"] = row.SessionKey
		}
		out = append(out, domain.AgentIdentity{
			ID:         id,
			GatewayURL: strings.TrimSpace(gatewayURL),
			Name:       name,
			Driver:     domain.GatewayDriverOpenClawWS,
			Version:    "1.0.1",
			Status:     domain.StatusOnline,
			LastSeen:   now,
			Platform: domain.PlatformInfo{
				OS:       "unknown",
				Arch:     "unknown",
				Hostname: "openclaw-gateway",
			},
			Metrics: domain.Metrics{},
			Custom:  custom,
		})
	}
	return out, nil
}
