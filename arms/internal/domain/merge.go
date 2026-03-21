package domain

import (
	"encoding/json"
	"strings"
)

// MergeShipState is persisted on workspace_merge_queue after a ship attempt.
type MergeShipState string

const (
	MergeShipNone     MergeShipState = ""          // never attempted / legacy row
	MergeShipMerged   MergeShipState = "merged"    // Git or local merge succeeded
	MergeShipSkipped  MergeShipState = "skipped"   // break-glass: advance queue without forge merge
	MergeShipConflict MergeShipState = "conflict"  // merge conflict (local or GitHub)
	MergeShipFailed   MergeShipState = "failed"    // other failure
)

// MergeShipResult is the outcome of a forge or local git merge attempt.
type MergeShipResult struct {
	State         MergeShipState
	MergedSHA     string
	ErrorMessage  string
	ConflictFiles []string
}

// MergePolicy is optional per-product JSON (merge_policy_json).
type MergePolicy struct {
	MergeMethod string `json:"merge_method"` // merge | squash | rebase (default merge)
	// MergeBackendOverride: when set, overrides process ARMS_MERGE_BACKEND for this product (github|local|noop).
	MergeBackendOverride string `json:"merge_backend,omitempty"`
	// RequireApprovedReview, when non-nil, overrides the default implied by [AutomationTier] for auto merge-queue ship.
	RequireApprovedReview *bool `json:"require_approved_review,omitempty"`
	// RequireCleanMergeable, when non-nil, overrides the default (GitHub mergeable_state must be "clean").
	RequireCleanMergeable *bool `json:"require_clean_mergeable,omitempty"`
}

// MergeExecutionGates controls unattended merge-queue completion (semi_auto / policy), not manual POST .../complete.
type MergeExecutionGates struct {
	RequireApprovedReview bool
	RequireCleanMergeable bool
}

// EffectiveMergeExecutionGates combines automation tier defaults with merge_policy_json overrides.
func EffectiveMergeExecutionGates(p *Product, pol MergePolicy) MergeExecutionGates {
	var defRev, defClean bool
	switch p.AutomationTier {
	case TierFullAuto:
		defRev, defClean = false, false
	case TierSemiAuto:
		defRev, defClean = true, true
	default:
		defRev, defClean = true, true
	}
	rev := defRev
	if pol.RequireApprovedReview != nil {
		rev = *pol.RequireApprovedReview
	}
	clean := defClean
	if pol.RequireCleanMergeable != nil {
		clean = *pol.RequireCleanMergeable
	}
	return MergeExecutionGates{RequireApprovedReview: rev, RequireCleanMergeable: clean}
}

// ParseMergePolicy unmarshals product.MergePolicyJSON; invalid JSON yields defaults.
func ParseMergePolicy(jsonStr string) MergePolicy {
	s := strings.TrimSpace(jsonStr)
	if s == "" {
		return MergePolicy{MergeMethod: "merge"}
	}
	var p MergePolicy
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return MergePolicy{MergeMethod: "merge"}
	}
	if strings.TrimSpace(p.MergeMethod) == "" {
		p.MergeMethod = "merge"
	}
	p.MergeMethod = strings.ToLower(strings.TrimSpace(p.MergeMethod))
	p.MergeBackendOverride = strings.ToLower(strings.TrimSpace(p.MergeBackendOverride))
	return p
}

// NormalizeMergeMethod returns a GitHub API merge_method string.
func NormalizeMergeMethod(m string) string {
	switch strings.ToLower(strings.TrimSpace(m)) {
	case "squash":
		return "squash"
	case "rebase":
		return "rebase"
	default:
		return "merge"
	}
}
