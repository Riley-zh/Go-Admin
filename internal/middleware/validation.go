package middleware

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationMiddleware represents the validation middleware
type ValidationMiddleware struct {
	validate *validator.Validate
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		validate: validator.New(),
	}
}

// Validate is the middleware function for request parameter validation
func (m *ValidationMiddleware) Validate(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON to the provided struct
		if err := c.ShouldBindJSON(obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Validate the struct
		if err := m.validate.Struct(obj); err != nil {
			// Convert validation errors to a more readable format
			errors := make(map[string]string)
			for _, err := range err.(validator.ValidationErrors) {
				field, _ := reflect.TypeOf(obj).Elem().FieldByName(err.StructField())
				fieldName := field.Tag.Get("json")
				if fieldName == "" {
					fieldName = err.StructField()
				}

				switch err.Tag() {
				case "required":
					errors[fieldName] = "This field is required"
				case "email":
					errors[fieldName] = "Invalid email format"
				case "min":
					if isStringType(err.Kind()) {
						errors[fieldName] = "Minimum length is " + err.Param()
					} else {
						errors[fieldName] = "Minimum value is " + err.Param()
					}
				case "max":
					if isStringType(err.Kind()) {
						errors[fieldName] = "Maximum length is " + err.Param()
					} else {
						errors[fieldName] = "Maximum value is " + err.Param()
					}
				default:
					errors[fieldName] = "Invalid value for " + fieldName
				}
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			c.Abort()
			return
		}

		// Set the validated object in context
		c.Set("validated", obj)

		// Continue to next handler
		c.Next()
	}
}

// isStringType checks if the kind is a string type
func isStringType(k reflect.Kind) bool {
	return k == reflect.String
}
