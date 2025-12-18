package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware provides CORS configuration for the application
func CORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		// Allow all origins in development
		AllowAllOrigins: true,
		// Allow common HTTP methods
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		// Allow common headers
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"X-CSRF-Token",
		},
		// Expose headers to frontend
		ExposeHeaders: []string{
			"Content-Length",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
		},
		// Allow credentials
		AllowCredentials: true,
	}

	return cors.New(config)
}
