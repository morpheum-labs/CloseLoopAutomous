package domain

import "time"

// IsDeleted reports whether the product is soft-deleted (DeletedAt is set).
// All adapters and application code should use this for a single definition of "active" vs tombstone.
func (p *Product) IsDeleted() bool {
	return !p.DeletedAt.IsZero()
}

// MarkDeleted sets DeletedAt and UpdatedAt to the same instant (UTC).
func (p *Product) MarkDeleted(at time.Time) {
	t := at.UTC()
	p.DeletedAt = t
	p.UpdatedAt = t
}

// ClearDeletion clears DeletedAt and sets UpdatedAt (UTC).
func (p *Product) ClearDeletion(at time.Time) {
	p.DeletedAt = time.Time{}
	p.UpdatedAt = at.UTC()
}
