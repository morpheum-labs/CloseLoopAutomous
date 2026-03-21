package domain

import (
	"testing"
)

func TestEffectiveMergeExecutionGates(t *testing.T) {
	pol := ParseMergePolicy(`{"merge_method":"squash"}`)
	p := &Product{AutomationTier: TierSemiAuto}
	g := EffectiveMergeExecutionGates(p, pol)
	if !g.RequireApprovedReview || !g.RequireCleanMergeable {
		t.Fatalf("semi_auto defaults: %#v", g)
	}

	p.AutomationTier = TierFullAuto
	g = EffectiveMergeExecutionGates(p, pol)
	if g.RequireApprovedReview || g.RequireCleanMergeable {
		t.Fatalf("full_auto defaults: %#v", g)
	}

	pol2 := ParseMergePolicy(`{"require_approved_review":true,"require_clean_mergeable":true}`)
	p.AutomationTier = TierFullAuto
	g = EffectiveMergeExecutionGates(p, pol2)
	if !g.RequireApprovedReview || !g.RequireCleanMergeable {
		t.Fatalf("json override on full_auto: %#v", g)
	}

	f := false
	pol3 := MergePolicy{MergeMethod: "merge", RequireApprovedReview: &f, RequireCleanMergeable: &f}
	p.AutomationTier = TierSemiAuto
	g = EffectiveMergeExecutionGates(p, pol3)
	if g.RequireApprovedReview || g.RequireCleanMergeable {
		t.Fatalf("json false override on semi: %#v", g)
	}
}
