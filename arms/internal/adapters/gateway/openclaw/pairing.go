package openclaw

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coder/websocket"
)

// ErrPairingRequired indicates the gateway closed the WebSocket with policy 1008 or an explicit pairing reason.
var ErrPairingRequired = errors.New("openclaw: device pairing required")

// PairingError carries gateway close metadata for Mission Control and tests.
type PairingError struct {
	RequestID string
	Reason    string
	Detail    string
}

func (e *PairingError) Error() string {
	if e == nil {
		return ""
	}
	return ErrPairingRequired.Error() + ": " + e.Detail
}

// Unwrap supports errors.Is(..., ErrPairingRequired).
func (e *PairingError) Unwrap() error { return ErrPairingRequired }

func extractRequestID(reason string) string {
	r := strings.TrimSpace(reason)
	lower := strings.ToLower(r)
	for _, marker := range []string{"requestid:", "request_id:"} {
		idx := strings.Index(lower, marker)
		if idx < 0 {
			continue
		}
		return strings.TrimSpace(r[idx+len(marker):])
	}
	return ""
}

func pairingApproveDetail(requestID string) string {
	if strings.TrimSpace(requestID) != "" {
		return fmt.Sprintf(
			"OpenClaw device pairing required (WebSocket close 1008 / policy). On the Gateway host run:\n"+
				"  openclaw devices list\n"+
				"  openclaw devices approve %s",
			requestID,
		)
	}
	return "OpenClaw device pairing required (WebSocket close 1008 / policy). On the Gateway host run:\n" +
		"  openclaw devices list\n" +
		"  openclaw devices approve <request-id>"
}

func isPairingClose(code websocket.StatusCode, reason string) bool {
	if strings.Contains(strings.ToLower(reason), "pairing") {
		return true
	}
	return code == websocket.StatusPolicyViolation
}

func closeReasonString(err error) string {
	var ce websocket.CloseError
	if errors.As(err, &ce) {
		return ce.Reason
	}
	return ""
}

func newPairingError(err error) error {
	if err == nil {
		return nil
	}
	code := websocket.CloseStatus(err)
	reason := closeReasonString(err)
	if !isPairingClose(code, reason) {
		return nil
	}
	rid := extractRequestID(reason)
	return &PairingError{
		RequestID: rid,
		Reason:    reason,
		Detail:    pairingApproveDetail(rid),
	}
}

func enrichReadError(err error) error {
	if err == nil {
		return nil
	}
	if pe := newPairingError(err); pe != nil {
		return fmt.Errorf("openclaw read: %w", pe)
	}
	return fmt.Errorf("openclaw read: %w", err)
}
