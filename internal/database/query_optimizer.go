package database

import (
	"fmt"
	"time"

	"go-admin/internal/logger"

	"gorm.io/gorm"
)

// QueryOptimizer provides database query optimization utilities
type QueryOptimizer struct {
	db *gorm.DB
}

// NewQueryOptimizer creates a new query optimizer
func NewQueryOptimizer() *QueryOptimizer {
	return &QueryOptimizer{
		db: GetDB(),
	}
}

// WithIndexHint adds an index hint to the query
func (qo *QueryOptimizer) WithIndexHint(table, index string) *gorm.DB {
	return qo.db.Table(table).Raw(fmt.Sprintf("SELECT * FROM %s USE INDEX (%s)", table, index))
}

// WithIndexHints adds multiple index hints to the query
func (qo *QueryOptimizer) WithIndexHints(hints map[string]string) *gorm.DB {
	// This is a simplified implementation
	// In a real-world scenario, you might want to build a more complex query builder
	return qo.db
}

// ExplainQuery returns the execution plan for a query
func (qo *QueryOptimizer) ExplainQuery(query string, args ...interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	// Add EXPLAIN prefix to the query
	explainQuery := fmt.Sprintf("EXPLAIN %s", query)
	
	err := qo.db.Raw(explainQuery, args...).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	
	return results, nil
}

// AnalyzeSlowQueries analyzes slow queries and provides optimization suggestions
func (qo *QueryOptimizer) AnalyzeSlowQueries() ([]SlowQueryAnalysis, error) {
	var analyses []SlowQueryAnalysis
	
	// Get slow query log from MySQL
	var slowQueries []struct {
		Start_time    time.Time
		User_host     string
		Query_time    time.Duration
		Lock_time     time.Duration
		Rows_sent     int
		Rows_examined int
		Db            string
		Last_insert_id int
		Insert_id     int
		Server_id     int
		SQL_text      string
	}
	
	err := qo.db.Raw("SELECT * FROM mysql.slow_log WHERE start_time > ? ORDER BY start_time DESC LIMIT 100", 
		time.Now().Add(-24*time.Hour)).Scan(&slowQueries).Error
	if err != nil {
		// If slow_log table doesn't exist, return empty result
		logger.DefaultStructuredLogger().
			WithError(err).
			Warn("Failed to get slow query log, table might not exist")
		return analyses, nil
	}
	
	// Analyze each slow query
	for _, sq := range slowQueries {
		analysis := SlowQueryAnalysis{
					Query:        sq.SQL_text,
					QueryTime:    sq.Query_time,
					RowsSent:     sq.Rows_sent,
					Rows_examined: sq.Rows_examined,
					Suggestions:  qo.generateOptimizationSuggestions(sq.SQL_text, sq.Query_time, sq.Rows_examined),
				}
		analyses = append(analyses, analysis)
	}
	
	return analyses, nil
}

// generateOptimizationSuggestions generates optimization suggestions for a query
func (qo *QueryOptimizer) generateOptimizationSuggestions(query string, queryTime time.Duration, rowsExamined int) []string {
	var suggestions []string
	
	// Check if query time is too long
	if queryTime > time.Second {
		suggestions = append(suggestions, "Query takes more than 1 second to execute")
	}
	
	// Check if query examines too many rows
	if rowsExamined > 1000 {
		suggestions = append(suggestions, "Query examines too many rows, consider adding indexes")
	}
	
	// Check for missing WHERE clause
	if !containsWord(query, "WHERE") {
		suggestions = append(suggestions, "Query lacks WHERE clause, might be scanning entire table")
	}
	
	// Check for missing LIMIT clause
	if !containsWord(query, "LIMIT") && startsWithWord(query, "SELECT") {
		suggestions = append(suggestions, "Consider adding LIMIT clause to limit result set")
	}
	
	// Check for SELECT *
	if containsSubstring(query, "SELECT *") {
		suggestions = append(suggestions, "Avoid using SELECT *, specify only needed columns")
	}
	
	// Check for ORDER BY without index
	if containsWord(query, "ORDER BY") && !containsWord(query, "USE INDEX") {
		suggestions = append(suggestions, "ORDER BY without index might cause performance issues")
	}
	
	return suggestions
}

// GetTableIndexes returns all indexes for a table
func (qo *QueryOptimizer) GetTableIndexes(tableName string) ([]TableIndex, error) {
	var indexes []TableIndex
	
	err := qo.db.Raw("SHOW INDEX FROM ?", tableName).Scan(&indexes).Error
	if err != nil {
		return nil, err
	}
	
	return indexes, nil
}

// SuggestMissingIndexes analyzes a table and suggests missing indexes
func (qo *QueryOptimizer) SuggestMissingIndexes(tableName string) ([]IndexSuggestion, error) {
	var suggestions []IndexSuggestion
	
	// Get table structure
	var columns []struct {
		Field   string
		Type    string
		Null    string
		Key     string
		Default interface{}
		Extra   string
	}
	
	err := qo.db.Raw("DESCRIBE ?", tableName).Scan(&columns).Error
	if err != nil {
		return nil, err
	}
	
	// Get existing indexes
	indexes, err := qo.GetTableIndexes(tableName)
	if err != nil {
		return nil, err
	}
	
	// Create a map of indexed columns
			indexedColumns := make(map[string]bool)
			for _, index := range indexes {
				indexedColumns[index.Column_name] = true
			}
	
	// Suggest indexes for commonly queried columns
	for _, col := range columns {
		// Skip already indexed columns
		if indexedColumns[col.Field] {
			continue
		}
		
		// Suggest index for ID columns
		if endsWith(col.Field, "_id") {
			suggestions = append(suggestions, IndexSuggestion{
				ColumnName: col.Field,
				Reason:     "Foreign key column should be indexed",
				IndexType:  "BTREE",
			})
		}
		
		// Suggest index for date/time columns
		if containsSubstring(col.Type, "date") || containsSubstring(col.Type, "time") {
			suggestions = append(suggestions, IndexSuggestion{
				ColumnName: col.Field,
				Reason:     "Date/time column should be indexed for range queries",
				IndexType:  "BTREE",
			})
		}
	}
	
	return suggestions, nil
}

// OptimizeTable runs table optimization
func (qo *QueryOptimizer) OptimizeTable(tableName string) error {
	// Run ANALYZE TABLE to update statistics
	err := qo.db.Exec("ANALYZE TABLE ?", tableName).Error
	if err != nil {
		return err
	}
	
	// Run OPTIMIZE TABLE to defragment
	err = qo.db.Exec("OPTIMIZE TABLE ?", tableName).Error
	if err != nil {
		return err
	}
	
	logger.DefaultStructuredLogger().
		WithField("table", tableName).
		Info("Table optimized successfully")
	
	return nil
}

// GetQueryStats returns statistics about query performance
func (qo *QueryOptimizer) GetQueryStats() (QueryStats, error) {
	var stats QueryStats
	
	// Get total number of queries
	err := qo.db.Raw("SHOW GLOBAL STATUS LIKE 'Queries'").Scan(&stats).Error
	if err != nil {
		return stats, err
	}
	
	// Get slow queries count
	var slowQueries struct {
		Variable_name string
		Value         string
	}
	err = qo.db.Raw("SHOW GLOBAL STATUS LIKE 'Slow_queries'").Scan(&slowQueries).Error
	if err != nil {
		return stats, err
	}
	
	// Parse slow queries count
	fmt.Sscanf(slowQueries.Value, "%d", &stats.SlowQueries)
	
	// Get uptime
	var uptime struct {
		Variable_name string
		Value         string
	}
	err = qo.db.Raw("SHOW GLOBAL STATUS LIKE 'Uptime'").Scan(&uptime).Error
	if err != nil {
		return stats, err
	}
	
	// Parse uptime
	fmt.Sscanf(uptime.Value, "%d", &stats.Uptime)
	
	// Calculate queries per second
	if stats.Uptime > 0 {
		stats.QueriesPerSecond = float64(stats.TotalQueries) / float64(stats.Uptime)
	}
	
	return stats, nil
}

// Helper functions

func containsWord(s, word string) bool {
	return containsSubstring(s, " "+word+" ") || 
		containsSubstring(s, " "+word+" ") || 
		containsSubstring(s, " "+word+" ") || 
		containsSubstring(s, " "+word+" ") ||
		containsSubstring(s, " "+word+" ") ||
		containsSubstring(s, " "+word+" ") ||
		containsSubstring(s, " "+word+" ") ||
		containsSubstring(s, " "+word+" ")
}

func startsWithWord(s, word string) bool {
	return len(s) >= len(word) && s[:len(word)] == word
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr, 0) != -1
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func findSubstring(s, substr string, start int) int {
	// Simple implementation of substring search
	// In a real-world scenario, you might want to use strings.Index
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// Data structures

type SlowQueryAnalysis struct {
	Query        string
	QueryTime    time.Duration
	RowsSent     int
	Rows_examined int
	Suggestions  []string
}

type TableIndex struct {
	Table        string
	Non_unique   int
	Key_name     string
	Seq_in_index int
	Column_name  string
	Collation    string
	Cardinality  int
	Sub_part     interface{}
	Packed       interface{}
	Null         string
	Index_type   string
	Comment      string
	Index_comment string
}

type IndexSuggestion struct {
	ColumnName string
	Reason     string
	IndexType  string
}

type QueryStats struct {
	TotalQueries      int64
	SlowQueries       int64
	Uptime            int64
	QueriesPerSecond  float64
}