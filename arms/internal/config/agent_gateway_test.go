package config

import (
	"testing"
	"time"
)

func TestResolveAgentGateway_autoEmptyURLIsStub(t *testing.T) {
	c := Config{AgentGatewayDriver: "auto", OpenClawDispatchTimeout: 45 * time.Second}
	got := c.ResolveAgentGateway()
	if got.Driver != AgentGatewayDriverStub {
		t.Fatalf("driver: got %q want stub", got.Driver)
	}
	if got.Timeout != 45*time.Second {
		t.Fatalf("timeout: got %v", got.Timeout)
	}
}

func TestResolveAgentGateway_autoWithURLIsOpenClaw(t *testing.T) {
	c := Config{
		AgentGatewayDriver:      "auto",
		OpenClawGatewayURL:      "wss://x/ws",
		OpenClawGatewayToken:    "t",
		OpenClawSessionKey:      "agent:main:x",
		OpenClawDispatchTimeout: 10 * time.Second,
		ArmsDeviceID:            "dev1",
	}
	got := c.ResolveAgentGateway()
	if got.Driver != AgentGatewayDriverOpenClawWS {
		t.Fatalf("driver: got %q", got.Driver)
	}
	if got.URL != "wss://x/ws" || got.Token != "t" || got.SessionKey != "agent:main:x" || got.DeviceID != "dev1" {
		t.Fatalf("fields: %+v", got)
	}
}

func TestResolveAgentGateway_nullclawUsesDedicatedEnvThenFallback(t *testing.T) {
	c := Config{
		AgentGatewayDriver:      "nullclaw_ws",
		NullClawGatewayURL:      "wss://null/ws",
		NullClawGatewayToken:    "nt",
		NullClawSessionKey:      "agent:null:1",
		OpenClawGatewayURL:      "wss://ignored/ws",
		OpenClawGatewayToken:    "ignored",
		OpenClawSessionKey:    "ignored",
		OpenClawDispatchTimeout: 30 * time.Second,
	}
	got := c.ResolveAgentGateway()
	if got.Driver != AgentGatewayDriverNullClawWS {
		t.Fatalf("driver: got %q", got.Driver)
	}
	if got.URL != "wss://null/ws" || got.Token != "nt" || got.SessionKey != "agent:null:1" {
		t.Fatalf("fields: %+v", got)
	}

	c2 := Config{
		AgentGatewayDriver:   "nullclaw",
		OpenClawGatewayURL:   "wss://fallback/ws",
		OpenClawGatewayToken: "ft",
		OpenClawSessionKey:   "sk",
	}
	got2 := c2.ResolveAgentGateway()
	if got2.URL != "wss://fallback/ws" || got2.Token != "ft" || got2.SessionKey != "sk" {
		t.Fatalf("fallback fields: %+v", got2)
	}
}

func TestResolveAgentGateway_openclawEmptyURLStub(t *testing.T) {
	c := Config{AgentGatewayDriver: "openclaw_ws"}
	if g := c.ResolveAgentGateway(); g.Driver != AgentGatewayDriverStub {
		t.Fatalf("want stub, got %+v", g)
	}
}

func TestResolveAgentGateway_unknownDriverStub(t *testing.T) {
	c := Config{AgentGatewayDriver: "unknown-thing", OpenClawGatewayURL: "wss://x/ws"}
	if g := c.ResolveAgentGateway(); g.Driver != AgentGatewayDriverStub {
		t.Fatalf("want stub, got %+v", g)
	}
}
