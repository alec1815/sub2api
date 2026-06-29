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
// Stubs
// ---------------------------------------------------------------------------

type entRepoStub struct {
	enterprise   *Enterprise
	getByIDErr   error
	createErr    error
	updateErr    error
	softDelErr   error
	balance      float64
	getBalErr    error
	listResult   []Enterprise
	listPag      *pagination.PaginationResult
	listErr      error
}

func (s *entRepoStub) Create(ctx context.Context, e *Enterprise) error {
	if s.createErr != nil {
		return s.createErr
	}
	if e != nil {
		e.ID = 1
	}
	return nil
}
func (s *entRepoStub) GetByID(ctx context.Context, id int64) (*Enterprise, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	if s.enterprise != nil {
		return s.enterprise, nil
	}
	panic("unexpected GetByID call")
}
func (s *entRepoStub) List(ctx context.Context, params pagination.PaginationParams, filters EnterpriseListFilters) ([]Enterprise, *pagination.PaginationResult, error) {
	if s.listErr != nil {
		return nil, nil, s.listErr
	}
	return s.listResult, s.listPag, nil
}
func (s *entRepoStub) Update(ctx context.Context, e *Enterprise) error { return s.updateErr }
func (s *entRepoStub) SoftDelete(ctx context.Context, id int64) error { return s.softDelErr }
func (s *entRepoStub) GetBalance(ctx context.Context, id int64) (float64, error) {
	return s.balance, s.getBalErr
}

type memberRepoStub struct {
	member       *EnterpriseMember
	getByIDErr   error
	getByUserErr error
	createErr    error
	updateErr    error
	unbindErr    error
	listResult   []EnterpriseMember
	listPag      *pagination.PaginationResult
	listErr      error
}

func (s *memberRepoStub) Create(ctx context.Context, m *EnterpriseMember) error {
	if s.createErr != nil {
		return s.createErr
	}
	if m != nil {
		m.ID = 100
	}
	return nil
}
func (s *memberRepoStub) GetByID(ctx context.Context, id int64) (*EnterpriseMember, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	if s.member != nil {
		return s.member, nil
	}
	panic("unexpected GetByID call")
}
func (s *memberRepoStub) GetByUserID(ctx context.Context, userID int64) (*EnterpriseMember, error) {
	if s.getByUserErr != nil {
		return nil, s.getByUserErr
	}
	if s.member != nil {
		return s.member, nil
	}
	return nil, ErrEnterpriseMemberNotFound
}
func (s *memberRepoStub) ListByEnterprise(ctx context.Context, enterpriseID int64, params pagination.PaginationParams, filters EnterpriseMemberListFilters) ([]EnterpriseMember, *pagination.PaginationResult, error) {
	if s.listErr != nil {
		return nil, nil, s.listErr
	}
	return s.listResult, s.listPag, nil
}
func (s *memberRepoStub) Update(ctx context.Context, m *EnterpriseMember) error { return s.updateErr }
func (s *memberRepoStub) Unbind(ctx context.Context, id int64) error           { return s.unbindErr }

type subRepoStub struct {
	subs          []EnterpriseSubscription
	listActiveErr error
	getByEntErr   error
	updateStatErr error
	createErr     error
	getByIDErr    error
	getByIDSub    *EnterpriseSubscription
}

func (s *subRepoStub) Create(ctx context.Context, sub *EnterpriseSubscription) error {
	if s.createErr != nil {
		return s.createErr
	}
	sub.ID = 200
	return nil
}
func (s *subRepoStub) GetByID(ctx context.Context, id int64) (*EnterpriseSubscription, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	if s.getByIDSub != nil {
		return s.getByIDSub, nil
	}
	panic("unexpected GetByID")
}
func (s *subRepoStub) GetByEnterprise(ctx context.Context, enterpriseID int64) ([]EnterpriseSubscription, error) {
	if s.getByEntErr != nil {
		return nil, s.getByEntErr
	}
	return s.subs, nil
}
func (s *subRepoStub) ListActive(ctx context.Context, enterpriseID int64) ([]EnterpriseSubscription, error) {
	if s.listActiveErr != nil {
		return nil, s.listActiveErr
	}
	return s.subs, nil
}
func (s *subRepoStub) UpdateStatus(ctx context.Context, id int64, status string) error { return s.updateStatErr }

type userRepoStubForEnt struct {
	user       *User
	getByEmailErr error
	createErr  error
	getByIDErr error
	updateErr  error
}

func (s *userRepoStubForEnt) Create(ctx context.Context, u *User) error {
	if s.createErr != nil {
		return s.createErr
	}
	u.ID = 500
	return nil
}
func (s *userRepoStubForEnt) GetByID(ctx context.Context, id int64) (*User, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	if s.user != nil {
		return s.user, nil
	}
	panic("unexpected GetByID")
}
func (s *userRepoStubForEnt) GetByEmail(ctx context.Context, email string) (*User, error) {
	if s.getByEmailErr != nil {
		return nil, s.getByEmailErr
	}
	return s.user, nil
}

// Remaining UserRepository methods — should not be called
func (s *userRepoStubForEnt) GetFirstAdmin(context.Context) (*User, error)                               { panic("unexpected") }
func (s *userRepoStubForEnt) Update(context.Context, *User) error                                         { return s.updateErr }
func (s *userRepoStubForEnt) Delete(context.Context, int64) error                                         { panic("unexpected") }
func (s *userRepoStubForEnt) GetUserAvatar(context.Context, int64) (*UserAvatar, error)                   { panic("unexpected") }
func (s *userRepoStubForEnt) UpsertUserAvatar(context.Context, int64, UpsertUserAvatarInput) (*UserAvatar, error) { panic("unexpected") }
func (s *userRepoStubForEnt) DeleteUserAvatar(context.Context, int64) error                              { panic("unexpected") }
func (s *userRepoStubForEnt) List(context.Context, pagination.PaginationParams) ([]User, *pagination.PaginationResult, error) { panic("unexpected") }
func (s *userRepoStubForEnt) ListWithFilters(context.Context, pagination.PaginationParams, UserListFilters) ([]User, *pagination.PaginationResult, error) { panic("unexpected") }
func (s *userRepoStubForEnt) GetLatestUsedAtByUserIDs(context.Context, []int64) (map[int64]*time.Time, error) { panic("unexpected") }
func (s *userRepoStubForEnt) GetLatestUsedAtByUserID(context.Context, int64) (*time.Time, error)          { panic("unexpected") }
func (s *userRepoStubForEnt) UpdateUserLastActiveAt(context.Context, int64, time.Time) error              { panic("unexpected") }
func (s *userRepoStubForEnt) UpdateBalance(context.Context, int64, float64) error                         { panic("unexpected") }
func (s *userRepoStubForEnt) DeductBalance(context.Context, int64, float64) error                         { panic("unexpected") }
func (s *userRepoStubForEnt) UpdateConcurrency(context.Context, int64, int) error                         { panic("unexpected") }
func (s *userRepoStubForEnt) BatchSetConcurrency(context.Context, []int64, int) (int, error)              { return 0, nil }
func (s *userRepoStubForEnt) BatchAddConcurrency(context.Context, []int64, int) (int, error)              { return 0, nil }
func (s *userRepoStubForEnt) ExistsByEmail(context.Context, string) (bool, error)                         { return false, nil }
func (s *userRepoStubForEnt) RemoveGroupFromAllowedGroups(context.Context, int64) (int64, error)          { return 0, nil }
func (s *userRepoStubForEnt) AddGroupToAllowedGroups(context.Context, int64, int64) error                 { return nil }
func (s *userRepoStubForEnt) RemoveGroupFromUserAllowedGroups(context.Context, int64, int64) error        { return nil }
func (s *userRepoStubForEnt) ListUserAuthIdentities(context.Context, int64) ([]UserAuthIdentityRecord, error) { return nil, nil }
func (s *userRepoStubForEnt) GetByIDIncludeDeleted(context.Context, int64) (*User, error)                 { panic("unexpected") }
func (s *userRepoStubForEnt) UpdateTotpSecret(context.Context, int64, *string) error                       { panic("unexpected") }
func (s *userRepoStubForEnt) EnableTotp(context.Context, int64) error                                      { panic("unexpected") }
func (s *userRepoStubForEnt) DisableTotp(context.Context, int64) error                                     { panic("unexpected") }
func (s *userRepoStubForEnt) UnbindUserAuthProvider(context.Context, int64, string) error                   { panic("unexpected") }

type deptRepoStubForEnt struct {
	dept          *Department
	getByIDErr    error
	createErr     error
	updateErr     error
	deleteErr     error
	listResult    []Department
	listPag       *pagination.PaginationResult
	listErr       error
	hasChildren   bool
	hasChildrenErr error
	hasMembers    bool
	hasMembersErr error
	treeResult    []Department
	treeErr       error
}

func (s *deptRepoStubForEnt) Create(ctx context.Context, d *Department) error { return s.createErr }
func (s *deptRepoStubForEnt) GetByID(ctx context.Context, id int64) (*Department, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	if s.dept != nil {
		return s.dept, nil
	}
	panic("unexpected GetByID")
}
func (s *deptRepoStubForEnt) TreeByEnterprise(ctx context.Context, enterpriseID int64) ([]Department, error) {
	return s.treeResult, s.treeErr
}
func (s *deptRepoStubForEnt) ListByEnterprise(ctx context.Context, enterpriseID int64, params pagination.PaginationParams, filters DepartmentListFilters) ([]Department, *pagination.PaginationResult, error) {
	return s.listResult, s.listPag, s.listErr
}
func (s *deptRepoStubForEnt) Update(ctx context.Context, d *Department) error     { return s.updateErr }
func (s *deptRepoStubForEnt) Delete(ctx context.Context, id int64) error           { return s.deleteErr }
func (s *deptRepoStubForEnt) HasChildren(ctx context.Context, id int64) (bool, error) { return s.hasChildren, s.hasChildrenErr }
func (s *deptRepoStubForEnt) HasMembers(ctx context.Context, id int64) (bool, error)  { return s.hasMembers, s.hasMembersErr }

func makeEntService(ent *entRepoStub, mem *memberRepoStub, sub *subRepoStub, usr *userRepoStubForEnt, dept *deptRepoStubForEnt) *EnterpriseService {
	return NewEnterpriseService(ent, mem, sub, usr, dept)
}

func makeActiveEnterprise() *Enterprise {
	return &Enterprise{
		ID: 1, Name: "Test Corp", Status: EnterpriseStatusActive,
		Balance: 1000, AdminUserID: 500,
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestCreateEnterprise_Success(t *testing.T) {
	ctx := context.Background()
	ent := &entRepoStub{enterprise: &Enterprise{}}
	mem := &memberRepoStub{} // GetByUserID returns ErrEnterpriseMemberNotFound
	usr := &userRepoStubForEnt{}
	sub := &subRepoStub{}
	dept := &deptRepoStubForEnt{}
	svc := makeEntService(ent, mem, sub, usr, dept)

	result, err := svc.CreateEnterprise(ctx, CreateEnterpriseRequest{
		Name: "New Corp", AdminEmail: "admin@test.com", AdminName: "Admin",
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "New Corp", result.Name)
	require.Equal(t, EnterpriseStatusActive, result.Status)
}

func TestCreateEnterprise_NameRequired(t *testing.T) {
	svc := makeEntService(&entRepoStub{}, &memberRepoStub{}, &subRepoStub{}, &userRepoStubForEnt{}, &deptRepoStubForEnt{})
	_, err := svc.CreateEnterprise(context.Background(), CreateEnterpriseRequest{AdminEmail: "a@b.com"})
	require.ErrorIs(t, err, ErrEnterpriseNameRequired)
}

func TestCreateEnterprise_AdminRequired(t *testing.T) {
	svc := makeEntService(&entRepoStub{}, &memberRepoStub{}, &subRepoStub{}, &userRepoStubForEnt{}, &deptRepoStubForEnt{})
	_, err := svc.CreateEnterprise(context.Background(), CreateEnterpriseRequest{Name: "X"})
	require.ErrorIs(t, err, ErrEnterpriseAdminRequired)
}

func TestCreateEnterprise_AdminAlreadyMember(t *testing.T) {
	ent := &entRepoStub{enterprise: &Enterprise{}}
	usr := &userRepoStubForEnt{user: &User{ID: 500, Email: "admin@test.com"}}
	// member already active in another enterprise
	mem := &memberRepoStub{
		member: &EnterpriseMember{ID: 1, UserID: 500, EnterpriseID: 99, Role: EnterpriseRoleAdmin, Status: StatusActive},
	}
	svc := makeEntService(ent, mem, &subRepoStub{}, usr, &deptRepoStubForEnt{})

	_, err := svc.CreateEnterprise(context.Background(), CreateEnterpriseRequest{
		Name: "X", AdminEmail: "admin@test.com",
	})
	require.ErrorIs(t, err, ErrEnterpriseMemberAlreadyActive)
}

func TestGetEnterprise_Success(t *testing.T) {
	ent := &entRepoStub{enterprise: makeActiveEnterprise()}
	svc := makeEntService(ent, &memberRepoStub{}, &subRepoStub{}, &userRepoStubForEnt{}, &deptRepoStubForEnt{})

	result, err := svc.GetEnterprise(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, int64(1), result.ID)
	require.Equal(t, "Test Corp", result.Name)
}

func TestGetEnterprise_NotFound(t *testing.T) {
	ent := &entRepoStub{getByIDErr: ErrEnterpriseNotFound}
	svc := makeEntService(ent, &memberRepoStub{}, &subRepoStub{}, &userRepoStubForEnt{}, &deptRepoStubForEnt{})

	_, err := svc.GetEnterprise(context.Background(), 999)
	require.Error(t, err)
}

func TestListEnterprises_Success(t *testing.T) {
	ent := &entRepoStub{
		listResult: []Enterprise{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}},
		listPag:    &pagination.PaginationResult{Total: 2, Page: 1, PageSize: 20, Pages: 1},
	}
	svc := makeEntService(ent, &memberRepoStub{}, &subRepoStub{}, &userRepoStubForEnt{}, &deptRepoStubForEnt{})

	result, pager, err := svc.ListEnterprises(context.Background(), pagination.DefaultPagination(), EnterpriseListFilters{})
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, int64(2), pager.Total)
}

func TestToggleStatus_ActiveToDisabled(t *testing.T) {
	ent := &entRepoStub{enterprise: makeActiveEnterprise()}
	svc := makeEntService(ent, &memberRepoStub{}, &subRepoStub{}, &userRepoStubForEnt{}, &deptRepoStubForEnt{})

	result, err := svc.ToggleStatus(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, EnterpriseStatusDisabled, result.Status)
}

func TestToggleStatus_DisabledToActive(t *testing.T) {
	enterprise := makeActiveEnterprise()
	enterprise.Status = EnterpriseStatusDisabled
	ent := &entRepoStub{enterprise: enterprise}
	svc := makeEntService(ent, &memberRepoStub{}, &subRepoStub{}, &userRepoStubForEnt{}, &deptRepoStubForEnt{})

	result, err := svc.ToggleStatus(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, EnterpriseStatusActive, result.Status)
}

func TestDeleteEnterprise_NotDisabled(t *testing.T) {
	ent := &entRepoStub{enterprise: makeActiveEnterprise()}
	svc := makeEntService(ent, &memberRepoStub{}, &subRepoStub{}, &userRepoStubForEnt{}, &deptRepoStubForEnt{})

	err := svc.DeleteEnterprise(context.Background(), 1)
	require.Error(t, err)
}

func TestDeleteEnterprise_Success(t *testing.T) {
	enterprise := makeActiveEnterprise()
	enterprise.Status = EnterpriseStatusDisabled
	ent := &entRepoStub{enterprise: enterprise}
	mem := &memberRepoStub{
		listResult: []EnterpriseMember{{ID: 1, Role: EnterpriseRoleMember, Status: StatusActive}},
		listPag:    &pagination.PaginationResult{Total: 1},
	}
	sub := &subRepoStub{
		subs: []EnterpriseSubscription{{ID: 1, Status: EnterpriseSubStatusActive}},
	}
	svc := makeEntService(ent, mem, sub, &userRepoStubForEnt{}, &deptRepoStubForEnt{})

	err := svc.DeleteEnterprise(context.Background(), 1)
	require.NoError(t, err)
}

func TestGetEnterpriseBalance_Success(t *testing.T) {
	ent := &entRepoStub{balance: 1234.56}
	svc := makeEntService(ent, &memberRepoStub{}, &subRepoStub{}, &userRepoStubForEnt{}, &deptRepoStubForEnt{})

	balance, err := svc.GetEnterpriseBalance(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, 1234.56, balance)
}
