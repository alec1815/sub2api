//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Stubs
// ---------------------------------------------------------------------------

type deptRepoStub struct {
	dept           *Department
	getByIDErr     error
	createErr      error
	updateErr      error
	deleteErr      error
	listResult     []Department
	listPag        *pagination.PaginationResult
	listErr        error
	hasChildren    bool
	hasChildrenErr error
	hasMembers     bool
	hasMembersErr  error
	treeResult     []Department
	treeErr        error
}

func (s *deptRepoStub) Create(ctx context.Context, d *Department) error {
	if s.createErr != nil {
		return s.createErr
	}
	d.ID = 10
	return nil
}
func (s *deptRepoStub) GetByID(ctx context.Context, id int64) (*Department, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	if s.dept != nil {
		return s.dept, nil
	}
	panic("unexpected GetByID")
}
func (s *deptRepoStub) TreeByEnterprise(ctx context.Context, enterpriseID int64) ([]Department, error) {
	return s.treeResult, s.treeErr
}
func (s *deptRepoStub) ListByEnterprise(ctx context.Context, enterpriseID int64, params pagination.PaginationParams, filters DepartmentListFilters) ([]Department, *pagination.PaginationResult, error) {
	return s.listResult, s.listPag, s.listErr
}
func (s *deptRepoStub) Update(ctx context.Context, d *Department) error          { return s.updateErr }
func (s *deptRepoStub) Delete(ctx context.Context, id int64) error                { return s.deleteErr }
func (s *deptRepoStub) HasChildren(ctx context.Context, id int64) (bool, error)   { return s.hasChildren, s.hasChildrenErr }
func (s *deptRepoStub) HasMembers(ctx context.Context, id int64) (bool, error)    { return s.hasMembers, s.hasMembersErr }

type entRepoStubForDept struct {
	enterprise *Enterprise
	getByIDErr error
}

func (s *entRepoStubForDept) Create(ctx context.Context, e *Enterprise) error          { panic("unexpected") }
func (s *entRepoStubForDept) GetByID(ctx context.Context, id int64) (*Enterprise, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}
	return s.enterprise, nil
}
func (s *entRepoStubForDept) List(ctx context.Context, params pagination.PaginationParams, filters EnterpriseListFilters) ([]Enterprise, *pagination.PaginationResult, error) { panic("unexpected") }
func (s *entRepoStubForDept) Update(ctx context.Context, e *Enterprise) error        { panic("unexpected") }
func (s *entRepoStubForDept) SoftDelete(ctx context.Context, id int64) error         { panic("unexpected") }
func (s *entRepoStubForDept) GetBalance(ctx context.Context, id int64) (float64, error) { return 0, nil }

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestDepartment_GetTree(t *testing.T) {
	dept := &deptRepoStub{
		treeResult: []Department{
			{ID: 1, Name: "Root", ParentID: 0},
			{ID: 2, Name: "Child", ParentID: 1},
		},
	}
	svc := NewDepartmentService(dept, &entRepoStubForDept{})

	result, err := svc.GetTree(context.Background(), 1)
	require.NoError(t, err)
	require.Len(t, result, 2)
}

func TestDepartment_ListDepartments(t *testing.T) {
	dept := &deptRepoStub{
		listResult: []Department{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}},
		listPag:    &pagination.PaginationResult{Total: 2, Page: 1, PageSize: 20, Pages: 1},
	}
	svc := NewDepartmentService(dept, &entRepoStubForDept{})

	result, pager, err := svc.ListDepartments(context.Background(), 1, pagination.DefaultPagination(), DepartmentListFilters{})
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, int64(2), pager.Total)
}

func TestDepartment_Create(t *testing.T) {
	ent := &entRepoStubForDept{enterprise: &Enterprise{ID: 1, Status: EnterpriseStatusActive}}
	dept := &deptRepoStub{}
	svc := NewDepartmentService(dept, ent)

	result, err := svc.CreateDepartment(context.Background(), CreateDepartmentRequest{
		EnterpriseID: 1, Name: "Engineering",
	})
	require.NoError(t, err)
	require.Equal(t, "Engineering", result.Name)
}

func TestDepartment_Create_NameRequired(t *testing.T) {
	svc := NewDepartmentService(&deptRepoStub{}, &entRepoStubForDept{})
	_, err := svc.CreateDepartment(context.Background(), CreateDepartmentRequest{EnterpriseID: 1})
	require.ErrorIs(t, err, ErrDepartmentNameRequired)
}

func TestDepartment_Update(t *testing.T) {
	dept := &deptRepoStub{dept: &Department{ID: 1, Name: "Old", EnterpriseID: 1}}
	svc := NewDepartmentService(dept, &entRepoStubForDept{})

	newName := "New Name"
	result, err := svc.UpdateDepartment(context.Background(), 1, UpdateDepartmentRequest{Name: &newName})
	require.NoError(t, err)
	require.Equal(t, "New Name", result.Name)
}

func TestDepartment_Update_CircularRef(t *testing.T) {
	dept := &deptRepoStub{dept: &Department{ID: 1, Name: "X", EnterpriseID: 1}}
	svc := NewDepartmentService(dept, &entRepoStubForDept{})

	selfParentID := int64(1)
	_, err := svc.UpdateDepartment(context.Background(), 1, UpdateDepartmentRequest{ParentID: &selfParentID})
	require.ErrorIs(t, err, ErrDepartmentCircularRef)
}

func TestDepartment_Delete_HasChildren(t *testing.T) {
	dept := &deptRepoStub{hasChildren: true}
	svc := NewDepartmentService(dept, &entRepoStubForDept{})

	err := svc.DeleteDepartment(context.Background(), 1)
	require.ErrorIs(t, err, ErrDepartmentHasChildrenErr)
}

func TestDepartment_Delete_HasMembers(t *testing.T) {
	dept := &deptRepoStub{hasChildren: false, hasMembers: true}
	svc := NewDepartmentService(dept, &entRepoStubForDept{})

	err := svc.DeleteDepartment(context.Background(), 1)
	require.ErrorIs(t, err, ErrDepartmentHasMembersErr)
}

func TestDepartment_Delete_Success(t *testing.T) {
	dept := &deptRepoStub{hasChildren: false, hasMembers: false}
	svc := NewDepartmentService(dept, &entRepoStubForDept{})

	err := svc.DeleteDepartment(context.Background(), 1)
	require.NoError(t, err)
}

func TestDepartment_Get(t *testing.T) {
	dept := &deptRepoStub{dept: &Department{ID: 1, Name: "Eng"}}
	svc := NewDepartmentService(dept, &entRepoStubForDept{})

	result, err := svc.GetDepartment(context.Background(), 1)
	require.NoError(t, err)
	require.Equal(t, "Eng", result.Name)
}
