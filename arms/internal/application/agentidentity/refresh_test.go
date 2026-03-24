package agentidentity

import (
	"context"
	"testing"
	"time"

	"github.com/morpheumstreet/CloseLoopAutomous/arms/internal/adapters/memory"
	"github.com/morpheumstreet/CloseLoopAutomous/arms/internal/domain"
)

type fakeRemoteSource struct {
	out []domain.AgentIdentity
	err error
}

func (f *fakeRemoteSource) ListRemoteProfiles(context.Context, *domain.GatewayEndpoint) ([]domain.AgentIdentity, error) {
	return f.out, f.err
}

func TestRefreshAll_RemoteListAndRegistryLink(t *testing.T) {
	ctx := context.Background()
	eps := memory.NewGatewayEndpointStore()
	profiles := memory.NewAgentProfileStore()
	reg := memory.NewExecutionAgentStore()
	now := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	_ = eps.Save(ctx, &domain.GatewayEndpoint{
		ID: "gw-a", Driver: domain.GatewayDriverOpenClawWS, GatewayURL: "wss://x/ws",
	})
	_ = reg.Save(ctx, &domain.ExecutionAgent{
		ID: "exec-1", DisplayName: "R1", EndpointID: "gw-a", SessionKey: "sess-1",
		CreatedAt: now,
	})
	src := &fakeRemoteSource{out: []domain.AgentIdentity{{
		ID:       "agent-a",
		Name:     "Alpha",
		Status:   domain.StatusOnline,
		LastSeen: now,
		Custom:   map[string]any{"session_key": "sess-1"},
	}}}
	svc := &Service{
		Endpoints: eps,
		Profiles:  profiles,
		Registry:  reg,
		Source:    src,
		Clock:     func() time.Time { return now },
	}
	if err := svc.RefreshAll(ctx); err != nil {
		t.Fatal(err)
	}
	list, err := profiles.List(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("len %d", len(list))
	}
	wantID := domain.FleetProfileID("gw-a", "sess-1")
	if list[0].ID != wantID {
		t.Fatalf("id %q want %q", list[0].ID, wantID)
	}
	if list[0].Custom["on_registry"] != true {
		t.Fatalf("on_registry = %v", list[0].Custom["on_registry"])
	}
	if list[0].Custom["execution_agent_id"] != "exec-1" {
		t.Fatalf("execution_agent_id = %v", list[0].Custom["execution_agent_id"])
	}
}

func TestRefreshAll_UnsupportedDriverRow(t *testing.T) {
	ctx := context.Background()
	eps := memory.NewGatewayEndpointStore()
	profiles := memory.NewAgentProfileStore()
	reg := memory.NewExecutionAgentStore()
	now := time.Now().UTC()
	_ = eps.Save(ctx, &domain.GatewayEndpoint{
		ID: "gw-http", Driver: domain.GatewayDriverMetaClawHTTP, GatewayURL: "https://api.example/v1",
	})
	src := &fakeRemoteSource{err: domain.ErrRemoteAgentListUnsupported}
	svc := &Service{
		Endpoints: eps,
		Profiles:  profiles,
		Registry:  reg,
		Source:    src,
		Clock:     func() time.Time { return now },
	}
	if err := svc.RefreshAll(ctx); err != nil {
		t.Fatal(err)
	}
	list, err := profiles.List(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("len %d", len(list))
	}
	if list[0].Custom["discovery_kind"] != "remote_list_unsupported" {
		t.Fatalf("kind %v", list[0].Custom["discovery_kind"])
	}
}
