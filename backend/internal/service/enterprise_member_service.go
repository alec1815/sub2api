package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// Enterprise member role constants
const (
	EnterpriseMemberRoleAdmin  = EnterpriseRoleAdmin
	EnterpriseMemberRoleMember = EnterpriseRoleMember
)

var (
	ErrMemberEnterpriseNotActive = infraerrors.BadRequest("ENTERPRISE_NOT_ACTIVE", "enterprise is not active")
	ErrMemberUserAlreadyMember   = infraerrors.Conflict("USER_ALREADY_MEMBER", "user is already a member of another enterprise")
	ErrMemberDepartmentInvalid   = infraerrors.BadRequest("DEPARTMENT_INVALID", "department does not belong to this enterprise")
	ErrMemberCannotUnbindSelf    = infraerrors.BadRequest("CANNOT_UNBIND_SELF", "enterprise admin cannot unbind themselves")
)

// CreateMemberRequest 创建企业成员请求（代注册）
type CreateMemberRequest struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DepartmentID *int64 `json:"department_id"`
	Concurrency  int    `json:"concurrency"`
	RPMLimit     int    `json:"rpm_limit"`
}

// UpdateMemberRequest 更新成员请求
type UpdateMemberRequest struct {
	Username     *string `json:"username"`
	Password     *string `json:"password"`
	DepartmentID *int64  `json:"department_id"`
	Concurrency  *int    `json:"concurrency"`
	RPMLimit     *int    `json:"rpm_limit"`
	Notes        *string `json:"notes"`
}

// EnterpriseMemberService 企业成员管理
type EnterpriseMemberService struct {
	memberRepo EnterpriseMemberRepository
	userRepo   UserRepository
	entRepo    EnterpriseRepository
	deptRepo   DepartmentRepository
}

// NewEnterpriseMemberService 构造函数
func NewEnterpriseMemberService(
	memberRepo EnterpriseMemberRepository,
	userRepo UserRepository,
	entRepo EnterpriseRepository,
	deptRepo DepartmentRepository,
) *EnterpriseMemberService {
	return &EnterpriseMemberService{
		memberRepo: memberRepo,
		userRepo:   userRepo,
		entRepo:    entRepo,
		deptRepo:   deptRepo,
	}
}

// CreateMember 创建企业成员（代注册）
func (s *EnterpriseMemberService) CreateMember(ctx context.Context, enterpriseID int64, req CreateMemberRequest) (*EnterpriseMember, error) {
	// 1. 检查企业状态
	ent, err := s.entRepo.GetByID(ctx, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("get enterprise: %w", err)
	}
	if ent.Status != EnterpriseStatusActive {
		return nil, ErrMemberEnterpriseNotActive
	}

	// 2. 校验 email 唯一性
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("check user email: %w", err)
	}
	if existingUser != nil {
		// 检查是否已是某企业成员
		existingMember, memberErr := s.memberRepo.GetByUserID(ctx, existingUser.ID)
		if memberErr != nil && !errors.Is(memberErr, ErrEnterpriseMemberNotFound) {
			return nil, fmt.Errorf("check existing member: %w", memberErr)
		}
		if existingMember != nil && existingMember.Status == StatusActive {
			return nil, ErrMemberUserAlreadyMember
		}
		// 已有用户但未加入任何企业 → 允许绑定
	}

	// 3. 校验 department_id 属于本企业
	if req.DepartmentID != nil && *req.DepartmentID > 0 {
		dept, deptErr := s.deptRepo.GetByID(ctx, *req.DepartmentID)
		if deptErr != nil {
			return nil, fmt.Errorf("check department: %w", deptErr)
		}
		if dept.EnterpriseID != enterpriseID {
			return nil, ErrMemberDepartmentInvalid
		}
	}

	// 4. 创建或复用 users
	var userID int64
	if existingUser != nil {
		userID = existingUser.ID
	} else {
		username := req.Username
		if username == "" {
			username = req.Email
		}
		user := &User{
			Email:    req.Email,
			Username: username,
			Role:     RoleUser,
			Status:   StatusActive,
		}
		if req.Password != "" {
			user.SetPassword(req.Password)
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("create user: %w", err)
		}
		userID = user.ID
	}

	// 5. 创建 enterprise_members 记录
	member := &EnterpriseMember{
		EnterpriseID: enterpriseID,
		UserID:       userID,
		Role:         EnterpriseMemberRoleMember,
		Status:       StatusActive,
		DepartmentID: req.DepartmentID,
		Concurrency:  req.Concurrency,
		RPMLimit:     req.RPMLimit,
		JoinedAt:     time.Now(),
	}
	if err := s.memberRepo.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("create member: %w", err)
	}

	return member, nil
}

// GetMember 获取成员详情
func (s *EnterpriseMemberService) GetMember(ctx context.Context, id int64) (*EnterpriseMember, error) {
	m, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get member: %w", err)
	}
	return m, nil
}

// ListMembers 成员列表
func (s *EnterpriseMemberService) ListMembers(
	ctx context.Context,
	enterpriseID int64,
	params pagination.PaginationParams,
	filters EnterpriseMemberListFilters,
) ([]EnterpriseMember, *pagination.PaginationResult, error) {
	members, result, err := s.memberRepo.ListByEnterprise(ctx, enterpriseID, params, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("list members: %w", err)
	}
	return members, result, nil
}

// UpdateMember 编辑成员信息
func (s *EnterpriseMemberService) UpdateMember(ctx context.Context, id int64, req UpdateMemberRequest) (*EnterpriseMember, error) {
	m, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get member: %w", err)
	}

	// 校验 department_id 属于本企业
	if req.DepartmentID != nil {
		if *req.DepartmentID > 0 {
			dept, deptErr := s.deptRepo.GetByID(ctx, *req.DepartmentID)
			if deptErr != nil {
				return nil, fmt.Errorf("check department: %w", deptErr)
			}
			if dept.EnterpriseID != m.EnterpriseID {
				return nil, ErrMemberDepartmentInvalid
			}
			m.DepartmentID = req.DepartmentID
		} else {
			// department_id = 0 表示清除
			zero := int64(0)
			m.DepartmentID = &zero
		}
	}
	if req.Concurrency != nil {
		m.Concurrency = *req.Concurrency
	}
	if req.RPMLimit != nil {
		m.RPMLimit = *req.RPMLimit
	}
	if req.Notes != nil {
		m.Notes = *req.Notes
	}

	if err := s.memberRepo.Update(ctx, m); err != nil {
		return nil, fmt.Errorf("update member: %w", err)
	}

	// 更新 users 表字段
	if req.Username != nil || req.Password != nil {
		user, userErr := s.userRepo.GetByID(ctx, m.UserID)
		if userErr != nil {
			return nil, fmt.Errorf("get user: %w", userErr)
		}
		if req.Username != nil {
			user.Username = *req.Username
		}
		if req.Password != nil {
			user.SetPassword(*req.Password)
		}
		if updateErr := s.userRepo.Update(ctx, user); updateErr != nil {
			return nil, fmt.Errorf("update user: %w", updateErr)
		}
	}

	return m, nil
}

// UnbindMember 解绑成员（级联禁用 Key）
func (s *EnterpriseMemberService) UnbindMember(ctx context.Context, id int64) error {
	m, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get member: %w", err)
	}

	// 企业管理员不能解绑自己
	if m.Role == EnterpriseMemberRoleAdmin {
		return ErrMemberCannotUnbindSelf
	}

	// 解绑：status → unbound
	if err := s.memberRepo.Unbind(ctx, id); err != nil {
		return fmt.Errorf("unbind member: %w", err)
	}

	// 级联禁用分配给该成员的所有 Key → P4 网关层处理
	// 此处仅标记 member 状态

	return nil
}

// GetMemberByUserID 根据平台用户 ID 获取企业成员信息
func (s *EnterpriseMemberService) GetMemberByUserID(ctx context.Context, userID int64) (*EnterpriseMember, error) {
	m, err := s.memberRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get member by user: %w", err)
	}
	return m, nil
}
