package enterprise

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// MemberHandler 企业成员管理（企业管理员视角）
type MemberHandler struct {
	memberService *service.EnterpriseMemberService
}

func NewMemberHandler(memberService *service.EnterpriseMemberService) *MemberHandler {
	return &MemberHandler{memberService: memberService}
}

// enterpriseIDFromContext 从中间件注入的 context 获取企业ID
func enterpriseIDFromContext(c *gin.Context) (int64, bool) {
	return middleware.GetEnterpriseIDFromContext(c)
}

// ListMembers GET /api/enterprise/members
func (h *MemberHandler) ListMembers(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
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

// CreateMember POST /api/enterprise/members
func (h *MemberHandler) CreateMember(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	var req service.CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	member, err := h.memberService.CreateMember(c.Request.Context(), enterpriseID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, member)
}

// UpdateMember PUT /api/enterprise/members/:id
func (h *MemberHandler) UpdateMember(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid member ID")
		return
	}

	var req service.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	member, err := h.memberService.UpdateMember(c.Request.Context(), memberID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, member)
}

// UnbindMember DELETE /api/enterprise/members/:id
func (h *MemberHandler) UnbindMember(c *gin.Context) {
	memberID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid member ID")
		return
	}

	if err := h.memberService.UnbindMember(c.Request.Context(), memberID); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"message": "Member unbound"})
}
