package domain

import "time"

// OperationLogEntry is one append-only audit row (Mission Control–style operations_log).
type OperationLogEntry struct {
	ID           int64
	CreatedAt    time.Time
	Actor        string
	Action       string
	ResourceType string
	ResourceID   string
	DetailJSON   string
	ProductID    ProductID
}
