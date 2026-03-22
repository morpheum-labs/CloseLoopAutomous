package domain

import (
	"testing"
	"time"
)

func TestProductIsDeletedMarkClear(t *testing.T) {
	var p Product
	if p.IsDeleted() {
		t.Fatal("zero product should not be deleted")
	}
	at := time.Unix(1700000000, 0).UTC()
	p.MarkDeleted(at)
	if !p.IsDeleted() || !p.DeletedAt.Equal(at) || !p.UpdatedAt.Equal(at) {
		t.Fatalf("MarkDeleted: DeletedAt=%v UpdatedAt=%v", p.DeletedAt, p.UpdatedAt)
	}
	later := at.Add(time.Hour)
	p.ClearDeletion(later)
	if p.IsDeleted() || !p.UpdatedAt.Equal(later) {
		t.Fatalf("ClearDeletion: DeletedAt=%v UpdatedAt=%v", p.DeletedAt, p.UpdatedAt)
	}
}
