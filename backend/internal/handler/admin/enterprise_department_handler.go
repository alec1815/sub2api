package admin

import (
	"context"
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// EnterpriseDepartmentHandler 部门管理 handler（平台运营方视角）
type EnterpriseDepartmentHandler struct {
	deptService *service.DepartmentService
}

func NewEnterpriseDepartmentHandler(deptService *service.DepartmentService) *EnterpriseDepartmentHandler {
	return &EnterpriseDepartmentHandler{
		deptService: deptService,
	}
}

// ListDepartments GET /api/admin/departments?enterprise_id=X
func (h *EnterpriseDepartmentHandler) ListDepartments(c *gin.Context) {
	enterpriseID, err := strconv.ParseInt(c.Query("enterprise_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid enterprise_id")
		return
	}

	tree, err := h.deptService.GetTree(c.Request.Context(), enterpriseID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, tree)
}

// CreateDepartment POST /api/admin/departments
func (h *EnterpriseDepartmentHandler) CreateDepartment(c *gin.Context) {
	var req service.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	executeAdminIdempotentJSON(c, "admin.departments.create", req,
		service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
			return h.deptService.CreateDepartment(ctx, req)
		})
}

// UpdateDepartment PUT /api/admin/departments/:id
func (h *EnterpriseDepartmentHandler) UpdateDepartment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid department ID")
		return
	}

	var req service.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	executeAdminIdempotentJSON(c, "admin.departments.update", req,
		service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
			return h.deptService.UpdateDepartment(ctx, id, req)
		})
}

// DeleteDepartment DELETE /api/admin/departments/:id
func (h *EnterpriseDepartmentHandler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid department ID")
		return
	}

	executeAdminIdempotentJSON(c, "admin.departments.delete", id,
		service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
			return nil, h.deptService.DeleteDepartment(ctx, id)
		})
}
