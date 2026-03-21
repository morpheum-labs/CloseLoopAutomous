package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/closeloopautomous/arms/internal/domain"
	"github.com/closeloopautomous/arms/internal/ports"
)

type ProductScheduleStore struct{ db *sql.DB }

func NewProductScheduleStore(db *sql.DB) *ProductScheduleStore { return &ProductScheduleStore{db: db} }

var _ ports.ProductScheduleRepository = (*ProductScheduleStore)(nil)

func (s *ProductScheduleStore) Get(ctx context.Context, productID domain.ProductID) (*domain.ProductSchedule, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT product_id, enabled, spec_json, updated_at FROM product_schedules WHERE product_id = ?`, string(productID))
	var pid, spec, atStr string
	var en int
	if err := row.Scan(&pid, &en, &spec, &atStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	t, perr := time.Parse(time.RFC3339Nano, atStr)
	if perr != nil {
		t, _ = time.Parse(time.RFC3339, atStr)
	}
	return &domain.ProductSchedule{
		ProductID: domain.ProductID(pid),
		Enabled:   en != 0,
		SpecJSON:  spec,
		UpdatedAt: t.UTC(),
	}, nil
}

func (s *ProductScheduleStore) Upsert(ctx context.Context, row *domain.ProductSchedule) error {
	atStr := row.UpdatedAt.UTC().Format(time.RFC3339Nano)
	en := 0
	if row.Enabled {
		en = 1
	}
	spec := row.SpecJSON
	if spec == "" {
		spec = "{}"
	}
	_, err := s.db.ExecContext(ctx, `
INSERT INTO product_schedules (product_id, enabled, spec_json, updated_at) VALUES (?, ?, ?, ?)
ON CONFLICT(product_id) DO UPDATE SET enabled = excluded.enabled, spec_json = excluded.spec_json, updated_at = excluded.updated_at
`, string(row.ProductID), en, spec, atStr)
	return err
}
