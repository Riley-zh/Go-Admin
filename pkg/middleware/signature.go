package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SignatureMiddleware 请求签名中间件
type SignatureMiddleware struct {
	secret string
}

// NewSignatureMiddleware 创建新的请求签名中间件
func NewSignatureMiddleware(secret string) *SignatureMiddleware {
	return &SignatureMiddleware{
		secret: secret,
	}
}

// GenerateSignature 生成请求签名
func (m *SignatureMiddleware) GenerateSignature(method, path, timestamp, body string) string {
	message := fmt.Sprintf("%s:%s:%s:%s", method, path, timestamp, body)
	hmac := hmac.New(sha256.New, []byte(m.secret))
	hmac.Write([]byte(message))
	return hex.EncodeToString(hmac.Sum(nil))
}

// ValidateSignature 验证请求签名
func (m *SignatureMiddleware) ValidateSignature(c *gin.Context) bool {
	// 获取请求头中的签名相关字段
	signature := c.GetHeader("X-Signature")
	timestampStr := c.GetHeader("X-Timestamp")
	nonce := c.GetHeader("X-Nonce")

	// 检查必需字段是否存在
	if signature == "" || timestampStr == "" || nonce == "" {
		return false
	}

	// 解析时间戳
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false
	}

	// 检查时间戳是否过期（允许5分钟的时间差）
	if time.Now().Unix()-timestamp > 300 {
		return false
	}

	// 读取请求体
	body, _ := c.GetRawData()
	c.Request.Body = nil // 重置请求体，以便后续处理

	// 生成期望的签名
	expectedSignature := m.GenerateSignature(
		c.Request.Method,
		c.Request.URL.Path,
		timestampStr,
		string(body),
	)

	// 比较签名
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// RequireSignature 是一个Gin中间件函数，用于要求请求签名
func (m *SignatureMiddleware) RequireSignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对于某些方法跳过签名验证
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 验证签名
		if !m.ValidateSignature(c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
			c.Abort()
			return
		}

		c.Next()
	}
}
