package repository

import (
	"context"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/apikeygroup"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type apiKeyGroupRepository struct {
	client *dbent.Client
}

func NewAPIKeyGroupRepository(client *dbent.Client) service.APIKeyGroupRepository {
	return &apiKeyGroupRepository{client: client}
}

func (r *apiKeyGroupRepository) SetGroups(ctx context.Context, apiKeyID int64, groupIDs []int64) error {
	client := clientFromContext(ctx, r.client)

	// Delete existing associations
	if _, err := client.APIKeyGroup.Delete().
		Where(apikeygroup.APIKeyIDEQ(apiKeyID)).
		Exec(ctx); err != nil {
		return err
	}

	// Insert new associations
	bulk := make([]*dbent.APIKeyGroupCreate, 0, len(groupIDs))
	for _, gid := range groupIDs {
		bulk = append(bulk, client.APIKeyGroup.Create().
			SetAPIKeyID(apiKeyID).
			SetGroupID(gid))
	}
	if len(bulk) > 0 {
		if _, err := client.APIKeyGroup.CreateBulk(bulk...).Save(ctx); err != nil {
			return translatePersistenceError(err, nil, nil)
		}
	}
	return nil
}

func (r *apiKeyGroupRepository) GetByKeyID(ctx context.Context, apiKeyID int64) ([]service.APIKeyGroup, error) {
	items, err := r.client.APIKeyGroup.Query().
		Where(apikeygroup.APIKeyIDEQ(apiKeyID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return apiKeyGroupEntitiesToService(items), nil
}

func (r *apiKeyGroupRepository) GetByGroupID(ctx context.Context, groupID int64) ([]service.APIKeyGroup, error) {
	items, err := r.client.APIKeyGroup.Query().
		Where(apikeygroup.GroupIDEQ(groupID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return apiKeyGroupEntitiesToService(items), nil
}

func (r *apiKeyGroupRepository) DeleteByKeyID(ctx context.Context, apiKeyID int64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.APIKeyGroup.Delete().
		Where(apikeygroup.APIKeyIDEQ(apiKeyID)).
		Exec(ctx)
	return err
}

// --- entity ↔ service mapping ---

func apiKeyGroupEntityToService(m *dbent.APIKeyGroup) *service.APIKeyGroup {
	if m == nil {
		return nil
	}
	return &service.APIKeyGroup{
		ID:        m.ID,
		APIKeyID:  m.APIKeyID,
		GroupID:   m.GroupID,
		CreatedAt: m.CreatedAt,
	}
}

func apiKeyGroupEntitiesToService(models []*dbent.APIKeyGroup) []service.APIKeyGroup {
	out := make([]service.APIKeyGroup, 0, len(models))
	for i := range models {
		if s := apiKeyGroupEntityToService(models[i]); s != nil {
			out = append(out, *s)
		}
	}
	return out
}
