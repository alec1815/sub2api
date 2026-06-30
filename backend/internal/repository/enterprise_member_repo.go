package repository

import (
	"context"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enterprisemember"
	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type enterpriseMemberRepository struct {
	client *dbent.Client
}

func NewEnterpriseMemberRepository(client *dbent.Client) service.EnterpriseMemberRepository {
	return &enterpriseMemberRepository{client: client}
}

func (r *enterpriseMemberRepository) activeQuery() *dbent.EnterpriseMemberQuery {
	return r.client.EnterpriseMember.Query().Where(enterprisemember.DeletedAtIsNil())
}

func (r *enterpriseMemberRepository) Create(ctx context.Context, m *service.EnterpriseMember) error {
	client := clientFromContext(ctx, r.client)
	created, err := client.EnterpriseMember.Create().
		SetEnterpriseID(m.EnterpriseID).
		SetUserID(m.UserID).
		SetRole(m.Role).
		SetStatus(m.Status).
		SetNillableDepartmentID(m.DepartmentID).
		SetConcurrency(m.Concurrency).
		SetRpmLimit(m.RPMLimit).
		SetNotes(m.Notes).
		SetJoinedAt(m.JoinedAt).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrEnterpriseMemberNotFound, service.ErrEnterpriseMemberAlreadyActive)
	}
	m.ID = created.ID
	m.CreatedAt = created.CreatedAt
	m.UpdatedAt = created.UpdatedAt
	return nil
}

func (r *enterpriseMemberRepository) GetByID(ctx context.Context, id int64) (*service.EnterpriseMember, error) {
	m, err := r.activeQuery().
		Where(enterprisemember.IDEQ(id)).
		WithUser().
		WithDepartment().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrEnterpriseMemberNotFound, nil)
	}
	return enterpriseMemberEntityToService(m), nil
}

func (r *enterpriseMemberRepository) GetByUserID(ctx context.Context, userID int64) (*service.EnterpriseMember, error) {
	m, err := r.activeQuery().
		Where(
			enterprisemember.UserIDEQ(userID),
			enterprisemember.StatusEQ(domain.StatusActive),
		).
		WithUser().
		WithDepartment().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrEnterpriseMemberNotFound, nil)
	}
	return enterpriseMemberEntityToService(m), nil
}

func (r *enterpriseMemberRepository) ListByEnterprise(
	ctx context.Context,
	enterpriseID int64,
	params pagination.PaginationParams,
	filters service.EnterpriseMemberListFilters,
) ([]service.EnterpriseMember, *pagination.PaginationResult, error) {
	q := r.activeQuery().
		Where(enterprisemember.EnterpriseIDEQ(enterpriseID))

	if filters.Role != "" {
		q = q.Where(enterprisemember.RoleEQ(filters.Role))
	}
	if filters.Status != "" {
		q = q.Where(enterprisemember.StatusEQ(filters.Status))
	}
	if filters.DepartmentID != nil {
		q = q.Where(enterprisemember.DepartmentIDEQ(*filters.DepartmentID))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	items, err := q.
		WithUser().
		WithDepartment().
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Desc(enterprisemember.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	out := enterpriseMemberEntitiesToService(items)
	return out, paginationResultFromTotal(int64(total), params), nil
}

func (r *enterpriseMemberRepository) Update(ctx context.Context, m *service.EnterpriseMember) error {
	client := clientFromContext(ctx, r.client)
	builder := client.EnterpriseMember.UpdateOneID(m.ID).
		SetRole(m.Role).
		SetNotes(m.Notes).
		SetConcurrency(m.Concurrency).
		SetRpmLimit(m.RPMLimit)

	if m.DepartmentID != nil {
		builder.SetDepartmentID(*m.DepartmentID)
	} else {
		builder.ClearDepartmentID()
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrEnterpriseMemberNotFound, nil)
	}
	m.UpdatedAt = updated.UpdatedAt
	return nil
}

func (r *enterpriseMemberRepository) Unbind(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	now := time.Now()
	return client.EnterpriseMember.UpdateOneID(id).
		SetStatus(domain.MemberStatusUnbound).
		SetUnboundAt(now).
		Exec(ctx)
}

// --- entity ↔ service mapping ---

func enterpriseMemberEntityToService(m *dbent.EnterpriseMember) *service.EnterpriseMember {
	if m == nil {
		return nil
	}
	s := &service.EnterpriseMember{
		ID:           m.ID,
		EnterpriseID: m.EnterpriseID,
		UserID:       m.UserID,
		Role:         m.Role,
		Status:       m.Status,
		DepartmentID: m.DepartmentID,
		Concurrency:  m.Concurrency,
		RPMLimit:     m.RpmLimit,
		Notes:        m.Notes,
		JoinedAt:     m.JoinedAt,
		UnboundAt:    m.UnboundAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
	// eager-loaded user fields
	if u := m.Edges.User; u != nil {
		s.Name = u.Username
		s.Email = u.Email
		if u.LastActiveAt != nil {
			s.LastActiveAt = u.LastActiveAt
		}
	}
	// eager-loaded department field
	if d := m.Edges.Department; d != nil {
		s.DepartmentName = d.Name
	}
	return s
}

func enterpriseMemberEntitiesToService(models []*dbent.EnterpriseMember) []service.EnterpriseMember {
	out := make([]service.EnterpriseMember, 0, len(models))
	for i := range models {
		if s := enterpriseMemberEntityToService(models[i]); s != nil {
			out = append(out, *s)
		}
	}
	return out
}
