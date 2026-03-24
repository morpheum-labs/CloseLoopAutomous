package agentidentity

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/morpheumstreet/CloseLoopAutomous/arms/internal/domain"
	"github.com/morpheumstreet/CloseLoopAutomous/arms/internal/ports"
)

// Service discovers remote agent profiles per gateway endpoint, upserts fleet cache rows, and links the execution registry.
type Service struct {
	Endpoints ports.GatewayEndpointRegistry
	Profiles  ports.AgentProfileRepository
	Registry  ports.ExecutionAgentRegistry
	Source    ports.RemoteAgentProfileSource
	Geo       ports.GeoIPResolver
	Events    ports.LiveActivityPublisher
	Clock     func() time.Time
}

// RefreshAll clears cached profiles per gateway, re-runs remote discovery (or synthetic rows), and upserts agent_profiles.
func (s *Service) RefreshAll(ctx context.Context) error {
	if s == nil || s.Endpoints == nil || s.Profiles == nil {
		return nil
	}
	list, err := s.Endpoints.List(ctx, 500)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	if s.Clock != nil {
		now = s.Clock().UTC()
	}
	var execAgents []domain.ExecutionAgent
	if s.Registry != nil {
		execAgents, err = s.Registry.List(ctx, 5000)
		if err != nil {
			return err
		}
	}
	regIndex := buildRegistryMatchIndex(execAgents)

	for i := range list {
		ep := &list[i]
		if err := s.Profiles.DeleteByGatewayID(ctx, ep.ID); err != nil {
			return err
		}
		hostGeo := lookupGatewayGeo(ctx, s.Geo, ep)

		var rows []domain.AgentIdentity
		if s.Source != nil {
			pctx, cancel := context.WithTimeout(ctx, 45*time.Second)
			discovered, derr := s.Source.ListRemoteProfiles(pctx, ep)
			cancel()
			switch {
			case derr != nil && errors.Is(derr, domain.ErrRemoteAgentListUnsupported):
				rows = []domain.AgentIdentity{unsupportedDiscoveryIdentity(ep, now, hostGeo)}
			case derr != nil:
				rows = []domain.AgentIdentity{gatewayScanErrorIdentity(ep, derr, now, hostGeo)}
			case len(discovered) == 0:
				rows = []domain.AgentIdentity{emptyRemoteListIdentity(ep, now, hostGeo)}
			default:
				rows = make([]domain.AgentIdentity, 0, len(discovered))
				for j := range discovered {
					final := finalizeFleetIdentity(ep, &discovered[j], regIndex, now, hostGeo)
					rows = append(rows, final)
				}
			}
		} else {
			rows = []domain.AgentIdentity{unsupportedDiscoveryIdentity(ep, now, hostGeo)}
		}

		for j := range rows {
			row := &rows[j]
			if err := s.Profiles.Upsert(ctx, ep.ID, row); err != nil {
				return err
			}
			if s.Events != nil {
				_ = s.Events.Publish(ctx, ports.LiveActivityEvent{
					Type: "agent_identity_updated",
					Ts:   now.Format(time.RFC3339Nano),
					Data: map[string]any{
						"identity_id": row.ID,
						"gateway_id":  ep.ID,
						"driver":      ep.Driver,
					},
				})
			}
		}
	}
	return nil
}

func lookupGatewayGeo(ctx context.Context, geo ports.GeoIPResolver, ep *domain.GatewayEndpoint) *domain.GeoLocation {
	if geo == nil || ep == nil || strings.TrimSpace(ep.GatewayURL) == "" {
		return nil
	}
	host := hostFromGatewayURL(ep.GatewayURL)
	if host == "" {
		return nil
	}
	gctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	g, err := geo.LookupHost(gctx, host)
	if err != nil {
		slog.Default().Debug("agentidentity geo lookup", "host", host, "err", err)
	}
	if g != nil && g.Source != "" && g.Source != "none" {
		return g
	}
	return nil
}

func hostFromGatewayURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" {
		return ""
	}
	return u.Hostname()
}
