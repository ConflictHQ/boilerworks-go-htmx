package model

import (
	"testing"
)

func TestNewPagination(t *testing.T) {
	t.Run("basic pagination", func(t *testing.T) {
		p := NewPagination(1, 20, 50)
		if p.Page != 1 {
			t.Errorf("expected page 1, got %d", p.Page)
		}
		if p.PerPage != 20 {
			t.Errorf("expected perPage 20, got %d", p.PerPage)
		}
		if p.Total != 50 {
			t.Errorf("expected total 50, got %d", p.Total)
		}
		if p.TotalPages != 3 {
			t.Errorf("expected 3 total pages, got %d", p.TotalPages)
		}
	})

	t.Run("exact fit", func(t *testing.T) {
		p := NewPagination(1, 10, 30)
		if p.TotalPages != 3 {
			t.Errorf("expected 3 total pages, got %d", p.TotalPages)
		}
	})

	t.Run("partial last page", func(t *testing.T) {
		p := NewPagination(1, 10, 25)
		if p.TotalPages != 3 {
			t.Errorf("expected 3 total pages, got %d", p.TotalPages)
		}
	})

	t.Run("zero total", func(t *testing.T) {
		p := NewPagination(1, 20, 0)
		if p.TotalPages != 0 {
			t.Errorf("expected 0 total pages, got %d", p.TotalPages)
		}
	})

	t.Run("invalid page defaults to 1", func(t *testing.T) {
		p := NewPagination(0, 20, 50)
		if p.Page != 1 {
			t.Errorf("expected page 1, got %d", p.Page)
		}
	})

	t.Run("invalid perPage defaults to 20", func(t *testing.T) {
		p := NewPagination(1, 0, 50)
		if p.PerPage != 20 {
			t.Errorf("expected perPage 20, got %d", p.PerPage)
		}
	})
}

func TestNewPaginationSingleItem(t *testing.T) {
	p := NewPagination(1, 20, 1)
	if p.TotalPages != 1 {
		t.Errorf("expected 1 total page, got %d", p.TotalPages)
	}
}
