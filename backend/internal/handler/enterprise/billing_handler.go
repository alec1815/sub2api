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
