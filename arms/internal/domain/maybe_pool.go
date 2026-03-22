package domain

import "time"

// MaybePoolEntry is one idea in the deferred "maybe" pool with re-evaluation metadata.
type MaybePoolEntry struct {
	IdeaID            IdeaID
	ProductID         ProductID
	CreatedAt         time.Time
	LastEvaluatedAt   time.Time // zero if never batch-reeval'd
	NextEvaluateAt    time.Time // zero if not scheduled
	EvaluationCount   int
	EvaluationNotes   string
}
