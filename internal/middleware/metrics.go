package middleware

import (
	"bytes"
	"strconv"
	"time"

	"go-admin/internal/metrics"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware 收集请求指标
func MetricsMiddleware(collector *metrics.MetricsCollector) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 读取请求体大小
		var requestSize int64
		if c.Request.Body != nil {
			requestSize = c.Request.ContentLength
		}

		// 创建响应写入器包装器以捕获响应大小
		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// 处理请求
		c.Next()

		// 计算请求持续时间
		duration := time.Since(startTime)

		// 获取用户ID（如果存在）
		userID := ""
		if userIDValue, exists := c.Get("user_id"); exists {
			if id, ok := userIDValue.(string); ok {
				userID = id
			}
		}

		// 创建请求指标
		requestMetrics := &metrics.RequestMetrics{
			Path:         c.Request.URL.Path,
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			Duration:     duration,
			Timestamp:    startTime,
			ResponseSize: int64(responseWriter.body.Len()),
			RequestSize:  requestSize,
			UserAgent:    c.Request.UserAgent(),
			ClientIP:     c.ClientIP(),
			UserID:       userID,
		}

		// 记录指标
		collector.RecordRequest(c.Request.Context(), requestMetrics)

		// 不在这里记录日志，避免与 RequestLoggerMiddleware 重复
		// 日志记录统一由 RequestLoggerMiddleware 处理
	}
}

// responseBodyWriter 包装 gin.ResponseWriter 以捕获响应体
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// MetricsMiddlewareWithConfig 带配置的指标中间件
func MetricsMiddlewareWithConfig(collector *metrics.MetricsCollector, config MetricsConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过健康检查端点
		if config.SkipHealthCheck && (c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ping") {
			c.Next()
			return
		}

		// 跳过静态文件
		if config.SkipStaticFiles && isStaticFile(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 记录请求开始时间
		startTime := time.Now()

		// 读取请求体大小
		var requestSize int64
		if c.Request.Body != nil {
			requestSize = c.Request.ContentLength
		}

		// 创建响应写入器包装器以捕获响应大小
		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = responseWriter

		// 处理请求
		c.Next()

		// 计算请求持续时间
		duration := time.Since(startTime)

		// 获取用户ID（如果存在）
		userID := ""
		if userIDValue, exists := c.Get("user_id"); exists {
			if id, ok := userIDValue.(string); ok {
				userID = id
			}
		}

		// 创建请求指标
		requestMetrics := &metrics.RequestMetrics{
			Path:         c.Request.URL.Path,
			Method:       c.Request.Method,
			StatusCode:   c.Writer.Status(),
			Duration:     duration,
			Timestamp:    startTime,
			ResponseSize: int64(responseWriter.body.Len()),
			RequestSize:  requestSize,
			UserAgent:    c.Request.UserAgent(),
			ClientIP:     c.ClientIP(),
			UserID:       userID,
		}

		// 记录指标
		collector.RecordRequest(c.Request.Context(), requestMetrics)

		// 不在这里记录日志，避免与 RequestLoggerMiddleware 重复
		// 日志记录统一由 RequestLoggerMiddleware 处理

		// 添加响应头
		if config.AddResponseHeaders {
			c.Header("X-Response-Time", strconv.FormatInt(duration.Milliseconds(), 10)+"ms")
			c.Header("X-Request-ID", getRequestID(c))
		}
	}
}

// MetricsConfig 指标中间件配置
type MetricsConfig struct {
	SlowRequestThreshold time.Duration
	SkipHealthCheck      bool
	SkipStaticFiles      bool
	AddResponseHeaders   bool
}

// DefaultMetricsConfig 返回默认配置
func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		SlowRequestThreshold: 2 * time.Second,
		SkipHealthCheck:      true,
		SkipStaticFiles:      true,
		AddResponseHeaders:   true,
	}
}

// isStaticFile 检查路径是否为静态文件
func isStaticFile(path string) bool {
	return len(path) > 4 && (
		path[len(path)-4:] == ".css" ||
		path[len(path)-3:] == ".js" ||
		path[len(path)-4:] == ".png" ||
		path[len(path)-4:] == ".jpg" ||
		path[len(path)-5:] == ".jpeg" ||
		path[len(path)-4:] == ".gif" ||
		path[len(path)-4:] == ".ico" ||
		path[len(path)-4:] == ".svg" ||
		path[len(path)-5:] == ".woff" ||
		path[len(path)-5:] == ".ttf")
}

// getRequestID 获取请求ID
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}