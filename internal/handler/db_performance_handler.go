package handler

import (
	"net/http"

	"go-admin/internal/database"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// DBPerformanceHandler handles database performance-related HTTP requests
type DBPerformanceHandler struct {
	queryOptimizer *database.QueryOptimizer
	indexManager   *database.IndexManager
}

// NewDBPerformanceHandler creates a new database performance handler
func NewDBPerformanceHandler() *DBPerformanceHandler {
	return &DBPerformanceHandler{
		queryOptimizer: database.NewQueryOptimizer(),
		indexManager:   database.NewIndexManager(),
	}
}

// GetQueryStats handles requests to get database query statistics
func (h *DBPerformanceHandler) GetQueryStats(c *gin.Context) {
	stats, err := h.queryOptimizer.GetQueryStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get query statistics")
		return
	}

	response.Success(c, "Query statistics retrieved successfully", stats)
}

// GetSlowQueries handles requests to get slow queries analysis
func (h *DBPerformanceHandler) GetSlowQueries(c *gin.Context) {
	analyses, err := h.queryOptimizer.AnalyzeSlowQueries()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to analyze slow queries")
		return
	}

	response.Success(c, "Slow queries analysis retrieved successfully", analyses)
}

// ExplainQuery handles requests to explain a query
func (h *DBPerformanceHandler) ExplainQuery(c *gin.Context) {
	var request struct {
		Query string   `json:"query" binding:"required"`
		Args  []interface{} `json:"args"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	results, err := h.queryOptimizer.ExplainQuery(request.Query, request.Args...)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to explain query")
		return
	}

	response.Success(c, "Query explanation retrieved successfully", results)
}

// GetTableIndexes handles requests to get table indexes
func (h *DBPerformanceHandler) GetTableIndexes(c *gin.Context) {
	tableName := c.Param("table")
	if tableName == "" {
		response.Error(c, http.StatusBadRequest, "Table name is required")
		return
	}

	indexes, err := h.indexManager.GetTableIndexes(tableName)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get table indexes")
		return
	}

	response.Success(c, "Table indexes retrieved successfully", indexes)
}

// AnalyzeTableIndexes handles requests to analyze table indexes
func (h *DBPerformanceHandler) AnalyzeTableIndexes(c *gin.Context) {
	tableName := c.Param("table")
	if tableName == "" {
		response.Error(c, http.StatusBadRequest, "Table name is required")
		return
	}

	analysis, err := h.indexManager.AnalyzeTableIndexes(tableName)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to analyze table indexes")
		return
	}

	response.Success(c, "Table indexes analysis retrieved successfully", analysis)
}

// CreateIndex handles requests to create an index
func (h *DBPerformanceHandler) CreateIndex(c *gin.Context) {
	var request struct {
		TableName string   `json:"table_name" binding:"required"`
		IndexName string   `json:"index_name" binding:"required"`
		Columns   []string `json:"columns" binding:"required,min=1"`
		IndexType string   `json:"index_type"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	err := h.indexManager.CreateIndex(request.TableName, request.IndexName, request.Columns, request.IndexType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create index")
		return
	}

	response.Success(c, "Index created successfully", nil)
}

// DropIndex handles requests to drop an index
func (h *DBPerformanceHandler) DropIndex(c *gin.Context) {
	indexName := c.Param("index")
	if indexName == "" {
		response.Error(c, http.StatusBadRequest, "Index name is required")
		return
	}

	err := h.indexManager.DropIndex(indexName)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to drop index")
		return
	}

	response.Success(c, "Index dropped successfully", nil)
}

// OptimizeTable handles requests to optimize a table
func (h *DBPerformanceHandler) OptimizeTable(c *gin.Context) {
	tableName := c.Param("table")
	if tableName == "" {
		response.Error(c, http.StatusBadRequest, "Table name is required")
		return
	}

	err := h.queryOptimizer.OptimizeTable(tableName)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to optimize table")
		return
	}

	response.Success(c, "Table optimized successfully", nil)
}

// GetIndexUsage handles requests to get index usage statistics
func (h *DBPerformanceHandler) GetIndexUsage(c *gin.Context) {
	usages, err := h.indexManager.GetIndexUsage()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get index usage")
		return
	}

	response.Success(c, "Index usage retrieved successfully", usages)
}

// SuggestMissingIndexes handles requests to suggest missing indexes
func (h *DBPerformanceHandler) SuggestMissingIndexes(c *gin.Context) {
	tableName := c.Param("table")
	if tableName == "" {
		response.Error(c, http.StatusBadRequest, "Table name is required")
		return
	}

	suggestions, err := h.queryOptimizer.SuggestMissingIndexes(tableName)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to suggest missing indexes")
		return
	}

	response.Success(c, "Missing index suggestions retrieved successfully", suggestions)
}

// CreateCompositeIndex handles requests to create a composite index
func (h *DBPerformanceHandler) CreateCompositeIndex(c *gin.Context) {
	var request struct {
		TableName string   `json:"table_name" binding:"required"`
		IndexName string   `json:"index_name" binding:"required"`
		Columns   []string `json:"columns" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	err := h.indexManager.CreateCompositeIndex(request.TableName, request.IndexName, request.Columns)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create composite index")
		return
	}

	response.Success(c, "Composite index created successfully", nil)
}

// CreateFullTextIndex handles requests to create a full-text index
func (h *DBPerformanceHandler) CreateFullTextIndex(c *gin.Context) {
	var request struct {
		TableName string   `json:"table_name" binding:"required"`
		IndexName string   `json:"index_name" binding:"required"`
		Columns   []string `json:"columns" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format")
		return
	}

	err := h.indexManager.CreateFullTextIndex(request.TableName, request.IndexName, request.Columns)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create full-text index")
		return
	}

	response.Success(c, "Full-text index created successfully", nil)
}

// RebuildIndex handles requests to rebuild an index
func (h *DBPerformanceHandler) RebuildIndex(c *gin.Context) {
	indexName := c.Param("index")
	if indexName == "" {
		response.Error(c, http.StatusBadRequest, "Index name is required")
		return
	}

	err := h.indexManager.RebuildIndex(indexName)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to rebuild index")
		return
	}

	response.Success(c, "Index rebuilt successfully", nil)
}