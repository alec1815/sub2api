package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/domain"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// Enterprise constants aliases
const (
	EnterpriseStatusActive   = domain.EnterpriseStatusActive
	EnterpriseStatusDisabled = domain.EnterpriseStatusDisabled

	EnterpriseSubStatusActive    = domain.EnterpriseSubStatusActive
	EnterpriseSubStatusExpired   = domain.EnterpriseSubStatusExpired
	EnterpriseSubStatusSuspended = domain.EnterpriseSubStatusSuspended

	PoolTypePersonal   = domain.PoolTypePersonal
	PoolTypeEnterprise = domain.PoolTypeEnterprise
)

// Enterprise-type aliases
type Enterprise = domain.Enterprise
type EnterpriseMember = domain.EnterpriseMember
type EnterpriseSubscription = domain.EnterpriseSubscription
type Department = domain.Department
type APIKeyGroup = domain.APIKeyGroup

// Enterprise error aliases
var (
	ErrEnterpriseNotFound            = domain.ErrEnterpriseNotFound
	ErrEnterpriseNameConflict        = domain.ErrEnterpriseNameConflict
	ErrEnterpriseMemberNotFound      = domain.ErrEnterpriseMemberNotFound
	ErrEnterpriseMemberAlreadyActive = domain.ErrEnterpriseMemberAlreadyActive
	ErrEnterpriseSubNotFound         = domain.ErrEnterpriseSubNotFound
	ErrDepartmentNotFound            = domain.ErrDepartmentNotFound
	ErrDepartmentNameExist           = domain.ErrDepartmentNameExist
	ErrDepartmentHasChildren         = domain.ErrDepartmentHasChildren
	ErrDepartmentHasMembers          = domain.ErrDepartmentHasMembers
)

// Service-layer input errors
var (
	ErrEnterpriseNameRequired = infraerrors.BadRequest("ENTERPRISE_NAME_REQUIRED", "enterprise name is required")
	ErrEnterpriseAdminRequired = infraerrors.BadRequest("ENTERPRISE_ADMIN_REQUIRED", "enterprise admin user is required")
	ErrEnterpriseInvalidStatus       = infraerrors.BadRequest("ENTERPRISE_STATUS_INVALID", "enterprise status is invalid")
	ErrEnterpriseInsufficientBalance = infraerrors.BadRequest("ENTERPRISE_INSUFFICIENT_BALANCE", "enterprise balance is insufficient")
	ErrMemberAlreadyInDept           = infraerrors.Conflict("MEMBER_ALREADY_IN_DEPT", "member already belongs to another department")
	ErrMemberInvalidRole             = infraerrors.BadRequest("MEMBER_ROLE_INVALID", "member role is invalid")
	ErrDepartmentNameRequired  = infraerrors.BadRequest("DEPARTMENT_NAME_REQUIRED", "department name is required")
)

// EnterpriseListFilters defines list query filters for enterprises.
type EnterpriseListFilters struct {
	Status string
	Search string
}

// EnterpriseRepository interface
type EnterpriseRepository interface {
	Create(ctx context.Context, e *Enterprise) error
	GetByID(ctx context.Context, id int64) (*Enterprise, error)
	List(ctx context.Context, params pagination.PaginationParams, filters EnterpriseListFilters) ([]Enterprise, *pagination.PaginationResult, error)
	Update(ctx context.Context, e *Enterprise) error
	SoftDelete(ctx context.Context, id int64) error
	GetBalance(ctx context.Context, id int64) (float64, error)
	// DeductBalance 扣减企业余额，返回 newBalance
	DeductBalance(ctx context.Context, id int64, amount float64) (float64, error)
}

// EnterpriseMemberListFilters defines list query filters for enterprise members.
type EnterpriseMemberListFilters struct {
	Role         string
	Status       string
	DepartmentID *int64
}

// EnterpriseMemberRepository interface
type EnterpriseMemberRepository interface {
	Create(ctx context.Context, m *EnterpriseMember) error
	GetByID(ctx context.Context, id int64) (*EnterpriseMember, error)
	GetByUserID(ctx context.Context, userID int64) (*EnterpriseMember, error)
	ListByEnterprise(ctx context.Context, enterpriseID int64, params pagination.PaginationParams, filters EnterpriseMemberListFilters) ([]EnterpriseMember, *pagination.PaginationResult, error)
	Update(ctx context.Context, m *EnterpriseMember) error
	Unbind(ctx context.Context, id int64) error
}

// EnterpriseSubscriptionRepository interface
type EnterpriseSubscriptionRepository interface {
	Create(ctx context.Context, s *EnterpriseSubscription) error
	GetByID(ctx context.Context, id int64) (*EnterpriseSubscription, error)
	GetByEnterprise(ctx context.Context, enterpriseID int64) ([]EnterpriseSubscription, error)
	ListActive(ctx context.Context, enterpriseID int64) ([]EnterpriseSubscription, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}

// APIKeyGroupRepository interface
type APIKeyGroupRepository interface {
	SetGroups(ctx context.Context, apiKeyID int64, groupIDs []int64) error
	GetByKeyID(ctx context.Context, apiKeyID int64) ([]APIKeyGroup, error)
	GetByGroupID(ctx context.Context, groupID int64) ([]APIKeyGroup, error)
	DeleteByKeyID(ctx context.Context, apiKeyID int64) error
}

// ============================================================================
// EnterpriseService — 企业账号管理
// ============================================================================

// CreateEnterpriseRequest 创建企业请求
type CreateEnterpriseRequest struct {
	Name         string `json:"name"`
	ShortName    string `json:"short_name"`
	CreditCode   string `json:"credit_code"`
	Address      string `json:"address"`
	Scale        string `json:"scale"`
	Industry     string `json:"industry"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	ContactEmail string `json:"contact_email"`
	AdminEmail   string `json:"admin_email"`
	AdminName    string `json:"admin_name"`
	Notes        string `json:"notes"`
	ParentID     int64  `json:"parent_id"`
}

// UpdateEnterpriseRequest 更新企业请求（全部指针=可选）
type UpdateEnterpriseRequest struct {
	Name         *string `json:"name"`
	ShortName    *string `json:"short_name"`
	CreditCode   *string `json:"credit_code"`
	Address      *string `json:"address"`
	Scale        *string `json:"scale"`
	Industry     *string `json:"industry"`
	ContactName  *string `json:"contact_name"`
	ContactPhone *string `json:"contact_phone"`
	ContactEmail *string `json:"contact_email"`
	Notes        *string `json:"notes"`
}

// EnterpriseService 企业业务逻辑层
type EnterpriseService struct {
	entRepo    EnterpriseRepository
	memberRepo EnterpriseMemberRepository
	subRepo    EnterpriseSubscriptionRepository
	userRepo   UserRepository
	deptRepo   DepartmentRepository
}

// NewEnterpriseService 构造函数
func NewEnterpriseService(
	entRepo EnterpriseRepository,
	memberRepo EnterpriseMemberRepository,
	subRepo EnterpriseSubscriptionRepository,
	userRepo UserRepository,
	deptRepo DepartmentRepository,
) *EnterpriseService {
	return &EnterpriseService{
		entRepo:    entRepo,
		memberRepo: memberRepo,
		subRepo:    subRepo,
		userRepo:   userRepo,
		deptRepo:   deptRepo,
	}
}

// CreateEnterprise 创建企业（含管理员自动注册）
func (s *EnterpriseService) CreateEnterprise(ctx context.Context, req CreateEnterpriseRequest) (*Enterprise, error) {
	if req.Name == "" {
		return nil, ErrEnterpriseNameRequired
	}
	if req.AdminEmail == "" {
		return nil, ErrEnterpriseAdminRequired
	}

	// 1. 查找或创建管理员 users
	admin, err := s.userRepo.GetByEmail(ctx, req.AdminEmail)
	if err != nil {
		return nil, fmt.Errorf("lookup admin user: %w", err)
	}
	if admin == nil {
		// 自动注册管理员用户
		adminName := req.AdminName
		if adminName == "" {
			adminName = req.AdminEmail
		}
		admin = &User{
			Email:    req.AdminEmail,
			Username: adminName,
			Role:     RoleUser, // platform role = user, enterprise role = enterprise_admin
			Status:   StatusActive,
		}
		if err := s.userRepo.Create(ctx, admin); err != nil {
			return nil, fmt.Errorf("create admin user: %w", err)
		}
	}

	// 2. 校验管理员未被其他企业绑定为管理员
	existing, err := s.memberRepo.GetByUserID(ctx, admin.ID)
	if err != nil && err != ErrEnterpriseMemberNotFound {
		return nil, fmt.Errorf("check existing member: %w", err)
	}
	if existing != nil && (existing.Role == EnterpriseRoleAdmin || existing.Status == StatusActive) {
		return nil, ErrEnterpriseMemberAlreadyActive
	}

	// 3. 创建 enterprise
	ent := &Enterprise{
		Name:         req.Name,
		ShortName:    req.ShortName,
		CreditCode:   req.CreditCode,
		Address:      req.Address,
		Scale:        req.Scale,
		Industry:     req.Industry,
		ParentID:     req.ParentID,
		Status:       EnterpriseStatusActive,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		Notes:        req.Notes,
		Balance:      0,
		TotalRecharged: 0,
		AdminUserID:  admin.ID,
	}
	if err := s.entRepo.Create(ctx, ent); err != nil {
		return nil, fmt.Errorf("create enterprise: %w", err)
	}

	// 4. 创建 enterprise_members 记录
	member := &EnterpriseMember{
		EnterpriseID: ent.ID,
		UserID:       admin.ID,
		Role:         EnterpriseRoleAdmin,
		Status:       StatusActive,
		JoinedAt:     time.Now(),
	}
	if err := s.memberRepo.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("create admin member: %w", err)
	}

	return ent, nil
}

// GetEnterprise 获取企业详情
func (s *EnterpriseService) GetEnterprise(ctx context.Context, id int64) (*Enterprise, error) {
	ent, err := s.entRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get enterprise: %w", err)
	}
	return ent, nil
}

// ListEnterprises 企业列表
func (s *EnterpriseService) ListEnterprises(ctx context.Context, params pagination.PaginationParams, filters EnterpriseListFilters) ([]Enterprise, *pagination.PaginationResult, error) {
	ents, result, err := s.entRepo.List(ctx, params, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("list enterprises: %w", err)
	}
	return ents, result, nil
}

// UpdateEnterprise 更新企业信息
func (s *EnterpriseService) UpdateEnterprise(ctx context.Context, id int64, req UpdateEnterpriseRequest) (*Enterprise, error) {
	ent, err := s.entRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get enterprise: %w", err)
	}

	if req.Name != nil {
		ent.Name = *req.Name
	}
	if req.ShortName != nil {
		ent.ShortName = *req.ShortName
	}
	if req.CreditCode != nil {
		ent.CreditCode = *req.CreditCode
	}
	if req.Address != nil {
		ent.Address = *req.Address
	}
	if req.Scale != nil {
		ent.Scale = *req.Scale
	}
	if req.Industry != nil {
		ent.Industry = *req.Industry
	}
	if req.ContactName != nil {
		ent.ContactName = *req.ContactName
	}
	if req.ContactPhone != nil {
		ent.ContactPhone = *req.ContactPhone
	}
	if req.ContactEmail != nil {
		ent.ContactEmail = *req.ContactEmail
	}
	if req.Notes != nil {
		ent.Notes = *req.Notes
	}

	if err := s.entRepo.Update(ctx, ent); err != nil {
		return nil, fmt.Errorf("update enterprise: %w", err)
	}
	return ent, nil
}

// ToggleStatus 启停企业（active ↔ disabled）
func (s *EnterpriseService) ToggleStatus(ctx context.Context, id int64) (*Enterprise, error) {
	ent, err := s.entRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get enterprise: %w", err)
	}

	newStatus := EnterpriseStatusDisabled
	if ent.Status == EnterpriseStatusDisabled {
		newStatus = EnterpriseStatusActive
	}
	ent.Status = newStatus

	if err := s.entRepo.Update(ctx, ent); err != nil {
		return nil, fmt.Errorf("toggle enterprise status: %w", err)
	}

	// 停用时级联：禁用该企业所有分配的 Key 由 P4 网关层处理
	// 此处仅更新企业状态

	return ent, nil
}

// DeleteEnterprise 删除企业（软删除）
// 前置条件：status 必须已为 disabled
// 级联操作：解绑所有成员、标记套餐过期
func (s *EnterpriseService) DeleteEnterprise(ctx context.Context, id int64) error {
	ent, err := s.entRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get enterprise: %w", err)
	}

	// 1. 前置校验：必须先停用
	if ent.Status != EnterpriseStatusDisabled {
		return infraerrors.BadRequest("ENTERPRISE_NOT_DISABLED", "enterprise must be disabled before deletion")
	}

	// 2. 级联解绑所有 enterprise_members
	members, _, err := s.memberRepo.ListByEnterprise(ctx, id, pagination.DefaultPagination(), EnterpriseMemberListFilters{
		Status: StatusActive,
	})
	if err != nil {
		return fmt.Errorf("list members: %w", err)
	}
	for _, m := range members {
		if unbindErr := s.memberRepo.Unbind(ctx, m.ID); unbindErr != nil {
			return fmt.Errorf("unbind member %d: %w", m.ID, unbindErr)
		}
	}

	// 3. 企业套餐标记 expired
	activeSubs, err := s.subRepo.ListActive(ctx, id)
	if err != nil {
		return fmt.Errorf("list active subs: %w", err)
	}
	for _, sub := range activeSubs {
		if subErr := s.subRepo.UpdateStatus(ctx, sub.ID, EnterpriseSubStatusExpired); subErr != nil {
			return fmt.Errorf("expire subscription %d: %w", sub.ID, subErr)
		}
	}

	// 4. 软删除企业
	if err := s.entRepo.SoftDelete(ctx, id); err != nil {
		return fmt.Errorf("soft delete enterprise: %w", err)
	}

	return nil
}

// GetEnterpriseBalance 获取企业余额
func (s *EnterpriseService) GetEnterpriseBalance(ctx context.Context, id int64) (float64, error) {
	balance, err := s.entRepo.GetBalance(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("get enterprise balance: %w", err)
	}
	return balance, nil
}
