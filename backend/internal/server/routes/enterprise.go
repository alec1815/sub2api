package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterEnterpriseRoutes 注册企业管理员视角路由 (P5)
// 依赖 RequireEnterpriseAdmin 中间件注入 enterprise_id 上下文
func RegisterEnterpriseRoutes(
	v1 *gin.RouterGroup,
	h *handler.Handlers,
	jwtAuth middleware.JWTAuthMiddleware,
	memberRepo service.EnterpriseMemberRepository,
	entRepo service.EnterpriseRepository,
) {
	// 企业管理员路由组：JWT + 企业管理员权限
	enterpriseGroup := v1.Group("/enterprise")
	enterpriseGroup.Use(gin.HandlerFunc(jwtAuth))
	enterpriseGroup.Use(middleware.RequireEnterpriseAdmin(memberRepo, entRepo))
	{
		// 成员管理
		enterpriseGroup.GET("/members", h.Enterprise.Member.ListMembers)
		enterpriseGroup.POST("/members", h.Enterprise.Member.CreateMember)
		enterpriseGroup.PUT("/members/:id", h.Enterprise.Member.UpdateMember)
		enterpriseGroup.DELETE("/members/:id", h.Enterprise.Member.UnbindMember)

	// 密钥管理
	enterpriseGroup.GET("/keys", h.Enterprise.Key.ListKeys)
	enterpriseGroup.POST("/keys", h.Enterprise.Key.CreateKey)
	enterpriseGroup.PUT("/keys/:id", h.Enterprise.Key.UpdateKey)
	enterpriseGroup.POST("/keys/:id/toggle", h.Enterprise.Key.ToggleKey)
	enterpriseGroup.DELETE("/keys/:id", h.Enterprise.Key.DeleteKey)

	// 企业财务
	enterpriseGroup.GET("/finance", h.Enterprise.Billing.GetFinance)
	enterpriseGroup.GET("/usage", h.Enterprise.Billing.GetUsage)
	enterpriseGroup.GET("/balance", h.Enterprise.Billing.GetBalance)
	enterpriseGroup.GET("/subscriptions", h.Enterprise.Billing.GetSubscriptions)
	enterpriseGroup.POST("/recharge", h.Enterprise.Billing.RechargeEnterprise)
	enterpriseGroup.POST("/subscribe", h.Enterprise.Billing.SubscribeEnterprise)

		// 企业设置
		enterpriseGroup.GET("/settings", h.Enterprise.Profile.GetSettings)
		enterpriseGroup.PUT("/settings", h.Enterprise.Profile.UpdateSettings)

		// 部门管理
		enterpriseGroup.GET("/departments", h.Enterprise.Department.ListDepartments)
		enterpriseGroup.POST("/departments", h.Enterprise.Department.CreateDepartment)
		enterpriseGroup.PUT("/departments/:id", h.Enterprise.Department.UpdateDepartment)
		enterpriseGroup.DELETE("/departments/:id", h.Enterprise.Department.DeleteDepartment)
	}

	// 企业成员路由组：JWT + 任意成员权限
	memberGroup := v1.Group("/enterprise")
	memberGroup.Use(gin.HandlerFunc(jwtAuth))
	memberGroup.Use(middleware.RequireEnterpriseMember(memberRepo, entRepo))
	{
		// Profile (只读)
		memberGroup.GET("/profile", h.Enterprise.Profile.GetProfile)
		// 企业仪表盘用量查询
		memberGroup.GET("/usage/models", h.Enterprise.Dashboard.GetModelStats)
		memberGroup.GET("/usage/trend", h.Enterprise.Dashboard.GetUsageTrend)
		memberGroup.GET("/dashboard/snapshot", h.Enterprise.Dashboard.GetSnapshotV2)
	}
}
