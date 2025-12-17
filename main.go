package main

import (
	"fmt"
	"os"

	"go-admin/internal/app"
)

// @title Go Admin API
// @version 1.0
// @description 这是一个使用Go语言开发的企业级后台管理系统API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Application startup failed: %v\n", err)
		os.Exit(1)
	}
}