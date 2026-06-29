package enterprise

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// KeyHandler 企业密钥管理（企业管理员视角）
type KeyHandler struct {
	keyService *service.EnterpriseKeyService
}

func NewKeyHandler(keyService *service.EnterpriseKeyService) *KeyHandler {
	return &KeyHandler{keyService: keyService}
}

// ListKeys GET /api/enterprise/keys
func (h *KeyHandler) ListKeys(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}

	items, paginationResult, err := h.keyService.ListEnterpriseKeys(c.Request.Context(), enterpriseID, params)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, paginationResult.Total, page, pageSize)
}

// CreateKey POST /api/enterprise/keys
func (h *KeyHandler) CreateKey(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req service.CreateEnterpriseKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	key, err := h.keyService.CreateEnterpriseKey(c.Request.Context(), subject.UserID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Created(c, key)
}

// UpdateKey PUT /api/enterprise/keys/:id
func (h *KeyHandler) UpdateKey(c *gin.Context) {
	keyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid key ID")
		return
	}

	var req service.UpdateEnterpriseKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	key, err := h.keyService.UpdateEnterpriseKey(c.Request.Context(), keyID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, key)
}

// DeleteKey DELETE /api/enterprise/keys/:id
func (h *KeyHandler) DeleteKey(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	keyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid key ID")
		return
	}

	if err := h.keyService.DeleteEnterpriseKeyWrapper(c.Request.Context(), subject.UserID, keyID); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Key deleted"})
}
// ToggleKey POST /api/enterprise/keys/:id/toggle
func (h *KeyHandler) ToggleKey(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	keyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid key ID")
		return
	}

	key, err := h.keyService.ToggleEnterpriseKey(c.Request.Context(), subject.UserID, keyID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, key)
}
