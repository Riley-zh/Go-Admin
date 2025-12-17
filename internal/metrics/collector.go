package metrics

import (
	"context"
	"sync"
	"time"
)

// MetricsCollector 收集和存储应用程序性能指标
type MetricsCollector struct {
	requestMetrics map[string]*RequestMetrics
	mu             sync.RWMutex
}

// RequestMetrics 存储单个请求的指标
type RequestMetrics struct {
	Path           string
	Method         string
	StatusCode     int
	Duration       time.Duration
	Timestamp      time.Time
	ResponseSize   int64
	RequestSize    int64
	UserAgent      string
	ClientIP       string
	UserID         string
}

// SystemMetrics 存储系统级别的指标
type SystemMetrics struct {
	TotalRequests       int64
	TotalErrors         int64
	AverageResponseTime time.Duration
	P95ResponseTime     time.Duration
	P99ResponseTime     time.Duration
	ActiveConnections   int64
	LastUpdated         time.Time
}

// MetricsSummary 提供指标摘要
type MetricsSummary struct {
	TotalRequests      int64                    `json:"total_requests"`
	TotalErrors        int64                    `json:"total_errors"`
	ErrorRate          float64                  `json:"error_rate"`
	AvgResponseTime    time.Duration            `json:"avg_response_time"`
	P95ResponseTime    time.Duration            `json:"p95_response_time"`
	P99ResponseTime    time.Duration            `json:"p99_response_time"`
	TopEndpoints       []EndpointStats          `json:"top_endpoints"`
	RecentErrors       []ErrorMetric            `json:"recent_errors"`
	RequestRate        float64                  `json:"request_rate_per_second"`
	ActiveConnections  int64                    `json:"active_connections"`
	SystemMetrics      map[string]interface{}   `json:"system_metrics"`
}

// EndpointStats 提供端点统计信息
type EndpointStats struct {
	Path         string        `json:"path"`
	Method       string        `json:"method"`
	RequestCount int64         `json:"request_count"`
	AvgTime      time.Duration `json:"avg_time"`
	ErrorCount   int64         `json:"error_count"`
	ErrorRate    float64       `json:"error_rate"`
}

// ErrorMetric 存储错误指标
type ErrorMetric struct {
	Timestamp    time.Time `json:"timestamp"`
	Path         string    `json:"path"`
	Method       string    `json:"method"`
	StatusCode   int       `json:"status_code"`
	ErrorMessage string    `json:"error_message"`
	ClientIP     string    `json:"client_ip"`
	UserID       string    `json:"user_id"`
}

// NewMetricsCollector 创建新的指标收集器
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		requestMetrics: make(map[string]*RequestMetrics),
	}
}

// RecordRequest 记录请求指标
func (mc *MetricsCollector) RecordRequest(ctx context.Context, metrics *RequestMetrics) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := metrics.Path + ":" + metrics.Method
	mc.requestMetrics[key] = metrics

	// 不记录日志，避免与 RequestLoggerMiddleware 重复
	// 指标收集只用于内部统计，不输出到日志
}

// GetMetricsSummary 获取指标摘要
func (mc *MetricsCollector) GetMetricsSummary() *MetricsSummary {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	// 计算端点统计
	endpointStats := make(map[string]*EndpointStats)
	var totalRequests int64
	var totalErrors int64
	var totalDuration time.Duration
	var durations []time.Duration

	for _, metrics := range mc.requestMetrics {
		key := metrics.Path + ":" + metrics.Method
		if _, exists := endpointStats[key]; !exists {
			endpointStats[key] = &EndpointStats{
				Path:   metrics.Path,
				Method: metrics.Method,
			}
		}

		stat := endpointStats[key]
		stat.RequestCount++
		totalDuration += metrics.Duration
		durations = append(durations, metrics.Duration)

		if metrics.StatusCode >= 400 {
			stat.ErrorCount++
			totalErrors++
		}

		totalRequests++
	}

	// 计算平均值和百分位数
	var avgTime time.Duration
	var p95Time time.Duration
	var p99Time time.Duration

	if len(durations) > 0 {
		avgTime = totalDuration / time.Duration(len(durations))
		p95Time = calculatePercentile(durations, 0.95)
		p99Time = calculatePercentile(durations, 0.99)
	}

	// 获取前10个最活跃的端点
	topEndpoints := make([]EndpointStats, 0, len(endpointStats))
	for _, stat := range endpointStats {
		if stat.RequestCount > 0 {
			stat.AvgTime = totalDuration / time.Duration(stat.RequestCount)
			stat.ErrorRate = float64(stat.ErrorCount) / float64(stat.RequestCount)
			topEndpoints = append(topEndpoints, *stat)
		}
	}

	// 计算错误率
	errorRate := 0.0
	if totalRequests > 0 {
		errorRate = float64(totalErrors) / float64(totalRequests)
	}

	// 计算请求速率（每秒）
	requestRate := 0.0
	if len(mc.requestMetrics) > 0 {
		// 简化计算：假设指标是在过去60秒内收集的
		requestRate = float64(totalRequests) / 60.0
	}

	return &MetricsSummary{
		TotalRequests:     totalRequests,
		TotalErrors:       totalErrors,
		ErrorRate:         errorRate,
		AvgResponseTime:   avgTime,
		P95ResponseTime:   p95Time,
		P99ResponseTime:   p99Time,
		TopEndpoints:      topEndpoints,
		RequestRate:       requestRate,
		ActiveConnections: 0, // 需要从其他地方获取
		SystemMetrics:     make(map[string]interface{}),
	}
}

// calculatePercentile 计算百分位数
func calculatePercentile(durations []time.Duration, percentile float64) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	// 简单排序（实际应用中应使用更高效的算法）
	for i := 0; i < len(durations); i++ {
		for j := i + 1; j < len(durations); j++ {
			if durations[i] > durations[j] {
				durations[i], durations[j] = durations[j], durations[i]
			}
		}
	}

	index := int(float64(len(durations)) * percentile)
	if index >= len(durations) {
		index = len(durations) - 1
	}

	return durations[index]
}

// ClearMetrics 清除所有指标
func (mc *MetricsCollector) ClearMetrics() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.requestMetrics = make(map[string]*RequestMetrics)
	// 不记录日志，避免不必要的输出
}

// GetMetricsByPath 获取特定路径的指标
func (mc *MetricsCollector) GetMetricsByPath(path string) []*RequestMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	var result []*RequestMetrics
	for _, metrics := range mc.requestMetrics {
		if metrics.Path == path {
			result = append(result, metrics)
		}
	}

	return result
}

// GetMetricsByTimeRange 获取指定时间范围内的指标
func (mc *MetricsCollector) GetMetricsByTimeRange(start, end time.Time) []*RequestMetrics {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	var result []*RequestMetrics
	for _, metrics := range mc.requestMetrics {
		if metrics.Timestamp.After(start) && metrics.Timestamp.Before(end) {
			result = append(result, metrics)
		}
	}

	return result
}