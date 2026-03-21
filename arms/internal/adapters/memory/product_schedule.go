package memory

import (
	"context"
	"sync"

	"github.com/closeloopautomous/arms/internal/domain"
	"github.com/closeloopautomous/arms/internal/ports"
)

type ProductScheduleStore struct {
	mu   sync.Mutex
	rows map[domain.ProductID]domain.ProductSchedule
}

func NewProductScheduleStore() *ProductScheduleStore {
	return &ProductScheduleStore{rows: make(map[domain.ProductID]domain.ProductSchedule)}
}

var _ ports.ProductScheduleRepository = (*ProductScheduleStore)(nil)

func (s *ProductScheduleStore) Get(_ context.Context, productID domain.ProductID) (*domain.ProductSchedule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.rows[productID]
	if !ok {
		return nil, nil
	}
	return &v, nil
}

func (s *ProductScheduleStore) Upsert(_ context.Context, row *domain.ProductSchedule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	row.UpdatedAt = row.UpdatedAt.UTC()
	s.rows[row.ProductID] = *row
	return nil
}
