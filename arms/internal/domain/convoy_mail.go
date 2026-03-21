package domain

import "time"

// ConvoyMailMessage is an append-only note attached to a convoy subtask (inter-subtask messaging baseline).
type ConvoyMailMessage struct {
	ID        string
	ConvoyID  ConvoyID
	SubtaskID SubtaskID
	Body      string
	CreatedAt time.Time
}
