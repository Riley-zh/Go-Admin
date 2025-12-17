package middleware

import (
	"net/http"
	"strings"

	"go-admin/internal/database"
	"go-admin/internal/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TransactionMiddleware 自动管理HTTP请求中的数据库事务
func TransactionMiddleware(transactionManager *database.TransactionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对需要事务的请求进行处理
		if !shouldUseTransaction(c) {
			c.Next()
			return
		}

		// 在事务中处理请求
		result, err := transactionManager.WithTransaction(c.Request.Context(), func(tx *gorm.DB) error {
			// 将事务存储到上下文中，以便处理器使用
			c.Set("db", tx)

			// 处理请求
			c.Next()

			// 如果响应已经写入且状态码表示错误，则返回错误以触发回滚
			if c.Writer.Written() && c.Writer.Status() >= 400 {
				return &TransactionError{
					StatusCode: c.Writer.Status(),
					Message:    "Request failed with status code " + string(rune(c.Writer.Status())),
				}
			}

			return nil
		})

		// 记录事务结果
		if err != nil {
			logger.DefaultStructuredLogger().
				WithError(err).
				WithField("path", c.Request.URL.Path).
				WithField("method", c.Request.Method).
				WithField("success", result.Success).
				WithField("retries", result.Retries).
				WithField("duration_ms", result.Duration.Milliseconds()).
				Error("Transaction failed")

			// 如果响应尚未写入，则写入错误响应
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal server error",
					"data":    nil,
				})
			}
		} else {
			logger.DefaultStructuredLogger().
				WithField("path", c.Request.URL.Path).
				WithField("method", c.Request.Method).
				WithField("success", result.Success).
				WithField("retries", result.Retries).
				WithField("duration_ms", result.Duration.Milliseconds()).
				Info("Transaction completed successfully")
		}
	}
}

// TransactionMiddlewareWithConfig 带配置的事务中间件
func TransactionMiddlewareWithConfig(transactionManager *database.TransactionManager, config TransactionConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否应该使用事务
		if !shouldUseTransactionWithConfig(c, config) {
			c.Next()
			return
		}

		// 在事务中处理请求
		result, err := transactionManager.WithTransaction(c.Request.Context(), func(tx *gorm.DB) error {
			// 将事务存储到上下文中，以便处理器使用
			c.Set("db", tx)

			// 处理请求
			c.Next()

			// 如果响应已经写入且状态码表示错误，则返回错误以触发回滚
			if c.Writer.Written() && c.Writer.Status() >= 400 {
				return &TransactionError{
					StatusCode: c.Writer.Status(),
					Message:    "Request failed with status code " + string(rune(c.Writer.Status())),
				}
			}

			return nil
		}, config.Options)

		// 记录事务结果
		if err != nil {
			logger.DefaultStructuredLogger().
				WithError(err).
				WithField("path", c.Request.URL.Path).
				WithField("method", c.Request.Method).
				WithField("success", result.Success).
				WithField("retries", result.Retries).
				WithField("duration_ms", result.Duration.Milliseconds()).
				Error("Transaction failed")

			// 如果响应尚未写入，则写入错误响应
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal server error",
					"data":    nil,
				})
			}
		} else {
			logger.DefaultStructuredLogger().
				WithField("path", c.Request.URL.Path).
				WithField("method", c.Request.Method).
				WithField("success", result.Success).
				WithField("retries", result.Retries).
				WithField("duration_ms", result.Duration.Milliseconds()).
				Info("Transaction completed successfully")
		}
	}
}

// TransactionConfig 事务中间件配置
type TransactionConfig struct {
	// Options 事务选项
	Options database.TransactionOptions
	
	// IncludePaths 包含的路径模式列表
	IncludePaths []string
	
	// ExcludePaths 排除的路径模式列表
	ExcludePaths []string
	
	// IncludeMethods 包含的HTTP方法列表
	IncludeMethods []string
	
	// ExcludeMethods 排除的HTTP方法列表
	ExcludeMethods []string
	
	// ReadOnlyPaths 只读事务的路径模式列表
	ReadOnlyPaths []string
}

// DefaultTransactionConfig 返回默认的事务配置
func DefaultTransactionConfig() TransactionConfig {
	return TransactionConfig{
		Options: database.DefaultTransactionOptions(),
		IncludePaths: []string{
			"/api/v1/users",
			"/api/v1/roles",
			"/api/v1/permissions",
		},
		ExcludePaths: []string{
			"/api/v1/auth/login",
			"/api/v1/auth/register",
			"/health",
			"/metrics",
		},
		IncludeMethods: []string{
			"POST",
			"PUT",
			"DELETE",
			"PATCH",
		},
		ExcludeMethods: []string{
			"GET",
			"HEAD",
			"OPTIONS",
		},
		ReadOnlyPaths: []string{
			"/api/v1/users",
			"/api/v1/roles",
		},
	}
}

// shouldUseTransaction 检查是否应该使用事务
func shouldUseTransaction(c *gin.Context) bool {
	// 只对写操作使用事务
	method := c.Request.Method
	if method != "POST" && method != "PUT" && method != "DELETE" && method != "PATCH" {
		return false
	}

	// 排除健康检查和指标端点
	path := c.Request.URL.Path
	if path == "/health" || path == "/metrics" || strings.HasPrefix(path, "/api/v1/auth") {
		return false
	}

	// 对API路径使用事务
	return strings.HasPrefix(path, "/api/v1/")
}

// shouldUseTransactionWithConfig 根据配置检查是否应该使用事务
func shouldUseTransactionWithConfig(c *gin.Context, config TransactionConfig) bool {
	method := c.Request.Method
	path := c.Request.URL.Path

	// 检查方法是否包含
	if len(config.IncludeMethods) > 0 {
		included := false
		for _, m := range config.IncludeMethods {
			if method == m {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	// 检查方法是否排除
	for _, m := range config.ExcludeMethods {
		if method == m {
			return false
		}
	}

	// 检查路径是否包含
	if len(config.IncludePaths) > 0 {
		included := false
		for _, p := range config.IncludePaths {
			if strings.HasPrefix(path, p) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	// 检查路径是否排除
	for _, p := range config.ExcludePaths {
		if strings.HasPrefix(path, p) {
			return false
		}
	}

	// 检查是否为只读事务
	if len(config.ReadOnlyPaths) > 0 {
		for _, p := range config.ReadOnlyPaths {
			if strings.HasPrefix(path, p) {
				// 使用只读事务
				return true
			}
		}
	}

	return true
}

// TransactionError 事务错误
type TransactionError struct {
	StatusCode int
	Message    string
}

func (e *TransactionError) Error() string {
	return e.Message
}

// GetTransactionFromContext 从上下文中获取事务
func GetTransactionFromContext(c *gin.Context) (*gorm.DB, bool) {
	db, exists := c.Get("db")
	if !exists {
		return nil, false
	}

	tx, ok := db.(*gorm.DB)
	return tx, ok
}

// WithTransaction 在处理器中使用事务
func WithTransaction(c *gin.Context, fn func(*gorm.DB) error) error {
	tx, exists := GetTransactionFromContext(c)
	if !exists {
		// 如果没有事务，使用普通数据库连接
		db := database.GetDB()
		return fn(db)
	}

	return fn(tx)
}