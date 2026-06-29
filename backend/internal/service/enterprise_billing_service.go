package service

import (
	"context"
	"fmt"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

// ============================================================================
// EnterpriseBillingService — 企业计费与财务聚合
// ============================================================================

var (
	ErrEnterpriseRechargeAmountInvalid = infraerrors.BadRequest("RECHARGE_AMOUNT_INVALID", "recharge amount must be greater than 0")
)

// EnterpriseFinanceOverview 企业财务汇总
type EnterpriseFinanceOverview struct {
	Balance       float64                `json:"balance"`
	TotalRecharged float64               `json:"total_recharged"`
	Subscriptions []EnterpriseSubscription `json:"subscriptions"`
	MonthlyUsage  *EnterpriseUsageSummary  `json:"monthly_usage"`
}

// EnterpriseUsageSummary 企业用量汇总
type EnterpriseUsageSummary struct {
	CallCount  int64   `json:"call_count"`
	TotalCost  float64 `json:"total_cost"`
	PoolType   string  `json:"pool_type"`
	Period     string  `json:"period"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

// EnterpriseUsageQuery 企业用量查询参数
type EnterpriseUsageQuery struct {
	EnterpriseID int64
	Start        time.Time
	End          time.Time
}

// EnterpriseBillingService 企业计费服务
type EnterpriseBillingService struct {
	entRepo  EnterpriseRepository
	subRepo  EnterpriseSubscriptionRepository
	memberRepo EnterpriseMemberRepository
}

// NewEnterpriseBillingService 构造函数
func NewEnterpriseBillingService(
	entRepo EnterpriseRepository,
	subRepo EnterpriseSubscriptionRepository,
	memberRepo EnterpriseMemberRepository,
) *EnterpriseBillingService {
	return &EnterpriseBillingService{
		entRepo:    entRepo,
		subRepo:    subRepo,
		memberRepo: memberRepo,
	}
}

// GetFinanceOverview 获取企业财务汇总
func (s *EnterpriseBillingService) GetFinanceOverview(ctx context.Context, enterpriseID int64) (*EnterpriseFinanceOverview, error) {
	// 1. 企业余额
	ent, err := s.entRepo.GetByID(ctx, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("get enterprise: %w", err)
	}

	// 2. 企业套餐
	subs, err := s.subRepo.ListActive(ctx, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	// 3. 月度用量（TODO: 需要 UsageLogRepository 支持企业用量查询）
	// 当前返回空，等 P4 网关层接入后实际查询 usage_logs
	summary := s.buildMonthlyUsageSummary(enterpriseID)

	return &EnterpriseFinanceOverview{
		Balance:        ent.Balance,
		TotalRecharged: ent.TotalRecharged,
		Subscriptions:  subs,
		MonthlyUsage:   summary,
	}, nil
}

// buildMonthlyUsageSummary 构建月度用量汇总（占位）
func (s *EnterpriseBillingService) buildMonthlyUsageSummary(_ int64) *EnterpriseUsageSummary {
	now := time.Now()
	return &EnterpriseUsageSummary{
		CallCount:  0,
		TotalCost:  0,
		PoolType:   PoolTypeEnterprise,
		Period:     "monthly",
		PeriodStart: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()),
		PeriodEnd:   now,
	}
}

// GetSubscriptions 获取企业套餐列表
func (s *EnterpriseBillingService) GetSubscriptions(ctx context.Context, enterpriseID int64) ([]EnterpriseSubscription, error) {
	subs, err := s.subRepo.GetByEnterprise(ctx, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	return subs, nil
}

// GetActiveSubscriptions 获取企业有效套餐
func (s *EnterpriseBillingService) GetActiveSubscriptions(ctx context.Context, enterpriseID int64) ([]EnterpriseSubscription, error) {
	subs, err := s.subRepo.ListActive(ctx, enterpriseID)
	if err != nil {
		return nil, fmt.Errorf("list active subscriptions: %w", err)
	}
	return subs, nil
}

// RechargeBody 企业充值输入
type EnterpriseRechargeInput struct {
	EnterpriseID int64   `json:"enterprise_id"`
	Amount       float64 `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}

// ValidateRecharge 校验充值参数
func (s *EnterpriseBillingService) ValidateRecharge(input EnterpriseRechargeInput) error {
	if input.Amount <= 0 {
		return ErrEnterpriseRechargeAmountInvalid
	}
	return nil
}

// GetBalance 获取企业余额
func (s *EnterpriseBillingService) GetBalance(ctx context.Context, enterpriseID int64) (float64, error) {
	balance, err := s.entRepo.GetBalance(ctx, enterpriseID)
	if err != nil {
		return 0, fmt.Errorf("get balance: %w", err)
	}
	return balance, nil
}

// GetUsageDetail 获取企业用量明细（分页）
// TODO(P6): 需要 UsageLogRepository 支持 enterprise_id + pool_type 查询，当前返回空
func (s *EnterpriseBillingService) GetUsageDetail(_ context.Context, _ int64) (*EnterpriseUsageSummary, error) {
	summary := s.buildMonthlyUsageSummary(0)
	return summary, nil
}

// EnterpriseSubscribeRequest 企业购买套餐请求
type EnterpriseSubscribeRequest struct {
	PlanID  int64 `json:"plan_id"`
	GroupID int64 `json:"group_id"`
}

// SubscribeEnterprise 企业购买套餐
// TODO(P6): 需要接入支付系统（复用现有 PaymentService），当前返回占位
func (s *EnterpriseBillingService) SubscribeEnterprise(_ context.Context, _ int64, _ EnterpriseSubscribeRequest) (map[string]interface{}, error) {
	return map[string]interface{}{
		"message": "Enterprise subscription is under development. Payment integration pending.",
		"status":  "not_implemented",
	}, nil
}

// RechargeEnterprise 企业充值
// TODO(P6): 需要接入支付系统（复用现有 PaymentService），当前返回占位
func (s *EnterpriseBillingService) RechargeEnterprise(_ context.Context, _ int64, _ EnterpriseRechargeInput) (map[string]interface{}, error) {
	return map[string]interface{}{
		"message": "Enterprise recharge is under development. Payment integration pending.",
		"status":  "not_implemented",
	}, nil
}
