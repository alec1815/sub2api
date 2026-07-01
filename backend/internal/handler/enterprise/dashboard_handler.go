package enterprise

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// EnterpriseDashboardHandler 企业仪表盘 handler
type EnterpriseDashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewEnterpriseDashboardHandler(dashboardService *service.DashboardService) *EnterpriseDashboardHandler {
	return &EnterpriseDashboardHandler{dashboardService: dashboardService}
}

// GetModelStats GET /api/enterprise/usage/models
func (h *EnterpriseDashboardHandler) GetModelStats(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "Enterprise context not found"})
		return
	}
	start, end := parseDateRange(c)
	stats, err := h.dashboardService.GetModelStatsForEnterprise(c.Request.Context(), start, end, enterpriseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "ERROR", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": stats})
}

// GetUsageTrend GET /api/enterprise/usage/trend
func (h *EnterpriseDashboardHandler) GetUsageTrend(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "Enterprise context not found"})
		return
	}
	start, end := parseDateRange(c)
	granularity := c.DefaultQuery("granularity", "hour")
	limit := 12
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil { limit = v }
	}
	trend, err := h.dashboardService.GetUserUsageTrendForEnterprise(c.Request.Context(), start, end, granularity, enterpriseID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "ERROR", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"trend": trend, "models": []any{}}})
}

func parseDateRange(c *gin.Context) (time.Time, time.Time) {
	end := time.Now()
	start := end.Add(-24 * time.Hour)
	if s := c.Query("start_date"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil { start = t }
	}
	if s := c.Query("end_date"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil { end = t.Add(24 * time.Hour) }
	}
	return start, end
}
