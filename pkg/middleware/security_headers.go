package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware 安全头中间件
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止XSS攻击
		c.Header("X-XSS-Protection", "1; mode=block")

		// 防止点击劫持
		c.Header("X-Frame-Options", "DENY")

		// 内容类型嗅探保护
		c.Header("X-Content-Type-Options", "nosniff")

		// CSP策略
		c.Header("Content-Security-Policy", "default-src 'self'")

		// 强制HTTPS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 防止referrer泄露
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 控制浏览器功能
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}
