package handler

import (
	"net/http"
	"strconv"
	"time"

	"go-admin/internal/metrics"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// MetricsHandler 处理指标相关的HTTP请求
type MetricsHandler struct {
	collector *metrics.MetricsCollector
}

// NewMetricsHandler 创建新的指标处理器
func NewMetricsHandler(collector *metrics.MetricsCollector) *MetricsHandler {
	return &MetricsHandler{
		collector: collector,
	}
}

// GetMetrics 获取应用程序指标摘要
func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	summary := h.collector.GetMetricsSummary()
	response.Success(c, "Metrics retrieved successfully", summary)
}

// GetMetricsByPath 获取特定路径的指标
func (h *MetricsHandler) GetMetricsByPath(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		response.Error(c, http.StatusBadRequest, "Path parameter is required")
		return
	}

	metricsData := h.collector.GetMetricsByPath(path)
	response.Success(c, "Path metrics retrieved successfully", metricsData)
}

// GetMetricsByTimeRange 获取指定时间范围内的指标
func (h *MetricsHandler) GetMetricsByTimeRange(c *gin.Context) {
	// 解析查询参数
	startTimeStr := c.Query("start")
	endTimeStr := c.Query("end")

	if startTimeStr == "" || endTimeStr == "" {
		response.Error(c, http.StatusBadRequest, "Both start and end parameters are required")
		return
	}

	// 解析时间
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid start time format, use RFC3339")
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid end time format, use RFC3339")
		return
	}

	// 验证时间范围
	if startTime.After(endTime) {
		response.Error(c, http.StatusBadRequest, "Start time must be before end time")
		return
	}

	// 限制时间范围为最多24小时
	if endTime.Sub(startTime) > 24*time.Hour {
		response.Error(c, http.StatusBadRequest, "Time range cannot exceed 24 hours")
		return
	}

	metricsData := h.collector.GetMetricsByTimeRange(startTime, endTime)
	response.Success(c, "Time range metrics retrieved successfully", metricsData)
}

// ClearMetrics 清除所有指标
func (h *MetricsHandler) ClearMetrics(c *gin.Context) {
	h.collector.ClearMetrics()
	response.Success(c, "Metrics cleared successfully", nil)
}

// GetSystemMetrics 获取系统级别的指标
func (h *MetricsHandler) GetSystemMetrics(c *gin.Context) {
	summary := h.collector.GetMetricsSummary()

	// 提取系统级别的指标
	systemMetrics := map[string]interface{}{
		"total_requests":        summary.TotalRequests,
		"total_errors":          summary.TotalErrors,
		"error_rate":            summary.ErrorRate,
		"avg_response_time_ms":  summary.AvgResponseTime.Milliseconds(),
		"p95_response_time_ms":  summary.P95ResponseTime.Milliseconds(),
		"p99_response_time_ms":  summary.P99ResponseTime.Milliseconds(),
		"request_rate_per_sec":  summary.RequestRate,
		"active_connections":    summary.ActiveConnections,
		"last_updated":          time.Now(),
	}

	response.Success(c, "System metrics retrieved successfully", systemMetrics)
}

// GetTopEndpoints 获取最活跃的端点
func (h *MetricsHandler) GetTopEndpoints(c *gin.Context) {
	summary := h.collector.GetMetricsSummary()

	// 解析限制参数
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// 限制返回的端点数量
	if len(summary.TopEndpoints) > limit {
		summary.TopEndpoints = summary.TopEndpoints[:limit]
	}

	response.Success(c, "Top endpoints retrieved successfully", summary.TopEndpoints)
}

// GetRecentErrors 获取最近的错误
func (h *MetricsHandler) GetRecentErrors(c *gin.Context) {
	summary := h.collector.GetMetricsSummary()

	// 解析限制参数
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	// 限制返回的错误数量
	if len(summary.RecentErrors) > limit {
		summary.RecentErrors = summary.RecentErrors[:limit]
	}

	response.Success(c, "Recent errors retrieved successfully", summary.RecentErrors)
}

// GetHealthStatus 获取健康状态
func (h *MetricsHandler) GetHealthStatus(c *gin.Context) {
	summary := h.collector.GetMetricsSummary()

	// 计算健康状态
	status := "healthy"
	if summary.ErrorRate > 0.1 { // 错误率超过10%
		status = "degraded"
	}
	if summary.ErrorRate > 0.5 { // 错误率超过50%
		status = "unhealthy"
	}

	// 计算平均响应时间状态
	responseTimeStatus := "good"
	if summary.AvgResponseTime > 1*time.Second {
		responseTimeStatus = "slow"
	}
	if summary.AvgResponseTime > 5*time.Second {
		responseTimeStatus = "very_slow"
	}

	healthStatus := map[string]interface{}{
		"status":              status,
		"error_rate":          summary.ErrorRate,
		"avg_response_time":   summary.AvgResponseTime.String(),
		"response_time_status": responseTimeStatus,
		"total_requests":      summary.TotalRequests,
		"timestamp":           time.Now(),
	}

	// 根据健康状态设置HTTP状态码
	httpStatus := http.StatusOK
	if status == "degraded" {
		httpStatus = http.StatusServiceUnavailable
	} else if status == "unhealthy" {
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, gin.H{
		"code":    httpStatus,
		"message": "Health status retrieved successfully",
		"data":    healthStatus,
	})
}