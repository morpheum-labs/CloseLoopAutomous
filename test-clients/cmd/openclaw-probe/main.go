// OpenClaw gateway probe using github.com/a3tai/openclaw-go: full connect handshake
// (connect.challenge → connect → hello-ok). Optional sessions.list or chat.send.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/a3tai/openclaw-go/gateway"
	"github.com/a3tai/openclaw-go/protocol"
)

func main() {
	rawURL := flag.String("url", "ws://127.0.0.1:18789/ws", "gateway WebSocket URL (include /ws path)")
	token := flag.String("token", "", "auth token (connect auth.token)")
	password := flag.String("password", "", "gateway password (connect auth.password; alternative to -token)")
	timeout := flag.Duration("timeout", 30*time.Second, "overall timeout for connect and optional RPCs")
	listSessions := flag.Bool("list-sessions", false, "call sessions.list after connect")
	chat := flag.String("chat", "", "if set, send one chat.send with this message and print the returned chat event")
	sessionKey := flag.String("session", "main", "session key for -chat")
	flag.Parse()

	if *token != "" && *password != "" {
		log.Fatal("use only one of -token or -password")
	}

	wsURL := normalizeWSURL(*rawURL)

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	opts := []gateway.Option{
		gateway.WithConnectTimeout(*timeout),
	}
	if t := *token; t != "" {
		opts = append(opts, gateway.WithToken(t))
	}
	if p := *password; p != "" {
		opts = append(opts, gateway.WithPassword(p))
	}

	client := gateway.NewClient(opts...)
	defer client.Close()

	if err := client.Connect(ctx, wsURL); err != nil {
		log.Fatalf("connect: %v", err)
	}

	hello := client.Hello()
	if hello == nil {
		log.Fatal("connected but Hello() is nil")
	}
	fmt.Printf("hello-ok: protocol=%d server=%s connId=%s\n",
		hello.Protocol, hello.Server.Version, hello.Server.ConnID)

	if *listSessions {
		limit := 20
		raw, err := client.SessionsList(ctx, protocol.SessionsListParams{Limit: &limit})
		if err != nil {
			log.Fatalf("sessions.list: %v", err)
		}
		var buf bytes.Buffer
		if err := json.Indent(&buf, raw, "", "  "); err != nil {
			fmt.Printf("sessions.list: %s\n", string(raw))
		} else {
			fmt.Printf("sessions.list:\n%s\n", buf.String())
		}
	}

	if msg := strings.TrimSpace(*chat); msg != "" {
		ev, err := client.ChatSend(ctx, protocol.ChatSendParams{
			SessionKey:     *sessionKey,
			Message:        msg,
			IdempotencyKey: fmt.Sprintf("probe-%d", time.Now().UnixNano()),
		})
		if err != nil {
			log.Fatalf("chat.send: %v", err)
		}
		out, _ := json.MarshalIndent(ev, "", "  ")
		fmt.Printf("chat.send result:\n%s\n", out)
	}
}

// normalizeWSURL ensures a path; OpenClaw gateways use /ws.
func normalizeWSURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if u.Path == "" || u.Path == "/" {
		u.Path = "/ws"
	}
	return u.String()
}
