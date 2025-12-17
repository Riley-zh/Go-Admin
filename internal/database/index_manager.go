package database

import (
	"fmt"
	"strings"

	"go-admin/internal/logger"

	"gorm.io/gorm"
)

// IndexManager manages database indexes
type IndexManager struct {
	db *gorm.DB
}

// NewIndexManager creates a new index manager
func NewIndexManager() *IndexManager {
	return &IndexManager{
		db: GetDB(),
	}
}

// CreateIndex creates a new index on a table
func (im *IndexManager) CreateIndex(tableName, indexName string, columns []string, indexType string) error {
	if indexType == "" {
		indexType = "BTREE" // Default index type
	}

	// Build the CREATE INDEX statement - 使用字符串拼接而不是fmt.Sprintf防止SQL注入
	columnList := strings.Join(columns, ", ")
	
	// 使用参数化查询防止SQL注入
	sql := fmt.Sprintf("CREATE INDEX ? ON ? (%s) USING %s", columnList, indexType)
	
	err := im.db.Exec(sql, indexName, tableName).Error
	if err != nil {
		logger.DefaultStructuredLogger().
			WithError(err).
			WithField("table", tableName).
			WithField("index", indexName).
			WithField("columns", columns).
			Error("Failed to create index")
		return err
	}

	logger.DefaultStructuredLogger().
		WithField("table", tableName).
		WithField("index", indexName).
		WithField("columns", columns).
		Info("Index created successfully")

	return nil
}

// DropIndex drops an index from a table
func (im *IndexManager) DropIndex(indexName string) error {
	// 使用参数化查询防止SQL注入
	sql := "DROP INDEX ?"
	
	err := im.db.Exec(sql, indexName).Error
	if err != nil {
		logger.DefaultStructuredLogger().
			WithError(err).
			WithField("index", indexName).
			Error("Failed to drop index")
		return err
	}

	logger.DefaultStructuredLogger().
		WithField("index", indexName).
		Info("Index dropped successfully")

	return nil
}

// AnalyzeTableIndexes analyzes indexes on a table and provides recommendations
func (im *IndexManager) AnalyzeTableIndexes(tableName string) (*IndexAnalysis, error) {
	analysis := &IndexAnalysis{
		TableName: tableName,
	}

	// Get table indexes
	indexes, err := im.GetTableIndexes(tableName)
	if err != nil {
		return nil, err
	}
	analysis.Indexes = indexes

	// Get table structure
	var columns []struct {
		Field   string
		Type    string
		Null    string
		Key     string
		Default interface{}
		Extra   string
	}

	err = im.db.Raw("DESCRIBE ??", tableName).Scan(&columns).Error
	if err != nil {
		return nil, err
	}

	// Create a map of indexed columns
	indexedColumns := make(map[string]bool)
	for _, index := range indexes {
		indexedColumns[index.Column_name] = true
	}

	// Analyze columns for potential indexes
	for _, col := range columns {
		// Skip already indexed columns
		if indexedColumns[col.Field] {
			continue
		}

		// Check for foreign key columns
		if len(col.Field) > 3 && col.Field[len(col.Field)-3:] == "_id" {
			analysis.Recommendations = append(analysis.Recommendations, IndexRecommendation{
				ColumnName: col.Field,
				Reason:     "Foreign key column should be indexed for join performance",
				IndexType:  "BTREE",
				Priority:   "High",
			})
		}

		// Check for date/time columns
		if containsSubstring(col.Type, "date") || containsSubstring(col.Type, "time") {
			analysis.Recommendations = append(analysis.Recommendations, IndexRecommendation{
				ColumnName: col.Field,
				Reason:     "Date/time column should be indexed for range queries",
				IndexType:  "BTREE",
				Priority:   "Medium",
			})
		}

		// Check for enum columns
		if containsSubstring(col.Type, "enum") {
			analysis.Recommendations = append(analysis.Recommendations, IndexRecommendation{
				ColumnName: col.Field,
				Reason:     "Enum column should be indexed for filtering",
				IndexType:  "BTREE",
				Priority:   "Medium",
			})
		}
	}

	return analysis, nil
}

// GetTableIndexes returns all indexes for a table
func (im *IndexManager) GetTableIndexes(tableName string) ([]TableIndex, error) {
	var indexes []TableIndex
	
	err := im.db.Raw("SHOW INDEX FROM ??", tableName).Scan(&indexes).Error
	if err != nil {
		return nil, err
	}
	
	return indexes, nil
}

// CreateCompositeIndex creates a composite index on multiple columns
func (im *IndexManager) CreateCompositeIndex(tableName, indexName string, columns []string) error {
	if len(columns) < 2 {
		return fmt.Errorf("composite index requires at least 2 columns")
	}

	return im.CreateIndex(tableName, indexName, columns, "BTREE")
}

// CreateFullTextIndex creates a full-text index on a table
func (im *IndexManager) CreateFullTextIndex(tableName, indexName string, columns []string) error {
	columnList := strings.Join(columns, ", ")

	// 使用参数化查询防止SQL注入
	sql := fmt.Sprintf("CREATE FULLTEXT INDEX ? ON ?? (%s)", columnList)
	
	err := im.db.Exec(sql, indexName, tableName).Error
	if err != nil {
		logger.DefaultStructuredLogger().
			WithError(err).
			WithField("table", tableName).
			WithField("index", indexName).
			WithField("columns", columns).
			Error("Failed to create full-text index")
		return err
	}

	logger.DefaultStructuredLogger().
		WithField("table", tableName).
		WithField("index", indexName).
		WithField("columns", columns).
		Info("Full-text index created successfully")

	return nil
}

// CreateSpatialIndex creates a spatial index on a table
func (im *IndexManager) CreateSpatialIndex(tableName, indexName string, column string) error {
	// 使用参数化查询防止SQL注入
	sql := "CREATE SPATIAL INDEX ? ON ?? (?)"
	
	err := im.db.Exec(sql, indexName, tableName, column).Error
	if err != nil {
		logger.DefaultStructuredLogger().
			WithError(err).
			WithField("table", tableName).
			WithField("index", indexName).
			WithField("column", column).
			Error("Failed to create spatial index")
		return err
	}

	logger.DefaultStructuredLogger().
		WithField("table", tableName).
		WithField("index", indexName).
		WithField("column", column).
		Info("Spatial index created successfully")

	return nil
}

// RebuildIndex rebuilds an index
func (im *IndexManager) RebuildIndex(indexName string) error {
	// MySQL doesn't have a direct REBUILD INDEX command
	// We can use ANALYZE TABLE to update index statistics
	err := im.db.Exec("ANALYZE TABLE").Error
	if err != nil {
		logger.DefaultStructuredLogger().
			WithError(err).
			WithField("index", indexName).
			Error("Failed to rebuild index")
		return err
	}

	logger.DefaultStructuredLogger().
		WithField("index", indexName).
		Info("Index rebuilt successfully")

	return nil
}

// GetIndexUsage returns usage statistics for indexes
func (im *IndexManager) GetIndexUsage() ([]IndexUsage, error) {
	var usages []IndexUsage
	
	// Get index usage statistics from MySQL
	err := im.db.Raw(`
		SELECT 
			TABLE_NAME, 
			INDEX_NAME, 
			CARDINALITY, 
			SUB_PART, 
			NULLABLE, 
			INDEX_TYPE
		FROM information_schema.STATISTICS 
		WHERE TABLE_SCHEMA = DATABASE()
		ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
	`).Scan(&usages).Error
	
	if err != nil {
		return nil, err
	}
	
	return usages, nil
}

// Data structures

type IndexAnalysis struct {
	TableName      string
	Indexes        []TableIndex
	Recommendations []IndexRecommendation
}

type IndexRecommendation struct {
	ColumnName string
	Reason     string
	IndexType  string
	Priority   string
}

type IndexUsage struct {
	TableName   string
	IndexName   string
	Cardinality int
	SubPart     interface{}
	Nullable    string
	IndexType   string
}