package domain

import (
	"strings"
	"time"
)

// ProductFeedback is external feedback attached to a product (MC product_feedback).
type ProductFeedback struct {
	ID          string
	ProductID   ProductID
	Source      string
	Content     string
	CustomerID  string
	Category    string
	Sentiment   string // positive, negative, neutral, mixed
	Processed   bool
	IdeaID      IdeaID // optional link to a promoted idea
	CreatedAt   time.Time
}

// NormalizeFeedbackSentiment maps to MC sentiment values.
func NormalizeFeedbackSentiment(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "positive", "negative", "neutral", "mixed":
		return strings.ToLower(strings.TrimSpace(s))
	default:
		return "neutral"
	}
}
