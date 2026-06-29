package service

import (
	"context"
	"fmt"
)

// EnterpriseBillingRouter 企业计费分流器。
// 在网关层判断 API Key 属于个人池还是企业池。
type EnterpriseBillingRouter struct {
	memberRepo EnterpriseMemberRepository
	entRepo    EnterpriseRepository
}

// NewEnterpriseBillingRouter 创建企业计费分流器
func NewEnterpriseBillingRouter(
	memberRepo EnterpriseMemberRepository,
	entRepo EnterpriseRepository,
) *EnterpriseBillingRouter {
	return &EnterpriseBillingRouter{
		memberRepo: memberRepo,
		entRepo:    entRepo,
	}
}

// BillingSource 计费来源
type BillingSource struct {
	PoolType     string // "personal" / "enterprise"
	EnterpriseID *int64 // 企业 ID (pool_type=enterprise 时)
	PayerUserID  int64  // 付费用户 ID (personal → key.UserID, enterprise → 企业ID 映射)
}

// DetermineBillingSource 根据 API Key 的 AssignedTo 字段判定计费来源。
//
// 规则：
//   - AssignedTo == nil → personal pool（扣创建者个人余额/套餐）
//   - AssignedTo != nil → enterprise pool（扣企业资金池）
func (r *EnterpriseBillingRouter) DetermineBillingSource(
	ctx context.Context,
	apiKey *APIKey,
) (*BillingSource, error) {
	// 个人 Key / 管理员自用 Key → 扣创建者个人余额
	if apiKey.AssignedTo == nil {
		return &BillingSource{
			PoolType:    PoolTypePersonal,
			PayerUserID: apiKey.UserID,
		}, nil
	}

	// 企业分配的 Key → 查询成员 → 查询企业
	member, err := r.memberRepo.GetByID(ctx, *apiKey.AssignedTo)
	if err != nil {
		return nil, fmt.Errorf("enterprise member %d not found: %w", *apiKey.AssignedTo, err)
	}

	enterprise, err := r.entRepo.GetByID(ctx, member.EnterpriseID)
	if err != nil {
		return nil, fmt.Errorf("enterprise %d not found: %w", member.EnterpriseID, err)
	}

	if enterprise.Status != EnterpriseStatusActive {
		return nil, fmt.Errorf("enterprise %d is not active", enterprise.ID)
	}

	return &BillingSource{
		PoolType:     PoolTypeEnterprise,
		EnterpriseID: &enterprise.ID,
		PayerUserID:  apiKey.UserID, // 创建 Key 的管理员
	}, nil
}
