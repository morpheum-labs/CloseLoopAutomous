// Package chatws dials Zeroclaw's dashboard-style chat WebSocket (/ws/chat).
// This path is separate from the OpenClaw-class JSON-RPC URL used by arms for zeroclaw_ws dispatch.
package chatws

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const subprotocolZeroclaw = "zeroclaw.v1"

// ChatPath is the HTTP path for the real-time chat socket on the Zeroclaw gateway.
const ChatPath = "/ws/chat"

// BuildURL returns ws:// or wss:// host/path for /ws/chat.
func BuildURL(host string, useTLS bool) string {
	scheme := "ws"
	if useTLS {
		scheme = "wss"
	}
	return fmt.Sprintf("%s://%s%s", scheme, host, ChatPath)
}

// Dial connects with Sec-WebSocket-Protocol: bearer.<token>, zeroclaw.v1 (same as the web dashboard).
func Dial(host, token string, useTLS bool) (*websocket.Conn, *http.Response, error) {
	u := BuildURL(host, useTLS)
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		Subprotocols:     []string{"bearer." + token, subprotocolZeroclaw},
	}
	return dialer.Dial(u, nil)
}
