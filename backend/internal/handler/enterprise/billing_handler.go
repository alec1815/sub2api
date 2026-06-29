package enterprise

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// BillingHandler 企业财务（企业管理员视角）
type BillingHandler struct {
	billingService *service.EnterpriseBillingService
}

func NewBillingHandler(billingService *service.EnterpriseBillingService) *BillingHandler {
	return &BillingHandler{billingService: billingService}
}

// GetFinance GET /api/enterprise/finance
func (h *BillingHandler) GetFinance(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	overview, err := h.billingService.GetFinanceOverview(c.Request.Context(), enterpriseID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, overview)
}

// GetBalance GET /api/enterprise/balance
func (h *BillingHandler) GetBalance(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	balance, err := h.billingService.GetBalance(c.Request.Context(), enterpriseID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"balance": balance})
}

// GetSubscriptions GET /api/enterprise/subscriptions
func (h *BillingHandler) GetSubscriptions(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	subs, err := h.billingService.GetActiveSubscriptions(c.Request.Context(), enterpriseID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, subs)
}

// GetUsage GET /api/enterprise/usage
// 企业用量明细（分页查询），TODO(P6): 接入 UsageLogRepository 分页查询
func (h *BillingHandler) GetUsage(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	summary, err := h.billingService.GetUsageDetail(c.Request.Context(), enterpriseID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, summary)
}

// RechargeEnterprise POST /api/enterprise/recharge
// 企业充值，TODO(P6): 接入支付系统
func (h *BillingHandler) RechargeEnterprise(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	var req service.EnterpriseRechargeInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	req.EnterpriseID = enterpriseID

	if err := h.billingService.ValidateRecharge(req); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	result, err := h.billingService.RechargeEnterprise(c.Request.Context(), enterpriseID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

// SubscribeEnterprise POST /api/enterprise/subscribe
// 企业购买套餐，TODO(P6): 接入支付系统
func (h *BillingHandler) SubscribeEnterprise(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	var req service.EnterpriseSubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.billingService.SubscribeEnterprise(c.Request.Context(), enterpriseID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}
