package service

import (
	"context"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// DepartmentListFilters defines list query filters for departments.
type DepartmentListFilters struct {
	Status string
}

// DepartmentRepository interface
type DepartmentRepository interface {
	Create(ctx context.Context, d *Department) error
	GetByID(ctx context.Context, id int64) (*Department, error)
	TreeByEnterprise(ctx context.Context, enterpriseID int64) ([]Department, error)
	ListByEnterprise(ctx context.Context, enterpriseID int64, params pagination.PaginationParams, filters DepartmentListFilters) ([]Department, *pagination.PaginationResult, error)
	Update(ctx context.Context, d *Department) error
	Delete(ctx context.Context, id int64) error
	HasChildren(ctx context.Context, id int64) (bool, error)
	HasMembers(ctx context.Context, id int64) (bool, error)
}
