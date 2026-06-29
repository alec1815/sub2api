package service

import (
	"context"
	"fmt"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

var (
	ErrDepartmentCircularRef = infraerrors.BadRequest("DEPARTMENT_CIRCULAR_REF", "parent department cannot be a descendant of the department itself")
	ErrDepartmentHasChildrenErr = ErrDepartmentHasChildren
	ErrDepartmentHasMembersErr  = ErrDepartmentHasMembers
)

// DepartmentListFilters defines list query filters for departments.
type DepartmentListFilters struct {
	Status string
}

// DepartmentRepository interface
type DepartmentRepository interface {
	Create(ctx context.Context, d *Department) error
	GetByID(ctx context.Context, id int64) (*Department, error)
	TreeByEnterprise(ctx context.Context, enterpriseID int64) ([]Department, error)
	ListByEnterprise(ctx context.Context, enterpriseID int64, params pagination.PaginationParams, filters DepartmentListFilters) ([]Department, *pagination.PaginationResult, error)
	Update(ctx context.Context, d *Department) error
	Delete(ctx context.Context, id int64) error
	HasChildren(ctx context.Context, id int64) (bool, error)
	HasMembers(ctx context.Context, id int64) (bool, error)
}

// ============================================================================
// DepartmentService — 部门管理
// ============================================================================

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	EnterpriseID int64  `json:"enterprise_id"`
	ParentID     int64  `json:"parent_id"`
	Name         string `json:"name"`
	OrderNum     int    `json:"order_num"`
	Leader       string `json:"leader"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	ParentID *int64  `json:"parent_id"`
	Name     *string `json:"name"`
	OrderNum *int    `json:"order_num"`
	Leader   *string `json:"leader"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
}

// DepartmentService 部门业务逻辑层
type DepartmentService struct {
	deptRepo DepartmentRepository
	entRepo  EnterpriseRepository
}

// NewDepartmentService 构造函数
func NewDepartmentService(deptRepo DepartmentRepository, entRepo EnterpriseRepository) *DepartmentService {
	return &DepartmentService{deptRepo: deptRepo, entRepo: entRepo}
}

// GetTree 获取部门树形列表
func (s *DepartmentService) GetTree(ctx context.Context, enterpriseID int64) ([]Department, error) {
	all, err := s.deptRepo.TreeByEnterprise(ctx, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("list departments: %w", err)
	}
	return all, nil
}

// ListDepartments 部门列表（分页）
func (s *DepartmentService) ListDepartments(
	ctx context.Context,
	enterpriseID int64,
	params pagination.PaginationParams,
	filters DepartmentListFilters,
) ([]Department, *pagination.PaginationResult, error) {
	depts, result, err := s.deptRepo.ListByEnterprise(ctx, enterpriseID, params, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("list departments: %w", err)
	}
	return depts, result, nil
}

// CreateDepartment 创建部门
func (s *DepartmentService) CreateDepartment(ctx context.Context, req CreateDepartmentRequest) (*Department, error) {
	if req.Name == "" {
		return nil, ErrDepartmentNameRequired
	}

	// 校验 enterprise 存在且 active
	ent, err := s.entRepo.GetByID(ctx, req.EnterpriseID)
	if err != nil {
		return nil, fmt.Errorf("get enterprise: %w", err)
	}
	if ent.Status != EnterpriseStatusActive {
		return nil, ErrMemberEnterpriseNotActive
	}

	dept := &Department{
		EnterpriseID: req.EnterpriseID,
		ParentID:     req.ParentID,
		Name:         req.Name,
		OrderNum:     req.OrderNum,
		Leader:       req.Leader,
		Phone:        req.Phone,
		Email:        req.Email,
		Status:       StatusActive,
	}
	if err := s.deptRepo.Create(ctx, dept); err != nil {
		return nil, fmt.Errorf("create department: %w", err)
	}
	return dept, nil
}

// UpdateDepartment 更新部门
func (s *DepartmentService) UpdateDepartment(ctx context.Context, id int64, req UpdateDepartmentRequest) (*Department, error) {
	dept, err := s.deptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get department: %w", err)
	}

	if req.ParentID != nil {
		// 防止循环引用
		if *req.ParentID == id {
			return nil, ErrDepartmentCircularRef
		}
		dept.ParentID = *req.ParentID
	}
	if req.Name != nil {
		dept.Name = *req.Name
	}
	if req.OrderNum != nil {
		dept.OrderNum = *req.OrderNum
	}
	if req.Leader != nil {
		dept.Leader = *req.Leader
	}
	if req.Phone != nil {
		dept.Phone = *req.Phone
	}
	if req.Email != nil {
		dept.Email = *req.Email
	}

	if err := s.deptRepo.Update(ctx, dept); err != nil {
		return nil, fmt.Errorf("update department: %w", err)
	}
	return dept, nil
}

// DeleteDepartment 删除部门（前置校验 + 级联）
func (s *DepartmentService) DeleteDepartment(ctx context.Context, id int64) error {
	// 1. 检查是否有子部门
	hasChildren, err := s.deptRepo.HasChildren(ctx, id)
	if err != nil {
		return fmt.Errorf("check children: %w", err)
	}
	if hasChildren {
		return ErrDepartmentHasChildrenErr
	}

	// 2. 检查是否有成员
	hasMembers, err := s.deptRepo.HasMembers(ctx, id)
	if err != nil {
		return fmt.Errorf("check members: %w", err)
	}
	if hasMembers {
		return ErrDepartmentHasMembersErr
	}

	// 3. 软删除
	if err := s.deptRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete department: %w", err)
	}
	return nil
}

// GetDepartment 获取部门详情
func (s *DepartmentService) GetDepartment(ctx context.Context, id int64) (*Department, error) {
	dept, err := s.deptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get department: %w", err)
	}
	return dept, nil
}
