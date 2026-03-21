package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/closeloopautomous/arms/internal/domain"
	"github.com/closeloopautomous/arms/internal/ports"
)

type ConvoyMailStore struct{ db *sql.DB }

func NewConvoyMailStore(db *sql.DB) *ConvoyMailStore { return &ConvoyMailStore{db: db} }

var _ ports.ConvoyMailRepository = (*ConvoyMailStore)(nil)

func (s *ConvoyMailStore) Append(ctx context.Context, convoyID domain.ConvoyID, subtaskID domain.SubtaskID, body string, at time.Time) error {
	id := uuid.NewString()
	atStr := at.UTC().Format(time.RFC3339Nano)
	_, err := s.db.ExecContext(ctx, `
INSERT INTO convoy_mail (id, convoy_id, subtask_id, body, created_at) VALUES (?, ?, ?, ?, ?)`,
		id, string(convoyID), string(subtaskID), body, atStr)
	return err
}

func (s *ConvoyMailStore) ListByConvoy(ctx context.Context, convoyID domain.ConvoyID, limit int) ([]domain.ConvoyMailMessage, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	rows, err := s.db.QueryContext(ctx, `
SELECT id, convoy_id, subtask_id, body, created_at FROM convoy_mail
WHERE convoy_id = ? ORDER BY created_at DESC LIMIT ?`, string(convoyID), limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var out []domain.ConvoyMailMessage
	for rows.Next() {
		var m domain.ConvoyMailMessage
		var atStr string
		if err := rows.Scan(&m.ID, &m.ConvoyID, &m.SubtaskID, &m.Body, &atStr); err != nil {
			return nil, err
		}
		t, perr := time.Parse(time.RFC3339Nano, atStr)
		if perr != nil {
			t, _ = time.Parse(time.RFC3339, atStr)
		}
		m.CreatedAt = t.UTC()
		out = append(out, m)
	}
	return out, rows.Err()
}
