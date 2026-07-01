package repository

import (
	"context"
	dbsql "database/sql"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enterprise"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type enterpriseRepository struct {
	client *dbent.Client
	sqlDB  *dbsql.DB
}

func NewEnterpriseRepository(client *dbent.Client, sqlDB *dbsql.DB) service.EnterpriseRepository {
	return &enterpriseRepository{client: client, sqlDB: sqlDB}
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
	case "balance":
		if isDesc {
			return enterprise.ByBalance(sql.OrderDesc())
		}
		return enterprise.ByBalance(sql.OrderAsc())
	case "concurrency":
		if isDesc {
			return enterprise.ByConcurrency(sql.OrderDesc())
		}
		return enterprise.ByConcurrency(sql.OrderAsc())
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
		SetBalance(e.Balance).
		SetTotalRecharged(e.TotalRecharged).
		SetConcurrency(e.Concurrency).
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
		Concurrency:    m.Concurrency,
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

// CreateBalanceLog 写入企业余额变更审计日志
func (r *enterpriseRepository) CreateBalanceLog(ctx context.Context, enterpriseID int64, amount float64, operation string, notes string) error {
	_, err := r.sqlDB.ExecContext(ctx,
		`INSERT INTO enterprise_balance_logs (enterprise_id, amount, operation, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		enterpriseID, amount, operation, notes, time.Now())
	return err
}

// GetBalanceLogs 查询企业余额变更历史（分页，按时间倒序）
func (r *enterpriseRepository) GetBalanceLogs(ctx context.Context, enterpriseID int64, page, pageSize int) ([]service.EnterpriseBalanceHistoryItem, int64, error) {
	// 总数
	var total int64
	if err := r.sqlDB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM enterprise_balance_logs WHERE enterprise_id = $1`, enterpriseID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := r.sqlDB.QueryContext(ctx,
		`SELECT id, amount, operation, notes, created_at
		 FROM enterprise_balance_logs WHERE enterprise_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		enterpriseID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []service.EnterpriseBalanceHistoryItem
	for rows.Next() {
		var item service.EnterpriseBalanceHistoryItem
		var note dbsql.NullString
		if err := rows.Scan(&item.ID, &item.Amount, &item.Operation, &note, &item.CreatedAt); err != nil {
			return nil, 0, err
		}
		if note.Valid {
			item.Notes = note.String
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

// GetPlatformQuotas 获取企业平台限额列表
func (r *enterpriseRepository) GetPlatformQuotas(ctx context.Context, enterpriseID int64) ([]service.EnterprisePlatformQuota, error) {
	rows, err := r.sqlDB.QueryContext(ctx,
		`SELECT id, enterprise_id, platform, daily_limit_usd, weekly_limit_usd, monthly_limit_usd
		 FROM enterprise_platform_quotas WHERE enterprise_id = $1 ORDER BY platform`, enterpriseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []service.EnterprisePlatformQuota
	for rows.Next() {
		var q service.EnterprisePlatformQuota
		if err := rows.Scan(&q.ID, &q.EnterpriseID, &q.Platform, &q.DailyLimit, &q.WeeklyLimit, &q.MonthlyLimit); err != nil {
			return nil, err
		}
		items = append(items, q)
	}
	return items, rows.Err()
}

// UpsertPlatformQuotas 全量更新企业平台限额
func (r *enterpriseRepository) UpsertPlatformQuotas(ctx context.Context, enterpriseID int64, quotas []service.EnterprisePlatformQuotaInput) error {
	for _, q := range quotas {
		_, err := r.sqlDB.ExecContext(ctx,
			`INSERT INTO enterprise_platform_quotas (enterprise_id, platform, daily_limit_usd, weekly_limit_usd, monthly_limit_usd, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (enterprise_id, platform) DO UPDATE SET
			   daily_limit_usd = EXCLUDED.daily_limit_usd,
			   weekly_limit_usd = EXCLUDED.weekly_limit_usd,
			   monthly_limit_usd = EXCLUDED.monthly_limit_usd,
			   updated_at = EXCLUDED.updated_at`,
			enterpriseID, q.Platform, q.DailyLimitUSD, q.WeeklyLimitUSD, q.MonthlyLimitUSD, time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}
