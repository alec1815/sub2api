package service

import (
	"context"

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
	ErrEnterpriseInvalidStatus = infraerrors.BadRequest("ENTERPRISE_STATUS_INVALID", "enterprise status is invalid")
	ErrMemberAlreadyInDept     = infraerrors.Conflict("MEMBER_ALREADY_IN_DEPT", "member already belongs to another department")
	ErrMemberInvalidRole       = infraerrors.BadRequest("MEMBER_ROLE_INVALID", "member role is invalid")
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
