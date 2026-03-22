package task

import (
	"context"
	"fmt"
	"strings"

	"github.com/closeloopautomous/arms/internal/domain"
)

// ApplyAgentWebhookOutcome handles a verified agent-completion webhook: optional Kanban advance for
// full_auto / semi_auto (testing / review) or mark done (default). source is stored on task_completed (done path only).
func (s *Service) ApplyAgentWebhookOutcome(ctx context.Context, taskID domain.TaskID, nextBoardStatus string, source string) error {
	nextBoardStatus = strings.TrimSpace(nextBoardStatus)
	if nextBoardStatus == "" || strings.EqualFold(nextBoardStatus, string(domain.StatusDone)) {
		return s.CompleteWithLiveActivity(ctx, taskID, source)
	}
	switch strings.ToLower(nextBoardStatus) {
	case "testing", "review":
	default:
		return fmt.Errorf("%w: next_board_status must be testing, review, or done (or omit for done)", domain.ErrInvalidInput)
	}
	to, err := domain.ParseTaskStatus(nextBoardStatus)
	if err != nil {
		return fmt.Errorf("%w: next_board_status", domain.ErrInvalidInput)
	}
	t, err := s.Tasks.ByID(ctx, taskID)
	if err != nil {
		return err
	}
	p, err := s.Products.ByID(ctx, t.ProductID)
	if err != nil {
		return err
	}
	if p.AutomationTier != domain.TierFullAuto && p.AutomationTier != domain.TierSemiAuto {
		return s.CompleteWithLiveActivity(ctx, taskID, source)
	}
	if !domain.AllowedKanbanTransition(t.Status, to) {
		return fmt.Errorf("%w: webhook cannot move %s -> %s", domain.ErrInvalidTransition, t.Status, to)
	}
	move := func() error {
		return s.SetKanbanStatus(ctx, taskID, to, "")
	}
	if s.Gate != nil {
		return s.Gate.WithLock(t.ProductID, move)
	}
	return move()
}
