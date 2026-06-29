package admin

import (
	"context"
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// EnterpriseMemberHandler 企业成员管理 handler（平台运营方视角）
type EnterpriseMemberHandler struct {
	memberService *service.EnterpriseMemberService
}

func NewEnterpriseMemberHandler(memberService *service.EnterpriseMemberService) *EnterpriseMemberHandler {
	return &EnterpriseMemberHandler{
		memberService: memberService,
	}
}

// ListMembers GET /api/admin/enterprises/:id/members
func (h *EnterpriseMemberHandler) ListMembers(c *gin.Context) {
	enterpriseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid enterprise ID")
		return
	}

	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}

	items, paginationResult, err := h.memberService.ListMembers(c.Request.Context(), enterpriseID, params, service.EnterpriseMemberListFilters{})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, items, paginationResult.Total, page, pageSize)
}

// CreateMember POST /api/admin/enterprises/:id/members
func (h *EnterpriseMemberHandler) CreateMember(c *gin.Context) {
	enterpriseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid enterprise ID")
		return
	}

	var req service.CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	executeAdminIdempotentJSON(c, "admin.enterprises.members.create", req,
		service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
			return h.memberService.CreateMember(ctx, enterpriseID, req)
		})
}

// UpdateMember PUT /api/admin/enterprises/:id/members/:mid
func (h *EnterpriseMemberHandler) UpdateMember(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("mid"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid member ID")
		return
	}

	var req service.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	executeAdminIdempotentJSON(c, "admin.enterprises.members.update", req,
		service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
			return h.memberService.UpdateMember(ctx, memberID, req)
		})
}

// UnbindMember DELETE /api/admin/enterprises/:id/members/:mid
func (h *EnterpriseMemberHandler) UnbindMember(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("mid"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid member ID")
		return
	}

	executeAdminIdempotentJSON(c, "admin.enterprises.members.unbind", memberID,
		service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
			return nil, h.memberService.UnbindMember(ctx, memberID)
		})
}
