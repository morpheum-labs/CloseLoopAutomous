package gateway

import (
	"context"
	"log/slog"

	"github.com/closeloopautomous/arms/internal/adapters/gateway/nullclaw"
	"github.com/closeloopautomous/arms/internal/adapters/gateway/openclaw"
	"github.com/closeloopautomous/arms/internal/config"
	"github.com/closeloopautomous/arms/internal/domain"
	"github.com/closeloopautomous/arms/internal/ports"
)

// NewFromConfig builds [ports.AgentGateway] from resolved config (driver + connection fields).
// cleanup closes WebSocket state when non-stub; invoke on process shutdown (e.g. App.Close).
func NewFromConfig(cfg config.Config, knowledgeForDispatch func(ctx context.Context, productID domain.ProductID, query string) (string, error)) (ports.AgentGateway, func()) {
	s := cfg.ResolveAgentGateway()
	opts := openclaw.Options{
		URL:                  s.URL,
		Token:                s.Token,
		DeviceID:             s.DeviceID,
		SessionKey:           s.SessionKey,
		Timeout:              s.Timeout,
		KnowledgeForDispatch: knowledgeForDispatch,
	}

	switch s.Driver {
	case config.AgentGatewayDriverStub:
		return &Stub{}, func() {}

	case config.AgentGatewayDriverOpenClawWS:
		slog.Info("arms agent gateway", "driver", string(s.Driver))
		c := openclaw.New(opts)
		return c, func() { _ = c.Close() }

	case config.AgentGatewayDriverNullClawWS:
		slog.Info("arms agent gateway", "driver", string(s.Driver))
		c := nullclaw.NewOpenClawCompatible(opts)
		return c, func() { _ = c.Close() }

	default:
		return &Stub{}, func() {}
	}
}
