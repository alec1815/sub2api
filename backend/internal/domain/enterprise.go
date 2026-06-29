package domain

import (
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

// Enterprise statuses
const (
	EnterpriseStatusActive   = "active"
	EnterpriseStatusDisabled = "disabled"
)

// Enterprise scale
const (
	EnterpriseScaleMicro  = "micro"
	EnterpriseScaleSmall  = "small"
	EnterpriseScaleMedium = "medium"
	EnterpriseScaleLarge  = "large"
)

// Enterprise member statuses
const (
	EnterpriseMemberStatusActive = "active"
	// MemberStatusUnbound defined in constants.go
)

// EnterpriseSubscription statuses
const (
	EnterpriseSubStatusActive    = "active"
	EnterpriseSubStatusExpired   = "expired"
	EnterpriseSubStatusSuspended = "suspended"
)

// Pool types for usage logs
const (
	PoolTypePersonal   = "personal"
	PoolTypeEnterprise = "enterprise"
)

// Enterprise errors
var (
	ErrEnterpriseNotFound            = infraerrors.NotFound("ENTERPRISE_NOT_FOUND", "enterprise not found")
	ErrEnterpriseNameConflict        = infraerrors.Conflict("ENTERPRISE_NAME_CONFLICT", "enterprise name already exists")
	ErrEnterpriseMemberNotFound      = infraerrors.NotFound("ENTERPRISE_MEMBER_NOT_FOUND", "enterprise member not found")
	ErrEnterpriseMemberAlreadyActive = infraerrors.Conflict(
		"ENTERPRISE_MEMBER_ALREADY_ACTIVE",
		"user is already an active member of an enterprise",
	)
	ErrEnterpriseSubNotFound = infraerrors.NotFound("ENTERPRISE_SUB_NOT_FOUND", "enterprise subscription not found")
	ErrDepartmentNotFound    = infraerrors.NotFound("DEPARTMENT_NOT_FOUND", "department not found")
	ErrDepartmentNameExist   = infraerrors.Conflict("DEPARTMENT_NAME_EXIST", "department name already exists in this enterprise")
	ErrDepartmentHasChildren = infraerrors.Conflict("DEPARTMENT_HAS_CHILDREN", "department has child departments")
	ErrDepartmentHasMembers  = infraerrors.Conflict("DEPARTMENT_HAS_MEMBERS", "department has members")
)

// Enterprise domain model
type Enterprise struct {
	ID              int64
	Name            string
	ShortName       string
	CreditCode      string
	Address         string
	Scale           string
	Industry        string
	ParentID        int64
	Status          string
	ContactName     string
	ContactPhone    string
	ContactEmail    string
	Notes           string
	Balance         float64
	TotalRecharged  float64
	AdminUserID     int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// EnterpriseMember domain model
type EnterpriseMember struct {
	ID           int64
	EnterpriseID int64
	UserID       int64
	Role         string
	Status       string
	DepartmentID *int64
	Concurrency  int
	RPMLimit     int
	Notes        string
	JoinedAt     time.Time
	UnboundAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Department domain model
type Department struct {
	ID           int64
	EnterpriseID int64
	ParentID     int64
	Name         string
	OrderNum     int
	Leader       string
	Phone        string
	Email        string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// EnterpriseSubscription domain model
type EnterpriseSubscription struct {
	ID             int64
	EnterpriseID   int64
	GroupID        int64
	PlanID         int64
	StartsAt       time.Time
	ExpiresAt      *time.Time
	Status         string
	DailyUsageUSD  float64
	WeeklyUsageUSD float64
	MonthlyUsageUSD float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// APIKeyGroup domain model (M:N join)
type APIKeyGroup struct {
	ID        int64
	APIKeyID  int64
	GroupID   int64
	CreatedAt time.Time
}
