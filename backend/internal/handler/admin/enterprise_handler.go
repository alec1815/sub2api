package admin

import (
	"context"
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// EnterpriseHandler 企业 CRUD 管理 handler（平台运营方视角）
type EnterpriseHandler struct {
	enterpriseService *service.EnterpriseService
}

func NewEnterpriseHandler(enterpriseService *service.EnterpriseService) *EnterpriseHandler {
	return &EnterpriseHandler{
		enterpriseService: enterpriseService,
	}
}

// ListEnterprises GET /api/admin/enterprises
func (h *EnterpriseHandler) ListEnterprises(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}

	filters := service.EnterpriseListFilters{}
	if status := c.Query("status"); status != "" {
		filters.Status = status
	}
	if search := c.Query("search"); search != "" {
		filters.Search = search
	}

	items, paginationResult, err := h.enterpriseService.ListEnterprises(c.Request.Context(), params, filters)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, paginationResult.Total, page, pageSize)
}

// GetEnterprise GET /api/admin/enterprises/:id
func (h *EnterpriseHandler) GetEnterprise(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid enterprise ID")
		return
	}

	ent, err := h.enterpriseService.GetEnterprise(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, ent)
}

// CreateEnterprise POST /api/admin/enterprises
func (h *EnterpriseHandler) CreateEnterprise(c *gin.Context) {
	var req service.CreateEnterpriseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	executeAdminIdempotentJSON(c, "admin.enterprises.create", req,
		service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
			return h.enterpriseService.CreateEnterprise(ctx, req)
		})
}

// UpdateEnterprise PUT /api/admin/enterprises/:id
func (h *EnterpriseHandler) UpdateEnterprise(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid enterprise ID")
		return
	}

	var req service.UpdateEnterpriseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ent, err := h.enterpriseService.UpdateEnterprise(c.Request.Context(), id, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, ent)
}

// ToggleEnterpriseStatus POST /api/admin/enterprises/:id/toggle
func (h *EnterpriseHandler) ToggleEnterpriseStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid enterprise ID")
		return
	}

	ent, err := h.enterpriseService.ToggleStatus(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, ent)
}

// DeleteEnterprise DELETE /api/admin/enterprises/:id
func (h *EnterpriseHandler) DeleteEnterprise(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid enterprise ID")
		return
	}

	executeAdminIdempotentJSON(c, "admin.enterprises.delete", id,
		service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
			return nil, h.enterpriseService.DeleteEnterprise(ctx, id)
		})
}

// ==================== 企业余额管理 ====================

// UpdateEnterpriseBalanceRequest 企业余额变更请求
type UpdateEnterpriseBalanceRequest struct {
	Balance   float64 `json:"balance" binding:"required,gt=0"`
	Operation string  `json:"operation" binding:"required,oneof=set add subtract"`
	Notes     string  `json:"notes"`
}

// UpdateBalance POST /api/admin/enterprises/:id/balance
func (h *EnterpriseHandler) UpdateBalance(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid enterprise ID")
		return
	}

	var req UpdateEnterpriseBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ent, err := h.enterpriseService.UpdateBalance(c.Request.Context(), id, req.Balance, req.Operation, req.Notes)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, ent)
}

// BalanceHistory POST /api/admin/enterprises/:id/balance-history
func (h *EnterpriseHandler) GetBalanceHistory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid enterprise ID")
		return
	}

	page, pageSize := response.ParsePagination(c)
	items, total, err := h.enterpriseService.GetBalanceHistory(c.Request.Context(), id, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, total, page, pageSize)
}
