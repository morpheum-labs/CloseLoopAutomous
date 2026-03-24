package ports

import (
	"context"

	"github.com/morpheumstreet/CloseLoopAutomous/arms/internal/domain"
)

// AgentProfileRepository persists unified [domain.AgentIdentity] rows keyed by [domain.FleetProfileID].
type AgentProfileRepository interface {
	Upsert(ctx context.Context, gatewayID string, ident *domain.AgentIdentity) error
	ByID(ctx context.Context, id string) (*domain.AgentIdentity, error)
	List(ctx context.Context, limit int) ([]domain.AgentIdentity, error)
	DeleteByGatewayID(ctx context.Context, gatewayID string) error
}

// RemoteAgentProfileSource lists remote agent profiles for a gateway endpoint (e.g. OpenClaw agents.list).
type RemoteAgentProfileSource interface {
	ListRemoteProfiles(ctx context.Context, ep *domain.GatewayEndpoint) ([]domain.AgentIdentity, error)
}

// GeoIPResolver resolves a hostname from a gateway URL to [domain.GeoLocation] using offline data (e.g. MaxMind GeoLite2).
// Implementations return (nil, nil) when lookup is not applicable or data is missing.
type GeoIPResolver interface {
	LookupHost(ctx context.Context, host string) (*domain.GeoLocation, error)
}
