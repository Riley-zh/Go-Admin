package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-admin/config"
	"go-admin/internal/cache"
	"go-admin/internal/database"
	"go-admin/internal/handler"
	"go-admin/internal/logger"
	"go-admin/internal/middleware"
	"go-admin/internal/metrics"

	"github.com/gin-gonic/gin"
)

var (
	rateLimiter *middleware.RateLimiter
)

// Run initializes and starts the application
func Run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.Log); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer logger.Sync()

	// Initialize database
	_, err = database.Init(cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer database.Close()

	// Initialize cache
	cache.Init(cfg.Cache)

	// Initialize metrics collector
	metricsCollector := metrics.NewMetricsCollector()

	// Initialize rate limiter
	rateLimitConfig := middleware.DefaultRateLimitConfig()
	rateLimitConfig.Requests = 100
	rateLimitConfig.Window = 1 * time.Minute
	rateLimiter = middleware.NewRateLimiter(rateLimitConfig)

	// Initialize response cache
	cacheConfig := middleware.DefaultCacheConfig()
	cacheConfig.CacheDuration = 5 * time.Minute

	// Initialize transaction manager
	db := database.GetDB()
	transactionManager := database.NewTransactionManager(db)

	// Create gin engine
	gin.SetMode(gin.ReleaseMode)
	if cfg.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Add middlewares
	router.Use(middleware.NewRecoveryMiddleware().Handle())
	router.Use(logger.GinLogger())
	router.Use(middleware.NewErrorHandlerMiddleware().Handle())
	router.Use(middleware.QueryPerformanceMiddleware())
	router.Use(middleware.MetricsMiddleware(metricsCollector))
	router.Use(rateLimiter.Limit())
	router.Use(middleware.RequestLoggerMiddleware())
	router.Use(middleware.ResponseCacheMiddleware(cacheConfig))
	router.Use(middleware.TransactionMiddleware(transactionManager))

	// Register routes
	registerRoutes(router, metricsCollector)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exiting")
	return nil
}

func registerRoutes(router *gin.Engine, metricsCollector *metrics.MetricsCollector) {
	// Health check endpoint
	metricsHandler := handler.NewMetricsHandler(metricsCollector)
	router.GET("/health", metricsHandler.GetHealthStatus)
	router.GET("/health/detailed", metricsHandler.GetSystemMetrics)

	// Metrics endpoints
	router.GET("/metrics", metricsHandler.GetMetrics)
	router.GET("/metrics/system", metricsHandler.GetSystemMetrics)
	router.GET("/metrics/health", metricsHandler.GetHealthStatus)
	router.GET("/metrics/endpoints", metricsHandler.GetTopEndpoints)
	router.GET("/metrics/errors", metricsHandler.GetRecentErrors)
	router.DELETE("/metrics", metricsHandler.ClearMetrics)
	router.GET("/metrics/path/:path", metricsHandler.GetMetricsByPath)
	router.GET("/metrics/timerange", metricsHandler.GetMetricsByTimeRange)

	// API version 1 group
	v1 := router.Group("/api/v1")
	{
		// Auth handlers
		authHandler := handler.NewAuthHandler()
		v1.POST("/register", authHandler.Register)
		v1.POST("/login", authHandler.Login)
		v1.POST("/logout", authHandler.Logout)
		v1.POST("/refresh", authHandler.RefreshToken)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.NewJWTMiddleware().Handle())
		protected.Use(middleware.NewCSRFMiddleware().Protect())
		{
			// User handlers
			userHandler := handler.NewUserHandler()
			protected.GET("/users/:id", userHandler.GetUserByID)
			protected.PUT("/users/:id", userHandler.UpdateUser)
			protected.DELETE("/users/:id", userHandler.DeleteUser)
			protected.GET("/users", userHandler.ListUsers)
			protected.PUT("/users/change-password", userHandler.ChangePassword)

			// Role handlers
			roleHandler := handler.NewRoleHandler()
			protected.POST("/roles", roleHandler.CreateRole)
			protected.GET("/roles/:id", roleHandler.GetRoleByID)
			protected.PUT("/roles/:id", roleHandler.UpdateRole)
			protected.DELETE("/roles/:id", roleHandler.DeleteRole)
			protected.GET("/roles", roleHandler.ListRoles)
			protected.POST("/roles/assign", roleHandler.AssignRole)
			protected.POST("/roles/remove", roleHandler.RemoveRole)
			protected.GET("/users/:id/roles", roleHandler.GetRolesByUserID)

			// Permission handlers
			permissionHandler := handler.NewPermissionHandler()
			protected.POST("/permissions", permissionHandler.CreatePermission)
			protected.GET("/permissions/:id", permissionHandler.GetPermissionByID)
			protected.PUT("/permissions/:id", permissionHandler.UpdatePermission)
			protected.DELETE("/permissions/:id", permissionHandler.DeletePermission)
			protected.GET("/permissions", permissionHandler.ListPermissions)
			protected.POST("/permissions/assign", permissionHandler.AssignPermission)
			protected.POST("/permissions/remove", permissionHandler.RemovePermission)
			protected.GET("/roles/:id/permissions", permissionHandler.GetPermissionsByRoleID)
			protected.GET("/users/:id/permissions", permissionHandler.GetPermissionsByUserID)

			// Menu handlers
			menuHandler := handler.NewMenuHandler()
			protected.POST("/menus", menuHandler.CreateMenu)
			protected.GET("/menus/:id", menuHandler.GetMenuByID)
			protected.PUT("/menus/:id", menuHandler.UpdateMenu)
			protected.DELETE("/menus/:id", menuHandler.DeleteMenu)
			protected.GET("/menus", menuHandler.ListMenus)
			protected.GET("/menus/tree", menuHandler.GetMenuTree)

			// Log handlers
			logHandler := handler.NewLogHandler()
			protected.GET("/logs/:id", logHandler.GetLogByID)
			protected.GET("/logs", logHandler.ListLogs)
			protected.DELETE("/logs/:id", logHandler.DeleteLog)
			protected.POST("/logs/clear", logHandler.ClearLogs)

			// Dictionary handlers
			dictHandler := handler.NewDictionaryHandler()
			protected.POST("/dictionaries", dictHandler.CreateDictionary)
			protected.GET("/dictionaries/:dictId", dictHandler.GetDictionaryByID)
			protected.PUT("/dictionaries/:dictId", dictHandler.UpdateDictionary)
			protected.DELETE("/dictionaries/:dictId", dictHandler.DeleteDictionary)
			protected.GET("/dictionaries", dictHandler.ListDictionaries)

			// Dictionary item handlers
			// Fixed route conflict by using consistent parameter names
			protected.POST("/dictionaries/:dictId/items", dictHandler.CreateDictionaryItem)
			protected.GET("/dictionaries/:dictId/items/:itemId", dictHandler.GetDictionaryItemByID)
			protected.PUT("/dictionaries/:dictId/items/:itemId", dictHandler.UpdateDictionaryItem)
			protected.DELETE("/dictionaries/:dictId/items/:itemId", dictHandler.DeleteDictionaryItem)
			protected.GET("/dictionaries/:dictId/items", dictHandler.ListDictionaryItems)
			protected.GET("/dictionaries/:dictId/items-all", dictHandler.ListAllDictionaryItems)

			// File handlers
			fileHandler := handler.NewFileHandler()
			protected.POST("/files/upload", fileHandler.UploadFile)
			protected.GET("/files/:id", fileHandler.GetFileByID)
			protected.GET("/files", fileHandler.ListFiles)
			protected.DELETE("/files/:id", fileHandler.DeleteFile)
			protected.GET("/files/:id/download", fileHandler.DownloadFile)

			// Notification handlers
			notificationHandler := handler.NewNotificationHandler()
			protected.POST("/notifications", notificationHandler.CreateNotification)
			protected.GET("/notifications/:id", notificationHandler.GetNotificationByID)
			protected.PUT("/notifications/:id", notificationHandler.UpdateNotification)
			protected.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
			protected.GET("/notifications", notificationHandler.ListNotifications)
			protected.GET("/notifications/active", notificationHandler.GetActiveNotifications)

			// Monitor handlers
			monitorHandler := handler.NewMonitorHandler()
			protected.GET("/monitor/info", monitorHandler.GetSystemInfo)
			protected.GET("/monitor/metrics", monitorHandler.GetSystemMetrics)
			protected.GET("/monitor/recent", monitorHandler.GetRecentMetrics)

			// Task handlers
			taskHandler := handler.NewTaskHandler()
			protected.POST("/tasks", taskHandler.CreateTask)
			protected.GET("/tasks/:id", taskHandler.GetTaskByID)
			protected.PUT("/tasks/:id", taskHandler.UpdateTask)
			protected.DELETE("/tasks/:id", taskHandler.DeleteTask)
			protected.GET("/tasks", taskHandler.ListTasks)
			protected.POST("/tasks/:id/run", taskHandler.RunTaskImmediately)

			// Import/Export handlers
			importExportHandler := handler.NewImportExportHandler()
			protected.GET("/export/users", importExportHandler.ExportUsers)
			protected.POST("/import/users", importExportHandler.ImportUsers)
			protected.GET("/export/data", importExportHandler.ExportData)

			// Cache handlers
			cacheHandler := handler.NewCacheHandler()
			protected.GET("/cache/stats", cacheHandler.GetCacheStats)
			protected.POST("/cache/reset-stats", cacheHandler.ResetCacheStats)
			protected.POST("/cache/clear", cacheHandler.ClearCache)

			// Database handlers
			dbHandler := handler.NewDBHandler()
			protected.GET("/db/stats", dbHandler.GetDBStats)

			// Database performance handlers
			dbPerfHandler := handler.NewDBPerformanceHandler()
			protected.GET("/db/performance/stats", dbPerfHandler.GetQueryStats)
			protected.GET("/db/performance/slow-queries", dbPerfHandler.GetSlowQueries)
			protected.POST("/db/performance/explain", dbPerfHandler.ExplainQuery)
			protected.GET("/db/performance/indexes/:table", dbPerfHandler.GetTableIndexes)
			protected.GET("/db/performance/indexes/:table/analyze", dbPerfHandler.AnalyzeTableIndexes)
			protected.POST("/db/performance/indexes", dbPerfHandler.CreateIndex)
			protected.DELETE("/db/performance/indexes/:index", dbPerfHandler.DropIndex)
			protected.POST("/db/performance/indexes/composite", dbPerfHandler.CreateCompositeIndex)
			protected.POST("/db/performance/indexes/fulltext", dbPerfHandler.CreateFullTextIndex)
			protected.POST("/db/performance/indexes/:index/rebuild", dbPerfHandler.RebuildIndex)
			protected.POST("/db/performance/tables/:table/optimize", dbPerfHandler.OptimizeTable)
			protected.GET("/db/performance/indexes/usage", dbPerfHandler.GetIndexUsage)
			protected.GET("/db/performance/tables/:table/suggest-indexes", dbPerfHandler.SuggestMissingIndexes)

			// Log level handlers
			logLevelHandler := handler.NewLogLevelHandler()
			protected.GET("/log/level", logLevelHandler.GetLogLevel)
			protected.POST("/log/level", logLevelHandler.SetLogLevel)

			// Security handlers
			securityHandler := handler.NewSecurityHandler()
			protected.GET("/security/csrf-token", securityHandler.GetCSRFToken)
			protected.GET("/security/rate-limit-config", securityHandler.GetRateLimitConfig)

			// Metrics handlers
			metricsHandler := handler.NewMetricsHandler(metricsCollector)
			protected.GET("/metrics", metricsHandler.GetMetrics)
			protected.GET("/health", metricsHandler.GetHealthStatus)
			protected.GET("/health/detailed", metricsHandler.GetSystemMetrics)

			// Config handlers
			configHandler := handler.NewConfigHandler()
			protected.GET("/config", configHandler.GetConfig)
			protected.POST("/config/reload", configHandler.ReloadConfig)

			// Protected ping endpoint
			protected.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})
		}
	}
}
