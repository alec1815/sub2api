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
// Stubs (lightweight, focusing on member service)
// ---------------------------------------------------------------------------

type entRepoStubForMember struct {
	enterprise *Enterprise
	getByIDErr error
}

func (s *entRepoStubForMember) Create(ctx context.Context, e *Enterprise) error                                     { panic("unexpected") }
func (s *entRepoStubForMember) GetByID(ctx context.Context, id int64) (*Enterprise, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	if s.enterprise != nil {
		return s.enterprise, nil
	}
	panic("unexpected GetByID")
}
func (s *entRepoStubForMember) List(ctx context.Context, params pagination.PaginationParams, filters EnterpriseListFilters) ([]Enterprise, *pagination.PaginationResult, error) { panic("unexpected") }
func (s *entRepoStubForMember) Update(ctx context.Context, e *Enterprise) error                                   { panic("unexpected") }
func (s *entRepoStubForMember) SoftDelete(ctx context.Context, id int64) error                                    { panic("unexpected") }
func (s *entRepoStubForMember) GetBalance(ctx context.Context, id int64) (float64, error)                         { return 0, nil }

type memberRepoStubForService struct {
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

func (s *memberRepoStubForService) Create(ctx context.Context, m *EnterpriseMember) error {
	if s.createErr != nil {
		return s.createErr
	}
	m.ID = 100
	return nil
}
func (s *memberRepoStubForService) GetByID(ctx context.Context, id int64) (*EnterpriseMember, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	return s.member, nil
}
func (s *memberRepoStubForService) GetByUserID(ctx context.Context, userID int64) (*EnterpriseMember, error) {
	if s.getByUserErr != nil {
		return nil, s.getByUserErr
	}
	if s.member != nil {
		return s.member, nil
	}
	return nil, ErrEnterpriseMemberNotFound
}
func (s *memberRepoStubForService) ListByEnterprise(ctx context.Context, enterpriseID int64, params pagination.PaginationParams, filters EnterpriseMemberListFilters) ([]EnterpriseMember, *pagination.PaginationResult, error) {
	if s.listErr != nil {
		return nil, nil, s.listErr
	}
	return s.listResult, s.listPag, nil
}
func (s *memberRepoStubForService) Update(ctx context.Context, m *EnterpriseMember) error { return s.updateErr }
func (s *memberRepoStubForService) Unbind(ctx context.Context, id int64) error             { return s.unbindErr }

type userRepoStubForMember struct {
	user       *User
	getByEmailErr error
	createErr  error
	getByIDErr error
	updateErr  error
}

func (s *userRepoStubForMember) Create(ctx context.Context, u *User) error {
	if s.createErr != nil {
		return s.createErr
	}
	u.ID = 600
	return nil
}
func (s *userRepoStubForMember) GetByID(ctx context.Context, id int64) (*User, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	return s.user, nil
}
func (s *userRepoStubForMember) GetByEmail(ctx context.Context, email string) (*User, error) {
	if s.getByEmailErr != nil {
		return nil, s.getByEmailErr
	}
	return s.user, nil
}
func (s *userRepoStubForMember) Update(ctx context.Context, u *User) error { return s.updateErr }

func (s *userRepoStubForMember) GetFirstAdmin(context.Context) (*User, error) { panic("unexpected") }
func (s *userRepoStubForMember) Delete(context.Context, int64) error          { panic("unexpected") }
func (s *userRepoStubForMember) GetUserAvatar(context.Context, int64) (*UserAvatar, error) { panic("unexpected") }
func (s *userRepoStubForMember) UpsertUserAvatar(context.Context, int64, UpsertUserAvatarInput) (*UserAvatar, error) { panic("unexpected") }
func (s *userRepoStubForMember) DeleteUserAvatar(context.Context, int64) error { panic("unexpected") }
func (s *userRepoStubForMember) List(context.Context, pagination.PaginationParams) ([]User, *pagination.PaginationResult, error) { panic("unexpected") }
func (s *userRepoStubForMember) ListWithFilters(context.Context, pagination.PaginationParams, UserListFilters) ([]User, *pagination.PaginationResult, error) { panic("unexpected") }
func (s *userRepoStubForMember) GetLatestUsedAtByUserIDs(context.Context, []int64) (map[int64]*time.Time, error) { panic("unexpected") }
func (s *userRepoStubForMember) GetLatestUsedAtByUserID(context.Context, int64) (*time.Time, error)          { panic("unexpected") }
func (s *userRepoStubForMember) UpdateUserLastActiveAt(context.Context, int64, time.Time) error              { panic("unexpected") }
func (s *userRepoStubForMember) UpdateBalance(context.Context, int64, float64) error                   { panic("unexpected") }
func (s *userRepoStubForMember) DeductBalance(context.Context, int64, float64) error                   { panic("unexpected") }
func (s *userRepoStubForMember) UpdateConcurrency(context.Context, int64, int) error                   { panic("unexpected") }
func (s *userRepoStubForMember) BatchSetConcurrency(context.Context, []int64, int) (int, error)         { return 0, nil }
func (s *userRepoStubForMember) BatchAddConcurrency(context.Context, []int64, int) (int, error)         { return 0, nil }
func (s *userRepoStubForMember) ExistsByEmail(context.Context, string) (bool, error)                    { return false, nil }
func (s *userRepoStubForMember) RemoveGroupFromAllowedGroups(context.Context, int64) (int64, error)     { return 0, nil }
func (s *userRepoStubForMember) AddGroupToAllowedGroups(context.Context, int64, int64) error            { return nil }
func (s *userRepoStubForMember) RemoveGroupFromUserAllowedGroups(context.Context, int64, int64) error   { return nil }
func (s *userRepoStubForMember) ListUserAuthIdentities(context.Context, int64) ([]UserAuthIdentityRecord, error) { return nil, nil }
func (s *userRepoStubForMember) GetByIDIncludeDeleted(context.Context, int64) (*User, error)            { panic("unexpected") }
func (s *userRepoStubForMember) UpdateTotpSecret(context.Context, int64, *string) error                  { panic("unexpected") }
func (s *userRepoStubForMember) EnableTotp(context.Context, int64) error                                 { panic("unexpected") }
func (s *userRepoStubForMember) DisableTotp(context.Context, int64) error                                { panic("unexpected") }
func (s *userRepoStubForMember) UnbindUserAuthProvider(context.Context, int64, string) error              { panic("unexpected") }

type deptRepoStubForMember struct {
	dept        *Department
	getByIDErr  error
}

func (s *deptRepoStubForMember) Create(ctx context.Context, d *Department) error { panic("unexpected") }
func (s *deptRepoStubForMember) GetByID(ctx context.Context, id int64) (*Department, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	return s.dept, nil
}
func (s *deptRepoStubForMember) TreeByEnterprise(ctx context.Context, enterpriseID int64) ([]Department, error) { panic("unexpected") }
func (s *deptRepoStubForMember) ListByEnterprise(ctx context.Context, enterpriseID int64, params pagination.PaginationParams, filters DepartmentListFilters) ([]Department, *pagination.PaginationResult, error) { panic("unexpected") }
func (s *deptRepoStubForMember) Update(ctx context.Context, d *Department) error     { panic("unexpected") }
func (s *deptRepoStubForMember) Delete(ctx context.Context, id int64) error           { panic("unexpected") }
func (s *deptRepoStubForMember) HasChildren(ctx context.Context, id int64) (bool, error) { return false, nil }
func (s *deptRepoStubForMember) HasMembers(ctx context.Context, id int64) (bool, error)  { return false, nil }

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestCreateMember_Success(t *testing.T) {
	ctx := context.Background()
	ent := &entRepoStubForMember{enterprise: &Enterprise{ID: 1, Status: EnterpriseStatusActive}}
	mem := &memberRepoStubForService{}
	usr := &userRepoStubForMember{} // no existing user
	dept := &deptRepoStubForMember{}
	svc := NewEnterpriseMemberService(mem, usr, ent, dept)

	result, err := svc.CreateMember(ctx, 1, CreateMemberRequest{Email: "new@test.com", Username: "New User"})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, StatusActive, result.Status)
	require.Equal(t, EnterpriseMemberRoleMember, result.Role)
}

func TestCreateMember_EnterpriseNotActive(t *testing.T) {
	ent := &entRepoStubForMember{enterprise: &Enterprise{ID: 1, Status: EnterpriseStatusDisabled}}
	svc := NewEnterpriseMemberService(&memberRepoStubForService{}, &userRepoStubForMember{}, ent, &deptRepoStubForMember{})

	_, err := svc.CreateMember(context.Background(), 1, CreateMemberRequest{Email: "x@y.com"})
	require.ErrorIs(t, err, ErrMemberEnterpriseNotActive)
}

func TestCreateMember_UserAlreadyMember(t *testing.T) {
	ent := &entRepoStubForMember{enterprise: &Enterprise{ID: 1, Status: EnterpriseStatusActive}}
	usr := &userRepoStubForMember{user: &User{ID: 1, Email: "existing@test.com"}}
	mem := &memberRepoStubForService{
		member: &EnterpriseMember{ID: 1, UserID: 1, EnterpriseID: 2, Status: StatusActive},
	}
	svc := NewEnterpriseMemberService(mem, usr, ent, &deptRepoStubForMember{})

	_, err := svc.CreateMember(context.Background(), 1, CreateMemberRequest{Email: "existing@test.com"})
	require.ErrorIs(t, err, ErrMemberUserAlreadyMember)
}

func TestCreateMember_DepartmentNotInEnterprise(t *testing.T) {
	ent := &entRepoStubForMember{enterprise: &Enterprise{ID: 1, Status: EnterpriseStatusActive}}
	mem := &memberRepoStubForService{}
	usr := &userRepoStubForMember{}
	wrongEntDeptID := int64(9)
	dept := &deptRepoStubForMember{
		dept: &Department{ID: 9, EnterpriseID: 99, Name: "Other Dept"}, // belongs to enterprise 99
	}
	svc := NewEnterpriseMemberService(mem, usr, ent, dept)

	_, err := svc.CreateMember(context.Background(), 1, CreateMemberRequest{
		Email: "x@y.com", DepartmentID: &wrongEntDeptID,
	})
	require.ErrorIs(t, err, ErrMemberDepartmentInvalid)
}

func TestGetMember_Success(t *testing.T) {
	expected := &EnterpriseMember{ID: 1, UserID: 600, Role: EnterpriseMemberRoleMember, Status: StatusActive}
	mem := &memberRepoStubForService{member: expected}
	svc := NewEnterpriseMemberService(mem, &userRepoStubForMember{}, &entRepoStubForMember{}, &deptRepoStubForMember{})

	result, err := svc.GetMember(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, int64(1), result.ID)
}

func TestListMembers_Success(t *testing.T) {
	mem := &memberRepoStubForService{
		listResult: []EnterpriseMember{{ID: 1}, {ID: 2}},
		listPag:    &pagination.PaginationResult{Total: 2},
	}
	svc := NewEnterpriseMemberService(mem, &userRepoStubForMember{}, &entRepoStubForMember{}, &deptRepoStubForMember{})

	result, pager, err := svc.ListMembers(context.Background(), 1, pagination.DefaultPagination(), EnterpriseMemberListFilters{})
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, int64(2), pager.Total)
}

func TestUnbindMember_AdminCannotUnbindSelf(t *testing.T) {
	mem := &memberRepoStubForService{
		member: &EnterpriseMember{ID: 1, UserID: 1, Role: EnterpriseRoleAdmin, Status: StatusActive},
	}
	svc := NewEnterpriseMemberService(mem, &userRepoStubForMember{}, &entRepoStubForMember{}, &deptRepoStubForMember{})

	err := svc.UnbindMember(context.Background(), 1)
	require.ErrorIs(t, err, ErrMemberCannotUnbindSelf)
}

func TestUnbindMember_Success(t *testing.T) {
	mem := &memberRepoStubForService{
		member: &EnterpriseMember{ID: 1, UserID: 1, Role: EnterpriseMemberRoleMember, Status: StatusActive},
	}
	svc := NewEnterpriseMemberService(mem, &userRepoStubForMember{}, &entRepoStubForMember{}, &deptRepoStubForMember{})

	err := svc.UnbindMember(context.Background(), 1)
	require.NoError(t, err)
}

func TestGetMemberByUserID_Success(t *testing.T) {
	expected := &EnterpriseMember{ID: 1, UserID: 600, Role: EnterpriseMemberRoleMember}
	mem := &memberRepoStubForService{member: expected}
	svc := NewEnterpriseMemberService(mem, &userRepoStubForMember{}, &entRepoStubForMember{}, &deptRepoStubForMember{})

	result, err := svc.GetMemberByUserID(context.Background(), 600)
	require.NoError(t, err)
	require.Equal(t, int64(600), result.UserID)
}
