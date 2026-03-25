// Streams Zeroclaw /api/events (Server-Sent Events) with Bearer auth.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	host := flag.String("host", "localhost:20013", "Zeroclaw gateway host:port")
	token := flag.String("token", "", "Bearer token (required)")
	tls := flag.Bool("tls", false, "use https://")
	flag.Parse()

	if *token == "" {
		log.Fatal("missing -token")
	}
	scheme := "http"
	if *tls {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s/api/events", scheme, *host)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+*token)

	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		log.Fatalf("unexpected status %s: %s", resp.Status, body)
	}

	if _, err := io.Copy(os.Stdout, resp.Body); err != nil && ctx.Err() == nil {
		log.Printf("stream ended: %v", err)
	}
}
