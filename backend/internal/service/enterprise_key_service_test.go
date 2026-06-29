//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Stubs — EnterpriseKeyService
// ---------------------------------------------------------------------------

type entKeyAPIKeyRepoStub struct {
	key         *APIKey
	getByIDErr  error
	createErr   error
	updateErr   error
	deleteErr   error
	listResult  []APIKey
	listPag     *pagination.PaginationResult
	listErr     error
}

func (s *entKeyAPIKeyRepoStub) Create(ctx context.Context, k *APIKey) error {
	if s.createErr != nil { return s.createErr }
	k.ID = 123
	k.Key = "sk-test-abc"
	return nil
}
func (s *entKeyAPIKeyRepoStub) GetByID(ctx context.Context, id int64) (*APIKey, error) {
	if s.getByIDErr != nil { return nil, s.getByIDErr }
	if s.key != nil { return s.key, nil }
	return nil, ErrAPIKeyNotFound
}
func (s *entKeyAPIKeyRepoStub) GetKeyAndOwnerID(ctx context.Context, id int64) (string, int64, error) {
	if s.key != nil { return s.key.Key, s.key.UserID, nil }
	return "", 0, ErrAPIKeyNotFound
}
func (s *entKeyAPIKeyRepoStub) GetByKey(ctx context.Context, key string) (*APIKey, error) {
	return nil, ErrAPIKeyNotFound
}
func (s *entKeyAPIKeyRepoStub) GetByKeyForAuth(ctx context.Context, key string) (*APIKey, error) {
	return nil, ErrAPIKeyNotFound
}
func (s *entKeyAPIKeyRepoStub) Update(ctx context.Context, k *APIKey) error          { return s.updateErr }
func (s *entKeyAPIKeyRepoStub) Delete(ctx context.Context, id int64) error            { return s.deleteErr }
func (s *entKeyAPIKeyRepoStub) DeleteWithAudit(ctx context.Context, id int64) error   { return s.deleteErr }
func (s *entKeyAPIKeyRepoStub) ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams, filters APIKeyListFilters) ([]APIKey, *pagination.PaginationResult, error) {
	return s.listResult, s.listPag, s.listErr
}
func (s *entKeyAPIKeyRepoStub) VerifyOwnership(ctx context.Context, userID int64, apiKeyIDs []int64) ([]int64, error) { return nil, nil }
func (s *entKeyAPIKeyRepoStub) CountByUserID(ctx context.Context, userID int64) (int64, error) { return 0, nil }
func (s *entKeyAPIKeyRepoStub) ExistsByKey(ctx context.Context, key string) (bool, error) { return false, nil }
func (s *entKeyAPIKeyRepoStub) ListByGroupID(ctx context.Context, groupID int64, params pagination.PaginationParams) ([]APIKey, *pagination.PaginationResult, error) { panic("unexpected") }
func (s *entKeyAPIKeyRepoStub) SearchAPIKeys(ctx context.Context, userID int64, keyword string, limit int) ([]APIKey, error) { return nil, nil }
func (s *entKeyAPIKeyRepoStub) ClearGroupIDByGroupID(ctx context.Context, groupID int64) (int64, error) { return 0, nil }
func (s *entKeyAPIKeyRepoStub) UpdateGroupIDByUserAndGroup(ctx context.Context, userID, oldGroupID, newGroupID int64) (int64, error) { return 0, nil }
func (s *entKeyAPIKeyRepoStub) CountByGroupID(ctx context.Context, groupID int64) (int64, error) { return 0, nil }
func (s *entKeyAPIKeyRepoStub) ListKeysByUserID(ctx context.Context, userID int64) ([]string, error) { return nil, nil }
func (s *entKeyAPIKeyRepoStub) ListKeysByGroupID(ctx context.Context, groupID int64) ([]string, error) { return nil, nil }
func (s *entKeyAPIKeyRepoStub) IncrementQuotaUsed(ctx context.Context, id int64, amount float64) (float64, error) { return 0, nil }
func (s *entKeyAPIKeyRepoStub) UpdateLastUsed(ctx context.Context, id int64, usedAt time.Time) error { return nil }
func (s *entKeyAPIKeyRepoStub) IncrementRateLimitUsage(ctx context.Context, id int64, cost float64) error { return nil }
func (s *entKeyAPIKeyRepoStub) ResetRateLimitWindows(ctx context.Context, id int64) error { return nil }
func (s *entKeyAPIKeyRepoStub) GetRateLimitData(ctx context.Context, id int64) (*APIKeyRateLimitData, error) { return nil, nil }

type entKeyEntRepoStub struct {
	enterprise *Enterprise
	getByIDErr error
}

func (s *entKeyEntRepoStub) Create(ctx context.Context, e *Enterprise) error { panic("unexpected") }
func (s *entKeyEntRepoStub) GetByID(ctx context.Context, id int64) (*Enterprise, error) {
	if s.getByIDErr != nil { return nil, s.getByIDErr }
	return s.enterprise, nil
}
func (s *entKeyEntRepoStub) List(ctx context.Context, params pagination.PaginationParams, filters EnterpriseListFilters) ([]Enterprise, *pagination.PaginationResult, error) { panic("unexpected") }
func (s *entKeyEntRepoStub) Update(ctx context.Context, e *Enterprise) error    { panic("unexpected") }
func (s *entKeyEntRepoStub) SoftDelete(ctx context.Context, id int64) error     { panic("unexpected") }
func (s *entKeyEntRepoStub) GetBalance(ctx context.Context, id int64) (float64, error) { return 0, nil }

type entKeyMemberRepoStub struct {
	member       *EnterpriseMember
	getByIDErr   error
	getByUserErr error
	listResult   []EnterpriseMember
	listPag      *pagination.PaginationResult
	listErr      error
}

func (s *entKeyMemberRepoStub) Create(ctx context.Context, m *EnterpriseMember) error { panic("unexpected") }
func (s *entKeyMemberRepoStub) GetByID(ctx context.Context, id int64) (*EnterpriseMember, error) {
	if s.getByIDErr != nil { return nil, s.getByIDErr }
	return s.member, nil
}
func (s *entKeyMemberRepoStub) GetByUserID(ctx context.Context, userID int64) (*EnterpriseMember, error) {
	if s.getByUserErr != nil { return nil, s.getByUserErr }
	return s.member, nil
}
func (s *entKeyMemberRepoStub) ListByEnterprise(ctx context.Context, enterpriseID int64, params pagination.PaginationParams, filters EnterpriseMemberListFilters) ([]EnterpriseMember, *pagination.PaginationResult, error) {
	return s.listResult, s.listPag, s.listErr
}
func (s *entKeyMemberRepoStub) Update(ctx context.Context, m *EnterpriseMember) error { panic("unexpected") }
func (s *entKeyMemberRepoStub) Unbind(ctx context.Context, id int64) error           { panic("unexpected") }

type entKeyGroupRepoStub struct {
	setGroupsErr error
}

func (s *entKeyGroupRepoStub) SetGroups(ctx context.Context, apiKeyID int64, groupIDs []int64) error { return s.setGroupsErr }
func (s *entKeyGroupRepoStub) GetByKeyID(ctx context.Context, apiKeyID int64) ([]APIKeyGroup, error) { return nil, nil }
func (s *entKeyGroupRepoStub) GetByGroupID(ctx context.Context, groupID int64) ([]APIKeyGroup, error) { return nil, nil }
func (s *entKeyGroupRepoStub) DeleteByKeyID(ctx context.Context, apiKeyID int64) error               { return nil }

func makeActiveKeyAdminMember() *EnterpriseMember {
	return &EnterpriseMember{
		ID: 1, EnterpriseID: 1, UserID: 500,
		Role: EnterpriseRoleAdmin, Status: StatusActive,
	}
}

func makeActiveEntForKey() *Enterprise {
	return &Enterprise{ID: 1, Status: EnterpriseStatusActive}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestCreateEnterpriseKey_Success(t *testing.T) {
	ctx := context.Background()
	apiKey := &entKeyAPIKeyRepoStub{}
	mem := &entKeyMemberRepoStub{member: makeActiveKeyAdminMember()}
	ent := &entKeyEntRepoStub{enterprise: makeActiveEntForKey()}
	group := &entKeyGroupRepoStub{}
	svc := NewEnterpriseKeyService(apiKey, mem, ent, group)

	result, err := svc.CreateEnterpriseKey(ctx, 500, CreateEnterpriseKeyRequest{
		Name: "Test Key", GroupIDs: []int64{1, 2},
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "sk-test-abc", result.Key)
}

func TestCreateEnterpriseKey_UsagePurposeTooLong(t *testing.T) {
	longStr := ""
	for i := 0; i < 250; i++ { longStr += "x" }
	svc := NewEnterpriseKeyService(&entKeyAPIKeyRepoStub{}, &entKeyMemberRepoStub{}, &entKeyEntRepoStub{}, &entKeyGroupRepoStub{})
	_, err := svc.CreateEnterpriseKey(context.Background(), 1, CreateEnterpriseKeyRequest{Name: "K", UsagePurpose: longStr})
	require.ErrorIs(t, err, ErrKeyUsagePurposeTooLong)
}

func TestCreateEnterpriseKey_InvalidTool(t *testing.T) {
	svc := NewEnterpriseKeyService(&entKeyAPIKeyRepoStub{}, &entKeyMemberRepoStub{}, &entKeyEntRepoStub{}, &entKeyGroupRepoStub{})
	_, err := svc.CreateEnterpriseKey(context.Background(), 1, CreateEnterpriseKeyRequest{Name: "K", BoundTool: "invalid_tool"})
	require.ErrorIs(t, err, ErrKeyBoundToolInvalid)
}

func TestCreateEnterpriseKey_NotAdmin(t *testing.T) {
	mem := &entKeyMemberRepoStub{
		member: &EnterpriseMember{ID: 2, EnterpriseID: 1, UserID: 500, Role: EnterpriseMemberRoleMember, Status: StatusActive},
	}
	svc := NewEnterpriseKeyService(&entKeyAPIKeyRepoStub{}, mem, &entKeyEntRepoStub{}, &entKeyGroupRepoStub{})
	_, err := svc.CreateEnterpriseKey(context.Background(), 500, CreateEnterpriseKeyRequest{Name: "K"})
	require.Error(t, err)
}

func TestUpdateEnterpriseKey_Success(t *testing.T) {
	apiKey := &entKeyAPIKeyRepoStub{key: &APIKey{ID: 1, Name: "Old", UserID: 500}}
	svc := NewEnterpriseKeyService(apiKey, &entKeyMemberRepoStub{}, &entKeyEntRepoStub{}, &entKeyGroupRepoStub{})
	newName := "Updated"
	result, err := svc.UpdateEnterpriseKey(context.Background(), 1, UpdateEnterpriseKeyRequest{Name: &newName})
	require.NoError(t, err)
	require.Equal(t, "Updated", result.Name)
}

func TestUpdateEnterpriseKey_BoundToolInvalid(t *testing.T) {
	apiKey := &entKeyAPIKeyRepoStub{key: &APIKey{ID: 1, UserID: 500}}
	svc := NewEnterpriseKeyService(apiKey, &entKeyMemberRepoStub{}, &entKeyEntRepoStub{}, &entKeyGroupRepoStub{})
	badTool := "bad_tool"
	_, err := svc.UpdateEnterpriseKey(context.Background(), 1, UpdateEnterpriseKeyRequest{BoundTool: &badTool})
	require.ErrorIs(t, err, ErrKeyBoundToolInvalid)
}

func TestListEnterpriseKeys_Empty(t *testing.T) {
	mem := &entKeyMemberRepoStub{
		listResult: []EnterpriseMember{},
		listPag:    &pagination.PaginationResult{Total: 0},
	}
	svc := NewEnterpriseKeyService(&entKeyAPIKeyRepoStub{}, mem, &entKeyEntRepoStub{}, &entKeyGroupRepoStub{})
	result, pager, err := svc.ListEnterpriseKeys(context.Background(), 1, pagination.DefaultPagination())
	require.NoError(t, err)
	require.Empty(t, result)
	require.Equal(t, int64(0), pager.Total)
}

func TestDeleteEnterpriseKey_Success_AssignedToNil(t *testing.T) {
	apiKey := &entKeyAPIKeyRepoStub{key: &APIKey{ID: 1, UserID: 500, AssignedTo: nil}}
	mem := &entKeyMemberRepoStub{member: makeActiveKeyAdminMember()}
	svc := NewEnterpriseKeyService(apiKey, mem, &entKeyEntRepoStub{}, &entKeyGroupRepoStub{})
	err := svc.DeleteEnterpriseKeyWrapper(context.Background(), 500, 1)
	require.NoError(t, err)
}

func TestDeleteEnterpriseKey_NotOwnedByAdmin(t *testing.T) {
	apiKey := &entKeyAPIKeyRepoStub{key: &APIKey{ID: 1, UserID: 999, AssignedTo: nil}}
	mem := &entKeyMemberRepoStub{member: makeActiveKeyAdminMember()}
	svc := NewEnterpriseKeyService(apiKey, mem, &entKeyEntRepoStub{}, &entKeyGroupRepoStub{})
	err := svc.DeleteEnterpriseKeyWrapper(context.Background(), 500, 1)
	require.Error(t, err)
}

func TestDeleteEnterpriseKey_WrongEnterprise(t *testing.T) {
	assignedToID := int64(2)
	apiKey := &entKeyAPIKeyRepoStub{key: &APIKey{ID: 1, UserID: 500, AssignedTo: &assignedToID}}
	mem := &entKeyMemberRepoStub{member: &EnterpriseMember{ID: 2, EnterpriseID: 99, UserID: 600}}
	svc := NewEnterpriseKeyService(apiKey, mem, &entKeyEntRepoStub{}, &entKeyGroupRepoStub{})
	err := svc.DeleteEnterpriseKeyWrapper(context.Background(), 500, 1)
	require.Error(t, err)
}
