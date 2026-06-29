package repository

import (
	"context"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/enterprisesubscription"
	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type enterpriseSubscriptionRepository struct {
	client *dbent.Client
}

func NewEnterpriseSubscriptionRepository(client *dbent.Client) service.EnterpriseSubscriptionRepository {
	return &enterpriseSubscriptionRepository{client: client}
}

func (r *enterpriseSubscriptionRepository) Create(ctx context.Context, s *service.EnterpriseSubscription) error {
	client := clientFromContext(ctx, r.client)
	created, err := client.EnterpriseSubscription.Create().
		SetEnterpriseID(s.EnterpriseID).
		SetGroupID(s.GroupID).
		SetPlanID(s.PlanID).
		SetStartsAt(s.StartsAt).
		SetNillableExpiresAt(s.ExpiresAt).
		SetStatus(s.Status).
		SetDailyUsageUsd(s.DailyUsageUSD).
		SetWeeklyUsageUsd(s.WeeklyUsageUSD).
		SetMonthlyUsageUsd(s.MonthlyUsageUSD).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrEnterpriseSubNotFound, nil)
	}
	s.ID = created.ID
	s.CreatedAt = created.CreatedAt
	s.UpdatedAt = created.UpdatedAt
	return nil
}

func (r *enterpriseSubscriptionRepository) GetByID(ctx context.Context, id int64) (*service.EnterpriseSubscription, error) {
	m, err := r.client.EnterpriseSubscription.Query().
		Where(enterprisesubscription.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrEnterpriseSubNotFound, nil)
	}
	return enterpriseSubscriptionEntityToService(m), nil
}

func (r *enterpriseSubscriptionRepository) GetByEnterprise(ctx context.Context, enterpriseID int64) ([]service.EnterpriseSubscription, error) {
	items, err := r.client.EnterpriseSubscription.Query().
		Where(enterprisesubscription.EnterpriseIDEQ(enterpriseID)).
		Order(dbent.Desc(enterprisesubscription.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return enterpriseSubscriptionEntitiesToService(items), nil
}

func (r *enterpriseSubscriptionRepository) ListActive(ctx context.Context, enterpriseID int64) ([]service.EnterpriseSubscription, error) {
	items, err := r.client.EnterpriseSubscription.Query().
		Where(
			enterprisesubscription.EnterpriseIDEQ(enterpriseID),
			enterprisesubscription.StatusEQ(domain.EnterpriseSubStatusActive),
		).
		Order(dbent.Desc(enterprisesubscription.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return enterpriseSubscriptionEntitiesToService(items), nil
}

func (r *enterpriseSubscriptionRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.EnterpriseSubscription.UpdateOneID(id).
		SetStatus(status).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrEnterpriseSubNotFound, nil)
	}
	return nil
}

// --- entity ↔ service mapping ---

func enterpriseSubscriptionEntityToService(m *dbent.EnterpriseSubscription) *service.EnterpriseSubscription {
	if m == nil {
		return nil
	}
	return &service.EnterpriseSubscription{
		ID:              m.ID,
		EnterpriseID:    m.EnterpriseID,
		GroupID:         m.GroupID,
		PlanID:          m.PlanID,
		StartsAt:        m.StartsAt,
		ExpiresAt:       m.ExpiresAt,
		Status:          m.Status,
		DailyUsageUSD:   m.DailyUsageUsd,
		WeeklyUsageUSD:  m.WeeklyUsageUsd,
		MonthlyUsageUSD: m.MonthlyUsageUsd,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func enterpriseSubscriptionEntitiesToService(models []*dbent.EnterpriseSubscription) []service.EnterpriseSubscription {
	out := make([]service.EnterpriseSubscription, 0, len(models))
	for i := range models {
		if s := enterpriseSubscriptionEntityToService(models[i]); s != nil {
			out = append(out, *s)
		}
	}
	return out
}
