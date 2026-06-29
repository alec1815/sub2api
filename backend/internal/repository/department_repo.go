package repository

import (
	"context"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/department"
	"github.com/Wei-Shaw/sub2api/ent/enterprisemember"
	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type departmentRepository struct {
	client *dbent.Client
}

func NewDepartmentRepository(client *dbent.Client) service.DepartmentRepository {
	return &departmentRepository{client: client}
}

func (r *departmentRepository) activeQuery() *dbent.DepartmentQuery {
	return r.client.Department.Query().Where(department.DeletedAtIsNil())
}

func (r *departmentRepository) Create(ctx context.Context, d *service.Department) error {
	client := clientFromContext(ctx, r.client)
	created, err := client.Department.Create().
		SetEnterpriseID(d.EnterpriseID).
		SetParentID(d.ParentID).
		SetName(d.Name).
		SetOrderNum(d.OrderNum).
		SetLeader(d.Leader).
		SetPhone(d.Phone).
		SetEmail(d.Email).
		SetStatus(d.Status).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrDepartmentNotFound, service.ErrDepartmentNameExist)
	}
	d.ID = created.ID
	d.CreatedAt = created.CreatedAt
	d.UpdatedAt = created.UpdatedAt
	return nil
}

func (r *departmentRepository) GetByID(ctx context.Context, id int64) (*service.Department, error) {
	m, err := r.activeQuery().
		Where(department.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrDepartmentNotFound, nil)
	}
	return departmentEntityToService(m), nil
}

func (r *departmentRepository) TreeByEnterprise(ctx context.Context, enterpriseID int64) ([]service.Department, error) {
	items, err := r.activeQuery().
		Where(department.EnterpriseIDEQ(enterpriseID)).
		Where(department.StatusEQ(domain.StatusActive)).
		Order(dbent.Asc(department.FieldOrderNum), dbent.Asc(department.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return departmentEntitiesToService(items), nil
}

func (r *departmentRepository) ListByEnterprise(
	ctx context.Context,
	enterpriseID int64,
	params pagination.PaginationParams,
	filters service.DepartmentListFilters,
) ([]service.Department, *pagination.PaginationResult, error) {
	q := r.activeQuery().
		Where(department.EnterpriseIDEQ(enterpriseID))

	if filters.Status != "" {
		q = q.Where(department.StatusEQ(filters.Status))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	items, err := q.
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Asc(department.FieldOrderNum), dbent.Asc(department.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	out := departmentEntitiesToService(items)
	return out, paginationResultFromTotal(int64(total), params), nil
}

func (r *departmentRepository) Update(ctx context.Context, d *service.Department) error {
	client := clientFromContext(ctx, r.client)
	updated, err := client.Department.UpdateOneID(d.ID).
		SetParentID(d.ParentID).
		SetName(d.Name).
		SetOrderNum(d.OrderNum).
		SetLeader(d.Leader).
		SetPhone(d.Phone).
		SetEmail(d.Email).
		SetStatus(d.Status).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrDepartmentNotFound, service.ErrDepartmentNameExist)
	}
	d.UpdatedAt = updated.UpdatedAt
	return nil
}

func (r *departmentRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	return client.Department.DeleteOneID(id).Exec(ctx)
}

func (r *departmentRepository) HasChildren(ctx context.Context, id int64) (bool, error) {
	count, err := r.activeQuery().
		Where(department.ParentIDEQ(id)).
		Count(ctx)
	return count > 0, err
}

func (r *departmentRepository) HasMembers(ctx context.Context, id int64) (bool, error) {
	count, err := r.client.EnterpriseMember.Query().
		Where(
			enterprisemember.DepartmentIDEQ(id),
			enterprisemember.StatusEQ(domain.StatusActive),
			enterprisemember.DeletedAtIsNil(),
		).
		Count(ctx)
	return count > 0, err
}

// --- entity ↔ service mapping ---

func departmentEntityToService(m *dbent.Department) *service.Department {
	if m == nil {
		return nil
	}
	return &service.Department{
		ID:           m.ID,
		EnterpriseID: m.EnterpriseID,
		ParentID:     m.ParentID,
		Name:         m.Name,
		OrderNum:     m.OrderNum,
		Leader:       m.Leader,
		Phone:        m.Phone,
		Email:        m.Email,
		Status:       m.Status,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func departmentEntitiesToService(models []*dbent.Department) []service.Department {
	out := make([]service.Department, 0, len(models))
	for i := range models {
		if s := departmentEntityToService(models[i]); s != nil {
			out = append(out, *s)
		}
	}
	return out
}
