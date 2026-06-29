package enterprise

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// DepartmentHandler 部门管理（企业管理员视角）
type DepartmentHandler struct {
	deptService *service.DepartmentService
}

func NewDepartmentHandler(deptService *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{deptService: deptService}
}

// ListDepartments GET /api/enterprise/departments
func (h *DepartmentHandler) ListDepartments(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	tree, err := h.deptService.GetTree(c.Request.Context(), enterpriseID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, tree)
}

// CreateDepartment POST /api/enterprise/departments
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	var req service.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	req.EnterpriseID = enterpriseID

	dept, err := h.deptService.CreateDepartment(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dept)
}

// UpdateDepartment PUT /api/enterprise/departments/:id
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
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

	dept, err := h.deptService.UpdateDepartment(c.Request.Context(), id, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, dept)
}

// DeleteDepartment DELETE /api/enterprise/departments/:id
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid department ID")
		return
	}

	if err := h.deptService.DeleteDepartment(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Department deleted"})
}
