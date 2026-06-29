package service

import (
	"context"
	"fmt"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

// ============================================================================
// EnterpriseProfileService — 企业 Profile 页
// ============================================================================

var (
	ErrProfileNotEnterpriseMember = infraerrors.Forbidden("NOT_ENTERPRISE_MEMBER", "user is not a member of this enterprise")
)

// EnterpriseProfile 企业 Profile 数据
type EnterpriseProfile struct {
	Enterprise   *Enterprise `json:"enterprise"`
	Member       *EnterpriseMember `json:"member"`
	Department   *Department `json:"department,omitempty"`
	MonthlyUsage *MonthlyUsageOverview `json:"monthly_usage"`
}

// MonthlyUsageOverview 月度用量概览
type MonthlyUsageOverview struct {
	CallCount  int64   `json:"call_count"`
	TotalCost  float64 `json:"total_cost"`
}

// EnterpriseProfileService 企业 Profile 服务
type EnterpriseProfileService struct {
	entRepo    EnterpriseRepository
	memberRepo EnterpriseMemberRepository
	deptRepo   DepartmentRepository
}

// NewEnterpriseProfileService 构造函数
func NewEnterpriseProfileService(
	entRepo EnterpriseRepository,
	memberRepo EnterpriseMemberRepository,
	deptRepo DepartmentRepository,
) *EnterpriseProfileService {
	return &EnterpriseProfileService{
		entRepo:    entRepo,
		memberRepo: memberRepo,
		deptRepo:   deptRepo,
	}
}

// GetProfile 获取企业 Profile（当前用户视角）
func (s *EnterpriseProfileService) GetProfile(ctx context.Context, userID int64) (*EnterpriseProfile, error) {
	// 1. 获取用户的企业成员身份
	member, err := s.memberRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get member: %w", err)
	}
	if member.Status != StatusActive {
		return nil, ErrProfileNotEnterpriseMember
	}

	// 2. 获取企业信息
	ent, err := s.entRepo.GetByID(ctx, member.EnterpriseID)
	if err != nil {
		return nil, fmt.Errorf("get enterprise: %w", err)
	}

	// 3. 获取部门信息（如果有）
	// TODO(P4): 月度用量需网关计费层接入后通过 enterprise_id 聚合 usage_logs 实现
	profile := &EnterpriseProfile{
		Enterprise: ent,
		Member:     member,
		MonthlyUsage: &MonthlyUsageOverview{
			CallCount: 0,
			TotalCost: 0,
		},
	}

	if member.DepartmentID != nil && *member.DepartmentID > 0 {
		dept, deptErr := s.deptRepo.GetByID(ctx, *member.DepartmentID)
		if deptErr == nil {
			profile.Department = dept
		}
	}

	return profile, nil
}

// GetProfileByEnterprise 获取指定企业的 Profile（以企业 ID 视角）
func (s *EnterpriseProfileService) GetProfileByEnterprise(ctx context.Context, enterpriseID int64) (*EnterpriseProfile, error) {
	ent, err := s.entRepo.GetByID(ctx, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("get enterprise: %w", err)
	}

	// TODO(P4): 月度用量需网关计费层接入后通过 enterprise_id 聚合 usage_logs 实现
	return &EnterpriseProfile{
		Enterprise: ent,
		MonthlyUsage: &MonthlyUsageOverview{
			CallCount: 0,
			TotalCost: 0,
		},
	}, nil
}

// GetProfileWithMember 获取企业 Profile（指定成员视角）
func (s *EnterpriseProfileService) GetProfileWithMember(ctx context.Context, memberID int64) (*EnterpriseProfile, error) {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("get member: %w", err)
	}

	ent, err := s.entRepo.GetByID(ctx, member.EnterpriseID)
	if err != nil {
		return nil, fmt.Errorf("get enterprise: %w", err)
	}

	// TODO(P4): 月度用量需网关计费层接入后通过 enterprise_id 聚合 usage_logs 实现
	profile := &EnterpriseProfile{
		Enterprise: ent,
		Member:     member,
		MonthlyUsage: &MonthlyUsageOverview{
			CallCount: 0,
			TotalCost: 0,
		},
	}

	if member.DepartmentID != nil && *member.DepartmentID > 0 {
		dept, deptErr := s.deptRepo.GetByID(ctx, *member.DepartmentID)
		if deptErr == nil {
			profile.Department = dept
		}
	}

	return profile, nil
}


