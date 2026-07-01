package enterprise

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/gin-gonic/gin"
)

// GetSnapshotV2 GET /api/enterprise/dashboard/snapshot
func (h *EnterpriseDashboardHandler) GetSnapshotV2(c *gin.Context) {
	enterpriseID, ok := enterpriseIDFromContext(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"code": "FORBIDDEN", "message": "Enterprise context not found"})
		return
	}

	start, end := parseEnterpriseTimeRange(c)
	granularity := c.DefaultQuery("granularity", "day")
	if granularity != "hour" {
		granularity = "day"
	}

	limit := 12
	if l := c.Query("users_trend_limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 50 {
			limit = v
		}
	}

	// Trend
	var trend []usagestats.TrendDataPoint
	if parseBoolQ(c.Query("include_trend"), true) {
		trend, _ = h.dashboardService.GetUsageTrendWithFilters(c.Request.Context(), start, end, granularity, 0, 0, 0, 0, enterpriseID, "", nil, nil, nil)
	}

	// Models
	var models []usagestats.ModelStat
	if parseBoolQ(c.Query("include_model_stats"), true) {
		models, _ = h.dashboardService.GetModelStatsWithFilters(c.Request.Context(), start, end, 0, 0, 0, 0, enterpriseID, nil, nil, nil)
	}

	// Users trend (Top 12)
	var usersTrend []usagestats.UserUsageTrendPoint
	if parseBoolQ(c.Query("include_users_trend"), true) {
		usersTrend, _ = h.dashboardService.GetUserUsageTrendForEnterprise(c.Request.Context(), start, end, granularity, enterpriseID, limit)
	}

	// Aggregated stats from trend+models (按企业维度)
	stats := aggregateEnterpriseStats(enterpriseID, trend, models, usersTrend)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"generated_at":  time.Now().UTC().Format(time.RFC3339),
			"start_date":    start.Format("2006-01-02"),
			"end_date":      end.Add(-24 * time.Hour).Format("2006-01-02"),
			"granularity":   granularity,
			"stats":         stats,
			"trend":         trend,
			"models":        models,
			"users_trend":   usersTrend,
			"enterprise_id": enterpriseID,
		},
	})
}

func aggregateEnterpriseStats(enterpriseID int64, trend []usagestats.TrendDataPoint, models []usagestats.ModelStat, usersTrend []usagestats.UserUsageTrendPoint) *usagestats.DashboardStats {
	stats := &usagestats.DashboardStats{
		ActiveUsers: int64(len(usersTrend)),
	}
	for _, t := range trend {
		stats.TotalRequests += t.Requests
		stats.TotalInputTokens += t.InputTokens
		stats.TotalOutputTokens += t.OutputTokens
		stats.TotalCacheCreationTokens += t.CacheCreationTokens
		stats.TotalCacheReadTokens += t.CacheReadTokens
		stats.TotalTokens += t.TotalTokens
		stats.TotalCost += t.Cost
		stats.TotalActualCost += t.ActualCost
	}
	for _, m := range models {
		stats.TotalRequests += m.Requests
		stats.TotalInputTokens += m.InputTokens
		stats.TotalOutputTokens += m.OutputTokens
		stats.TotalCacheCreationTokens += m.CacheCreationTokens
		stats.TotalCacheReadTokens += m.CacheReadTokens
		stats.TotalTokens += m.TotalTokens
		stats.TotalCost += m.Cost
		stats.TotalActualCost += m.ActualCost
	}
	// Today = 全部
	stats.TodayRequests = stats.TotalRequests
	stats.TodayInputTokens = stats.TotalInputTokens
	stats.TodayOutputTokens = stats.TotalOutputTokens
	stats.TodayCacheCreationTokens = stats.TotalCacheCreationTokens
	stats.TodayCacheReadTokens = stats.TotalCacheReadTokens
	stats.TodayTokens = stats.TotalTokens
	stats.TodayCost = stats.TotalCost
	stats.TodayActualCost = stats.TotalActualCost
	return stats
}

func parseBoolQ(s string, def bool) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return def
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return def
	}
	return v
}

// parseEnterpriseTimeRange 解析 start_date/end_date/timezone 参数
func parseEnterpriseTimeRange(c *gin.Context) (time.Time, time.Time) {
	userTZ := c.Query("timezone")
	now := timezone.NowInUserLocation(userTZ)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var startTime, endTime time.Time
	if startDate != "" {
		if t, err := timezone.ParseInUserLocation("2006-01-02", startDate, userTZ); err == nil {
			startTime = t
		}
	}
	if endDate != "" {
		if t, err := timezone.ParseInUserLocation("2006-01-02", endDate, userTZ); err == nil {
			endTime = t.Add(24 * time.Hour)
		}
	}

	if startTime.IsZero() && endTime.IsZero() {
		endTime = now
		startTime = now.Add(-24 * time.Hour)
	} else if startTime.IsZero() {
		startTime = endTime.AddDate(0, 0, -7)
	} else if endTime.IsZero() {
		endTime = now
	}
	return startTime, endTime
}
