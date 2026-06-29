package repository

import (
	"context"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enterprise"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type enterpriseRepository struct {
	client *dbent.Client
}

func NewEnterpriseRepository(client *dbent.Client) service.EnterpriseRepository {
	return &enterpriseRepository{client: client}
}

// activeQuery returns the base query with soft-delete filter.
func (r *enterpriseRepository) activeQuery() *dbent.EnterpriseQuery {
	return r.client.Enterprise.Query().Where(enterprise.DeletedAtIsNil())
}

func (r *enterpriseRepository) Create(ctx context.Context, e *service.Enterprise) error {
	client := clientFromContext(ctx, r.client)
	created, err := client.Enterprise.Create().
		SetName(e.Name).
		SetShortName(e.ShortName).
		SetCreditCode(e.CreditCode).
		SetAddress(e.Address).
		SetScale(e.Scale).
		SetIndustry(e.Industry).
		SetParentID(e.ParentID).
		SetStatus(e.Status).
		SetContactName(e.ContactName).
		SetContactPhone(e.ContactPhone).
		SetContactEmail(e.ContactEmail).
		SetNotes(e.Notes).
		SetBalance(e.Balance).
		SetTotalRecharged(e.TotalRecharged).
		SetAdminUserID(e.AdminUserID).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrEnterpriseNotFound, service.ErrEnterpriseNameConflict)
	}
	e.ID = created.ID
	e.CreatedAt = created.CreatedAt
	e.UpdatedAt = created.UpdatedAt
	return nil
}

func (r *enterpriseRepository) GetByID(ctx context.Context, id int64) (*service.Enterprise, error) {
	m, err := r.activeQuery().
		Where(enterprise.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrEnterpriseNotFound, nil)
	}
	return enterpriseEntityToService(m), nil
}

func (r *enterpriseRepository) List(
	ctx context.Context,
	params pagination.PaginationParams,
	filters service.EnterpriseListFilters,
) ([]service.Enterprise, *pagination.PaginationResult, error) {
	q := r.activeQuery()

	if filters.Status != "" {
		q = q.Where(enterprise.StatusEQ(filters.Status))
	}
	if filters.Search != "" {
		q = q.Where(
			enterprise.Or(
				enterprise.NameContainsFold(filters.Search),
				enterprise.ShortNameContainsFold(filters.Search),
				enterprise.ContactNameContainsFold(filters.Search),
				enterprise.ContactEmailContainsFold(filters.Search),
			),
		)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	items, err := q.
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(enterpriseListOrder(params)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	out := enterpriseEntitiesToService(items)
	return out, paginationResultFromTotal(int64(total), params), nil
}

func enterpriseListOrder(params pagination.PaginationParams) enterprise.OrderOption {
	sortBy := strings.ToLower(strings.TrimSpace(params.SortBy))
	sortOrder := params.NormalizedSortOrder(pagination.SortOrderDesc)

	isDesc := sortOrder == pagination.SortOrderDesc

	switch sortBy {
	case "name":
		if isDesc {
			return enterprise.ByName(sql.OrderDesc())
		}
		return enterprise.ByName(sql.OrderAsc())
	case "status":
		if isDesc {
			return enterprise.ByStatus(sql.OrderDesc())
		}
		return enterprise.ByStatus(sql.OrderAsc())
	case "id":
		if isDesc {
			return enterprise.ByID(sql.OrderDesc())
		}
		return enterprise.ByID(sql.OrderAsc())
	case "", "created_at":
		fallthrough
	default:
		if isDesc {
			return enterprise.ByCreatedAt(sql.OrderDesc())
		}
		return enterprise.ByCreatedAt(sql.OrderAsc())
	}
}

func (r *enterpriseRepository) Update(ctx context.Context, e *service.Enterprise) error {
	client := clientFromContext(ctx, r.client)
	updated, err := client.Enterprise.UpdateOneID(e.ID).
		SetName(e.Name).
		SetShortName(e.ShortName).
		SetCreditCode(e.CreditCode).
		SetAddress(e.Address).
		SetScale(e.Scale).
		SetIndustry(e.Industry).
		SetParentID(e.ParentID).
		SetStatus(e.Status).
		SetContactName(e.ContactName).
		SetContactPhone(e.ContactPhone).
		SetContactEmail(e.ContactEmail).
		SetNotes(e.Notes).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrEnterpriseNotFound, service.ErrEnterpriseNameConflict)
	}
	e.UpdatedAt = updated.UpdatedAt
	return nil
}

func (r *enterpriseRepository) SoftDelete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	return client.Enterprise.DeleteOneID(id).Exec(ctx)
}

func (r *enterpriseRepository) GetBalance(ctx context.Context, id int64) (float64, error) {
	m, err := r.activeQuery().
		Where(enterprise.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return 0, translatePersistenceError(err, service.ErrEnterpriseNotFound, nil)
	}
	return m.Balance, nil
}

// DeductBalance 原子扣减企业余额，返回扣减后的余额。
// 使用 WHERE balance >= amount 防止超扣。
func (r *enterpriseRepository) DeductBalance(ctx context.Context, id int64, amount float64) (float64, error) {
	client := clientFromContext(ctx, r.client)
	n, err := client.Enterprise.Update().
		Where(
			enterprise.IDEQ(id),
			enterprise.DeletedAtIsNil(),
			enterprise.BalanceGTE(amount),
		).
		AddBalance(-amount).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("deduct enterprise balance: %w", err)
	}
	if n == 0 {
		return 0, service.ErrEnterpriseInsufficientBalance
	}
	// 查询扣减后的余额
	bal, err := r.GetBalance(ctx, id)
	if err != nil {
		return 0, err
	}
	return bal, nil
}

// --- entity ↔ service mapping ---

func enterpriseEntityToService(m *dbent.Enterprise) *service.Enterprise {
	if m == nil {
		return nil
	}
	return &service.Enterprise{
		ID:             m.ID,
		Name:           m.Name,
		ShortName:      m.ShortName,
		CreditCode:     m.CreditCode,
		Address:        m.Address,
		Scale:          m.Scale,
		Industry:       m.Industry,
		ParentID:       m.ParentID,
		Status:         m.Status,
		ContactName:    m.ContactName,
		ContactPhone:   m.ContactPhone,
		ContactEmail:   m.ContactEmail,
		Notes:          m.Notes,
		Balance:        m.Balance,
		TotalRecharged: m.TotalRecharged,
		AdminUserID:    m.AdminUserID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func enterpriseEntitiesToService(models []*dbent.Enterprise) []service.Enterprise {
	out := make([]service.Enterprise, 0, len(models))
	for i := range models {
		if s := enterpriseEntityToService(models[i]); s != nil {
			out = append(out, *s)
		}
	}
	return out
}
