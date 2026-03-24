package agentidentity

import (
	"strings"
	"time"

	"github.com/morpheumstreet/CloseLoopAutomous/arms/internal/domain"
)

func buildRegistryMatchIndex(agents []domain.ExecutionAgent) map[string]string {
	out := make(map[string]string)
	for i := range agents {
		a := &agents[i]
		key := strings.TrimSpace(a.EndpointID) + "\x00" + strings.TrimSpace(a.SessionKey)
		out[key] = a.ID
	}
	return out
}

func remoteStableKey(ident *domain.AgentIdentity) string {
	if ident == nil {
		return "_"
	}
	if ident.Custom != nil {
		if sk, ok := ident.Custom["session_key"].(string); ok && strings.TrimSpace(sk) != "" {
			return strings.TrimSpace(sk)
		}
		if r, ok := ident.Custom["remote_agent_id"].(string); ok && strings.TrimSpace(r) != "" {
			return strings.TrimSpace(r)
		}
	}
	if strings.TrimSpace(ident.ID) != "" {
		return strings.TrimSpace(ident.ID)
	}
	return "_"
}

func finalizeFleetIdentity(ep *domain.GatewayEndpoint, raw *domain.AgentIdentity, regIndex map[string]string, now time.Time, geo *domain.GeoLocation) domain.AgentIdentity {
	if raw == nil {
		return domain.AgentIdentity{}
	}
	origKey := remoteStableKey(raw)
	shortID := strings.TrimSpace(raw.ID)
	if raw.Custom == nil {
		raw.Custom = map[string]any{}
	}
	if _, has := raw.Custom["remote_agent_id"]; !has {
		if shortID != "" {
			raw.Custom["remote_agent_id"] = shortID
		} else {
			raw.Custom["remote_agent_id"] = origKey
		}
	}
	raw.Custom["gateway_endpoint_id"] = ep.ID

	if _, has := raw.Custom["suggested_session_key"]; !has {
		suggested := origKey
		if sk, ok := raw.Custom["session_key"].(string); ok && strings.TrimSpace(sk) != "" {
			suggested = strings.TrimSpace(sk)
		}
		raw.Custom["suggested_session_key"] = suggested
	}
	suggested := strings.TrimSpace(asString(raw.Custom["suggested_session_key"]))

	raw.ID = domain.FleetProfileID(ep.ID, origKey)
	raw.Driver = domain.NormalizeGatewayDriver(ep.Driver)
	if strings.TrimSpace(raw.GatewayURL) == "" {
		raw.GatewayURL = ep.GatewayURL
	}
	raw.LastSeen = now
	if geo != nil && geo.Source != "" && geo.Source != "none" {
		raw.Geo = geo
	}

	matchKey := strings.TrimSpace(ep.ID) + "\x00" + strings.TrimSpace(suggested)
	if execID, ok := regIndex[matchKey]; ok && execID != "" {
		raw.Custom["on_registry"] = true
		raw.Custom["execution_agent_id"] = execID
	} else {
		raw.Custom["on_registry"] = false
		delete(raw.Custom, "execution_agent_id")
	}
	if raw.Custom["discovery_kind"] == nil {
		raw.Custom["discovery_kind"] = "remote_list"
	}
	return *raw
}

func gatewayScanErrorIdentity(ep *domain.GatewayEndpoint, err error, now time.Time, geo *domain.GeoLocation) domain.AgentIdentity {
	cust := map[string]any{
		"gateway_endpoint_id": ep.ID,
		"discovery_kind":      "gateway_scan_error",
		"discovery_error":     err.Error(),
		"on_registry":         false,
	}
	return domain.AgentIdentity{
		ID:         domain.FleetProfileID(ep.ID, "_scan_error"),
		GatewayURL: ep.GatewayURL,
		Name:       "Gateway scan failed",
		Driver:     domain.NormalizeGatewayDriver(ep.Driver),
		Status:     domain.StatusError,
		LastSeen:   now,
		Custom:     cust,
		Geo:        geo,
	}
}

func unsupportedDiscoveryIdentity(ep *domain.GatewayEndpoint, now time.Time, geo *domain.GeoLocation) domain.AgentIdentity {
	return domain.AgentIdentity{
		ID:         domain.FleetProfileID(ep.ID, "_remote_list_unsupported"),
		GatewayURL: ep.GatewayURL,
		Name:       "Remote agent list not supported",
		Driver:     domain.NormalizeGatewayDriver(ep.Driver),
		Status:     domain.StatusOffline,
		LastSeen:   now,
		Custom: map[string]any{
			"gateway_endpoint_id": ep.ID,
			"discovery_kind":      "remote_list_unsupported",
			"on_registry":         false,
		},
		Geo: geo,
	}
}

func emptyRemoteListIdentity(ep *domain.GatewayEndpoint, now time.Time, geo *domain.GeoLocation) domain.AgentIdentity {
	return domain.AgentIdentity{
		ID:         domain.FleetProfileID(ep.ID, "_no_remote_profiles"),
		GatewayURL: ep.GatewayURL,
		Name:       "No remote agent profiles",
		Driver:     domain.NormalizeGatewayDriver(ep.Driver),
		Status:     domain.StatusOffline,
		LastSeen:   now,
		Custom: map[string]any{
			"gateway_endpoint_id": ep.ID,
			"discovery_kind":      "no_remote_profiles",
			"on_registry":         false,
		},
		Geo: geo,
	}
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}
