// Interactive OpenClaw gateway chat via github.com/a3tai/openclaw-go (WebSocket + chat.send).
// Streams assistant deltas from chat events on stdout; not the Zeroclaw /ws/chat protocol.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/a3tai/openclaw-go/gateway"
	"github.com/a3tai/openclaw-go/protocol"
)

func main() {
	rawURL := flag.String("url", "ws://127.0.0.1:18789/ws", "gateway WebSocket URL (include /ws path)")
	token := flag.String("token", "", "auth token (required unless -password)")
	password := flag.String("password", "", "gateway password (alternative to -token)")
	sessionKey := flag.String("session", "main", "chat session key")
	flag.Parse()

	if *token == "" && *password == "" {
		log.Fatal("missing -token or -password")
	}
	if *token != "" && *password != "" {
		log.Fatal("use only one of -token or -password")
	}

	wsURL := normalizeWSURL(*rawURL)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var outMu sync.Mutex
	printlnLocked := func(a ...any) {
		outMu.Lock()
		defer outMu.Unlock()
		fmt.Println(a...)
	}
	printLocked := func(a ...any) {
		outMu.Lock()
		defer outMu.Unlock()
		fmt.Print(a...)
	}

	opts := []gateway.Option{
		gateway.WithConnectTimeout(30 * time.Second),
	}
	if t := *token; t != "" {
		opts = append(opts, gateway.WithToken(t))
	}
	if p := *password; p != "" {
		opts = append(opts, gateway.WithPassword(p))
	}
	opts = append(opts, gateway.WithOnEvent(func(ev protocol.Event) {
		if ev.EventName != protocol.EventChat {
			return
		}
		var data map[string]any
		if json.Unmarshal(ev.Payload, &data) != nil {
			return
		}
		state, _ := data["state"].(string)
		switch state {
		case "delta":
			switch msg := data["message"].(type) {
			case string:
				printLocked(msg)
			default:
				b, _ := json.Marshal(msg)
				printLocked(string(b))
			}
		case "final":
			printLocked("\n")
			printlnLocked("  [chat] done")
		case "error":
			errMsg, _ := data["errorMessage"].(string)
			printlnLocked("  [chat] error:", errMsg)
		}
	}))

	client := gateway.NewClient(opts...)
	defer client.Close()

	log.Printf("dialing %s", wsURL)
	if err := client.Connect(ctx, wsURL); err != nil {
		log.Fatalf("connect: %v", err)
	}
	h := client.Hello()
	if h != nil {
		log.Printf("connected (protocol=%d server=%s)", h.Protocol, h.Server.Version)
	}

	fmt.Println("Enter lines to send as chat messages (empty line exits). Ctrl+C to quit.")

	lines := make(chan string)
	go func() {
		sc := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !sc.Scan() {
				if err := sc.Err(); err != nil {
					log.Printf("stdin: %v", err)
				}
				close(lines)
				return
			}
			lines <- sc.Text()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-lines:
			if !ok {
				return
			}
			if strings.TrimSpace(line) == "" {
				return
			}
			_, err := client.ChatSend(ctx, protocol.ChatSendParams{
				SessionKey:     *sessionKey,
				Message:        line,
				IdempotencyKey: fmt.Sprintf("cli-%d", time.Now().UnixNano()),
			})
			if err != nil {
				log.Printf("chat.send: %v", err)
				continue
			}
		}
	}
}

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
