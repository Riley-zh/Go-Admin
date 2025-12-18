package validation

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"go-admin/pkg/errors"
	"go-admin/pkg/httpclient"
	"go-admin/pkg/jsonutils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Validator provides enhanced validation capabilities
type Validator struct {
	validate *validator.Validate
	client   *httpclient.Client
}

// NewValidator creates a new enhanced validator
func NewValidator(client *httpclient.Client) *Validator {
	v := validator.New()
	
	// Register custom validators
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("url", validateURL)
	v.RegisterValidation("json", validateJSON)
	v.RegisterValidation("future_date", validateFutureDate)
	v.RegisterValidation("past_date", validatePastDate)
	
	// Register custom validation functions
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	
	return &Validator{
		validate: v,
		client:   client,
	}
}

// Validate validates a struct and returns detailed error information
func (v *Validator) Validate(obj interface{}) error {
	if err := v.validate.Struct(obj); err != nil {
		return v.formatValidationError(err)
	}
	return nil
}

// ValidateJSON validates JSON data against a schema
func (v *Validator) ValidateJSON(jsonData []byte, schema interface{}) error {
	// First validate if it's valid JSON
	if err := jsonutils.ValidateJSON(jsonData); err != nil {
		return errors.New(400, "Invalid JSON format", err.Error())
	}
	
	// Then validate against the schema
	if err := jsonutils.Unmarshal(jsonData, schema); err != nil {
		return errors.New(400, "JSON validation failed", err.Error())
	}
	
	// Finally validate the struct
	return v.Validate(schema)
}

// ValidateAPIResponse validates an API response
func (v *Validator) ValidateAPIResponse(resp *http.Response, expectedSchema interface{}) error {
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode >= 400 {
		return errors.New(resp.StatusCode, "API request failed", resp.Status)
	}
	
	// Validate JSON response
	if expectedSchema != nil {
		if err := jsonutils.DecodeJSON(resp.Body, expectedSchema); err != nil {
			return errors.New(500, "Failed to decode API response", err.Error())
		}
		
		if err := v.Validate(expectedSchema); err != nil {
			return errors.New(400, "API response validation failed", err.Error())
		}
	}
	
	return nil
}

// ValidateWithRules validates with custom rules
func (v *Validator) ValidateWithRules(obj interface{}, rules map[string]string) error {
	// First perform standard validation
	if err := v.Validate(obj); err != nil {
		return err
	}
	
	// Then apply custom rules
	for field, rule := range rules {
		if err := v.applyCustomRule(obj, field, rule); err != nil {
			return err
		}
	}
	
	return nil
}

// formatValidationError formats validation errors into a more readable format
func (v *Validator) formatValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorDetails := make(map[string]string)
		
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()
			
			var message string
			switch tag {
			case "required":
				message = "This field is required"
			case "min":
				if e.Kind() == reflect.String {
					message = fmt.Sprintf("Minimum length is %s", param)
				} else {
					message = fmt.Sprintf("Minimum value is %s", param)
				}
			case "max":
				if e.Kind() == reflect.String {
					message = fmt.Sprintf("Maximum length is %s", param)
				} else {
					message = fmt.Sprintf("Maximum value is %s", param)
				}
			case "email":
				message = "Invalid email format"
			case "password":
				message = "Password must be at least 8 characters with uppercase, lowercase, number, and special character"
			case "username":
				message = "Username must be 3-50 characters with letters, numbers, and underscores only"
			case "phone":
				message = "Invalid phone number format"
			case "url":
				message = "Invalid URL format"
			case "json":
				message = "Invalid JSON format"
			case "future_date":
				message = "Date must be in the future"
			case "past_date":
				message = "Date must be in the past"
			default:
				message = fmt.Sprintf("Invalid value for %s", field)
			}
			
			errorDetails[field] = message
		}
		
		return errors.New(400, "Validation failed", fmt.Sprintf("%+v", errorDetails))
	}
	
	return err
}

// applyCustomRule applies a custom validation rule to a field
func (v *Validator) applyCustomRule(obj interface{}, field, rule string) error {
	// Use reflection to get the field value
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	fieldVal := val.FieldByName(field)
	if !fieldVal.IsValid() {
		return nil // Field doesn't exist, skip validation
	}
	
	// Apply the rule based on its type
	switch rule {
	case "unique":
		// This would typically involve a database check
		// For now, just return nil as a placeholder
		return nil
	case "exists":
		// This would typically involve a database check
		// For now, just return nil as a placeholder
		return nil
	default:
		// Unknown rule, skip
		return nil
	}
}

// Custom validation functions
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	// Password must be at least 8 characters
	if len(password) < 8 {
		return false
	}
	
	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return false
	}
	
	// Check for at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return false
	}
	
	// Check for at least one number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return false
	}
	
	// Check for at least one special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		return false
	}
	
	return true
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	
	// Username must be 3-50 characters
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	
	// Username must contain only letters, numbers, and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	return matched
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	
	// Basic phone validation - can be enhanced based on requirements
	matched, _ := regexp.MatchString(`^\+?[0-9]{10,15}$`, phone)
	return matched
}

func validateURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	
	// Basic URL validation
	matched, _ := regexp.MatchString(`^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`, url)
	return matched
}

func validateJSON(fl validator.FieldLevel) bool {
	jsonStr := fl.Field().String()
	
	return jsonutils.ValidateJSON([]byte(jsonStr)) == nil
}

func validateFutureDate(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	
	return date.After(time.Now())
}

func validatePastDate(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	
	return date.Before(time.Now())
}

// ValidationMiddleware creates a Gin middleware for request validation
func ValidationMiddleware(validator *Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the validation target from context or use a default
		// This would typically be set by the route handler
		if validationTarget, exists := c.Get("validation_target"); exists {
			if err := c.ShouldBindJSON(validationTarget); err != nil {
				c.JSON(400, gin.H{
					"error":   "Invalid request format",
					"details": err.Error(),
				})
				c.Abort()
				return
			}
			
			if err := validator.Validate(validationTarget); err != nil {
				c.JSON(400, gin.H{
					"error":   "Validation failed",
					"details": err.Error(),
				})
				c.Abort()
				return
			}
			
			// Store the validated object in context
			c.Set("validated", validationTarget)
		}
		
		c.Next()
	}
}

// APIValidationMiddleware creates a middleware for API response validation
func APIValidationMiddleware(validator *Validator, expectedSchema interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process the request
		c.Next()
		
		// Validate the response if there was no error
		if c.Writer.Status() < 400 {
			// This is a simplified approach
			// In a real implementation, you would need to capture the response
			// and validate it before sending it to the client
			if expectedSchema != nil {
				// Placeholder for response validation
				// This would require more complex middleware to intercept responses
			}
		}
	}
}

// RequestValidator provides a convenient way to validate requests
type RequestValidator struct {
	validator *Validator
}

// NewRequestValidator creates a new request validator
func NewRequestValidator(validator *Validator) *RequestValidator {
	return &RequestValidator{
		validator: validator,
	}
}

// Validate validates a request body against a schema
func (rv *RequestValidator) Validate(c *gin.Context, schema interface{}) bool {
	if err := c.ShouldBindJSON(schema); err != nil {
		c.JSON(400, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return false
	}
	
	if err := rv.validator.Validate(schema); err != nil {
		c.JSON(400, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return false
	}
	
	return true
}

// ValidateWithRules validates with custom rules
func (rv *RequestValidator) ValidateWithRules(c *gin.Context, schema interface{}, rules map[string]string) bool {
	if !rv.Validate(c, schema) {
		return false
	}
	
	if err := rv.validator.ValidateWithRules(schema, rules); err != nil {
		c.JSON(400, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return false
	}
	
	return true
}