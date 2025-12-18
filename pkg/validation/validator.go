package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidationMiddleware 输入验证中间件
type ValidationMiddleware struct {
	validator *validator.Validate
}

// NewValidationMiddleware 创建新的输入验证中间件
func NewValidationMiddleware() *ValidationMiddleware {
	v := validator.New()

	// 注册自定义验证器
	v.RegisterValidation("password", validatePassword)

	return &ValidationMiddleware{
		validator: v,
	}
}

// ValidateStruct 验证结构体
func (m *ValidationMiddleware) ValidateStruct(s interface{}) error {
	return m.validator.Struct(s)
}

// validatePassword 验证密码强度
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// 检查长度至少为8位
	if len(password) < 8 {
		return false
	}

	// 必须包含大写字母
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}

	// 必须包含小写字母
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}

	// 必须包含数字
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	}

	// 必须包含特殊字符
	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password) {
		return false
	}

	return true
}
