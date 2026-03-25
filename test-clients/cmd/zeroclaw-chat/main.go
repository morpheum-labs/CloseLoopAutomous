// Zeroclaw chat WebSocket test client: /ws/chat with bearer.<token> + zeroclaw.v1 subprotocols.
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"

	"github.com/morpheumstreet/CloseLoopAutomous/test-clients/pkg/zeroclaw/chatws"
)

func main() {
	host := flag.String("host", "localhost:20013", "Zeroclaw gateway host:port (no scheme)")
	token := flag.String("token", "", "Zeroclaw bearer token (required)")
	tls := flag.Bool("tls", false, "use wss:// instead of ws://")
	flag.Parse()

	if *token == "" {
		log.Fatal("missing -token (pair via dashboard at http://<host> and read localStorage zeroclaw_token)")
	}

	u := chatws.BuildURL(*host, *tls)
	log.Printf("dialing %s", u)

	conn, resp, err := chatws.Dial(*host, *token, *tls)
	if err != nil {
		if resp != nil {
			log.Printf("HTTP response: %s", resp.Status)
		}
		log.Fatalf("websocket dial: %v", err)
	}
	defer conn.Close()

	log.Printf("connected; negotiated subprotocol: %q", conn.Subprotocol())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	readDone := make(chan struct{})
	go func() {
		defer close(readDone)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Printf("read: %v", err)
				}
				return
			}
			fmt.Printf("<- %s\n", message)
		}
	}()

	fmt.Println("Enter lines to send as text frames (empty line exits). Ctrl+C to quit.")

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
			_ = conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return
		case <-readDone:
			return
		case line, ok := <-lines:
			if !ok {
				return
			}
			if line == "" {
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
				log.Fatalf("write: %v", err)
			}
		}
	}
}
