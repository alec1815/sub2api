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
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	ShortName       string     `json:"short_name"`
	CreditCode      string     `json:"credit_code"`
	Address         string     `json:"address"`
	Scale           string     `json:"scale"`
	Industry        string     `json:"industry"`
	ParentID        int64      `json:"parent_id"`
	Status          string     `json:"status"`
	ContactName     string     `json:"contact_name"`
	ContactPhone    string     `json:"contact_phone"`
	ContactEmail    string     `json:"contact_email"`
	Notes           string     `json:"notes"`
	Balance         float64    `json:"balance"`
	TotalRecharged  float64    `json:"total_recharged"`
	AdminUserID     int64      `json:"admin_user_id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// EnterpriseMember domain model
type EnterpriseMember struct {
	ID             int64      `json:"id"`
	EnterpriseID   int64      `json:"enterprise_id"`
	UserID         int64      `json:"user_id"`
	// from users join
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	LastActiveAt   *time.Time  `json:"last_active_at,omitempty"`
	Role           string     `json:"role"`
	Status         string     `json:"status"`
	DepartmentID   *int64     `json:"department_id"`
	// from departments join
	DepartmentName string     `json:"department_name"`
	Concurrency    int        `json:"concurrency"`
	RPMLimit       int        `json:"rpm_limit"`
	Notes          string     `json:"notes"`
	JoinedAt       time.Time  `json:"joined_at"`
	UnboundAt      *time.Time `json:"unbound_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// Department domain model
type Department struct {
	ID           int64     `json:"id"`
	EnterpriseID int64     `json:"enterprise_id"`
	ParentID     int64     `json:"parent_id"`
	Name         string    `json:"name"`
	OrderNum     int       `json:"order_num"`
	Leader       string    `json:"leader"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// EnterpriseSubscription domain model
type EnterpriseSubscription struct {
	ID              int64      `json:"id"`
	EnterpriseID    int64      `json:"enterprise_id"`
	GroupID         int64      `json:"group_id"`
	PlanID          int64      `json:"plan_id"`
	StartsAt        time.Time  `json:"starts_at"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	Status          string     `json:"status"`
	DailyUsageUSD   float64    `json:"daily_usage_usd"`
	WeeklyUsageUSD  float64    `json:"weekly_usage_usd"`
	MonthlyUsageUSD float64    `json:"monthly_usage_usd"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// APIKeyGroup domain model (M:N join)
type APIKeyGroup struct {
	ID        int64
	APIKeyID  int64
	GroupID   int64
	CreatedAt time.Time
}
