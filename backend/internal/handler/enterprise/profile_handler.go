package enterprise

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ProfileHandler 企业Profile（企业成员视角，只读）
type ProfileHandler struct {
	profileService   *service.EnterpriseProfileService
	enterpriseService *service.EnterpriseService
}

func NewProfileHandler(
	profileService *service.EnterpriseProfileService,
	enterpriseService *service.EnterpriseService,
) *ProfileHandler {
	return &ProfileHandler{
		profileService:   profileService,
		enterpriseService: enterpriseService,
	}
}

// GetProfile GET /api/enterprise/profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	profile, err := h.profileService.GetProfile(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, profile)
}

// GetSettings GET /api/enterprise/settings
func (h *ProfileHandler) GetSettings(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	profile, err := h.profileService.GetProfileByEnterprise(c.Request.Context(), enterpriseID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, profile)
}

// UpdateSettings PUT /api/enterprise/settings
func (h *ProfileHandler) UpdateSettings(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		response.Forbidden(c, "Enterprise context not found")
		return
	}

	var req service.UpdateEnterpriseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ent, err := h.enterpriseService.UpdateEnterprise(c.Request.Context(), enterpriseID, req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, ent)
}
