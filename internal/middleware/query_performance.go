package middleware

import (
	"context"
	"time"

	"go-admin/internal/database"
	"go-admin/internal/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// QueryPerformanceMiddleware tracks database query performance
func QueryPerformanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with query tracking
		ctx := context.WithValue(c.Request.Context(), "queryTracker", &QueryTracker{
			queries: make([]QueryInfo, 0),
		})
		c.Request = c.Request.WithContext(ctx)

		// Get database instance
		db := database.GetDB()
		if db == nil {
			c.Next()
			return
		}

		// Register callback to track queries
		var queryStartTime time.Time
		
		beforeCallback := func(db *gorm.DB) {
			queryStartTime = time.Now()
		}
		
		afterCallback := func(db *gorm.DB) {
			if tracker, ok := db.Statement.Context.Value("queryTracker").(*QueryTracker); ok {
				tracker.AddQuery(db.Statement.SQL.String(), db.Statement.Vars, time.Since(queryStartTime))
			}
		}

		// Register the callback
		db.Callback().Query().Before("gorm:query").Register("track_query_before", beforeCallback)
		db.Callback().Query().After("gorm:query").Register("track_query_after", afterCallback)

		// Process request
		c.Next()

		// Log query performance
		if tracker, ok := c.Request.Context().Value("queryTracker").(*QueryTracker); ok {
			tracker.LogPerformance(c)
		}

		// Unregister callbacks to avoid memory leaks
		db.Callback().Query().Remove("track_query_before")
		db.Callback().Query().Remove("track_query_after")
	}
}

// QueryTracker tracks database queries during a request
type QueryTracker struct {
	queries []QueryInfo
}

// QueryInfo contains information about a database query
type QueryInfo struct {
	SQL       string
	Vars      []interface{}
	Duration  time.Duration
	Timestamp time.Time
}

// AddQuery adds a query to the tracker
func (qt *QueryTracker) AddQuery(sql string, vars []interface{}, duration time.Duration) {
	qt.queries = append(qt.queries, QueryInfo{
		SQL:       sql,
		Vars:      vars,
		Duration:  duration,
		Timestamp: time.Now(),
	})
}

// GetQueries returns all tracked queries
func (qt *QueryTracker) GetQueries() []QueryInfo {
	return qt.queries
}

// GetTotalDuration returns the total duration of all queries
func (qt *QueryTracker) GetTotalDuration() time.Duration {
	var total time.Duration
	for _, query := range qt.queries {
		total += query.Duration
	}
	return total
}

// GetSlowQueries returns queries that took longer than the threshold
func (qt *QueryTracker) GetSlowQueries(threshold time.Duration) []QueryInfo {
	var slowQueries []QueryInfo
	for _, query := range qt.queries {
		if query.Duration > threshold {
			slowQueries = append(slowQueries, query)
		}
	}
	return slowQueries
}

// LogPerformance logs the performance of tracked queries
func (qt *QueryTracker) LogPerformance(c *gin.Context) {
	if len(qt.queries) == 0 {
		return
	}

	// Calculate statistics
	totalDuration := qt.GetTotalDuration()
	avgDuration := totalDuration / time.Duration(len(qt.queries))
	slowQueries := qt.GetSlowQueries(100 * time.Millisecond) // Queries taking more than 100ms are considered slow

	// Log performance metrics
	logger.DefaultStructuredLogger().
		WithField("path", c.Request.URL.Path).
		WithField("method", c.Request.Method).
		WithField("query_count", len(qt.queries)).
		WithField("total_duration_ms", totalDuration.Milliseconds()).
		WithField("avg_duration_ms", avgDuration.Milliseconds()).
		WithField("slow_query_count", len(slowQueries)).
		Info("Database query performance")

	// Log slow queries
	for _, query := range slowQueries {
		logger.DefaultStructuredLogger().
			WithField("path", c.Request.URL.Path).
			WithField("method", c.Request.Method).
			WithField("sql", query.SQL).
			WithField("duration_ms", query.Duration.Milliseconds()).
			Warn("Slow database query detected")
	}
}