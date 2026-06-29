package middleware

import (
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ContextKey 扩展 —— 企业认证上下文键
const (
	// ContextKeyEnterpriseID 当前用户所属企业ID (int64)
	ContextKeyEnterpriseID ContextKey = "enterprise_id"
	// ContextKeyEnterpriseMemberID 当前用户的企业成员记录ID (int64)
	ContextKeyEnterpriseMemberID ContextKey = "enterprise_member_id"
	// ContextKeyEnterpriseRole 当前用户在企业中的角色 ("enterprise_admin" / "enterprise_member")
	ContextKeyEnterpriseRole ContextKey = "enterprise_role"
)

// GetEnterpriseIDFromContext 从 gin.Context 获取企业ID
func GetEnterpriseIDFromContext(c *gin.Context) (int64, bool) {
	value, exists := c.Get(string(ContextKeyEnterpriseID))
	if !exists {
		return 0, false
	}
	id, ok := value.(int64)
	return id, ok
}

// GetEnterpriseMemberIDFromContext 从 gin.Context 获取企业成员ID
func GetEnterpriseMemberIDFromContext(c *gin.Context) (int64, bool) {
	value, exists := c.Get(string(ContextKeyEnterpriseMemberID))
	if !exists {
		return 0, false
	}
	id, ok := value.(int64)
	return id, ok
}

// GetEnterpriseRoleFromContext 从 gin.Context 获取企业角色
func GetEnterpriseRoleFromContext(c *gin.Context) (string, bool) {
	value, exists := c.Get(string(ContextKeyEnterpriseRole))
	if !exists {
		return "", false
	}
	role, ok := value.(string)
	return role, ok
}

// RequireEnterpriseAdmin 企业管理员权限中间件。
// 必须在 JWTAuth 中间件之后使用。
// 验证当前用户是某激活企业的 enterprise_admin。
func RequireEnterpriseAdmin(
	memberRepo service.EnterpriseMemberRepository,
	entRepo service.EnterpriseRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject, ok := GetAuthSubjectFromContext(c)
		if !ok {
			AbortWithError(c, 401, "UNAUTHORIZED", "User not authenticated")
			return
		}

		member, err := memberRepo.GetByUserID(c.Request.Context(), subject.UserID)
		if err != nil {
			AbortWithError(c, 403, "FORBIDDEN", "Not an enterprise member")
			return
		}

		if member.Role != service.EnterpriseRoleAdmin {
			AbortWithError(c, 403, "FORBIDDEN", "Enterprise admin access required")
			return
		}

		if member.Status != service.StatusActive {
			AbortWithError(c, 403, "FORBIDDEN", "Enterprise member is not active")
			return
		}

		enterprise, err := entRepo.GetByID(c.Request.Context(), member.EnterpriseID)
		if err != nil {
			AbortWithError(c, 403, "FORBIDDEN", "Enterprise not found")
			return
		}

		if enterprise.Status != service.EnterpriseStatusActive {
			AbortWithError(c, 403, "FORBIDDEN", "Enterprise is not active")
			return
		}

		c.Set(string(ContextKeyEnterpriseID), member.EnterpriseID)
		c.Set(string(ContextKeyEnterpriseMemberID), member.ID)
		c.Set(string(ContextKeyEnterpriseRole), member.Role)
		c.Next()
	}
}

// RequireEnterpriseMember 企业成员权限中间件。
// 必须在 JWTAuth 中间件之后使用。
// 验证当前用户是某激活企业的成员（不限角色）。
func RequireEnterpriseMember(
	memberRepo service.EnterpriseMemberRepository,
	entRepo service.EnterpriseRepository,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		subject, ok := GetAuthSubjectFromContext(c)
		if !ok {
			AbortWithError(c, 401, "UNAUTHORIZED", "User not authenticated")
			return
		}

		member, err := memberRepo.GetByUserID(c.Request.Context(), subject.UserID)
		if err != nil {
			AbortWithError(c, 403, "FORBIDDEN", "Not an enterprise member")
			return
		}

		if member.Status != service.StatusActive {
			AbortWithError(c, 403, "FORBIDDEN", "Enterprise member is not active")
			return
		}

		enterprise, err := entRepo.GetByID(c.Request.Context(), member.EnterpriseID)
		if err != nil {
			AbortWithError(c, 403, "FORBIDDEN", "Enterprise not found")
			return
		}

		if enterprise.Status != service.EnterpriseStatusActive {
			AbortWithError(c, 403, "FORBIDDEN", "Enterprise is not active")
			return
		}

		c.Set(string(ContextKeyEnterpriseID), member.EnterpriseID)
		c.Set(string(ContextKeyEnterpriseMemberID), member.ID)
		c.Set(string(ContextKeyEnterpriseRole), member.Role)
		c.Next()
	}
}
