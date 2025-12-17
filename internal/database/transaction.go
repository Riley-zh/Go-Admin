package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-admin/internal/logger"

	"gorm.io/gorm"
)

// TransactionManager 管理数据库事务
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager 创建新的事务管理器
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

// TransactionOptions 事务选项
type TransactionOptions struct {
	Timeout         time.Duration
	IsolationLevel  sql.IsolationLevel
	ReadOnly        bool
	RetryAttempts   int
	RetryDelay      time.Duration
}

// DefaultTransactionOptions 返回默认的事务选项
func DefaultTransactionOptions() TransactionOptions {
	return TransactionOptions{
		Timeout:        30 * time.Second,
		IsolationLevel: sql.LevelReadCommitted,
		ReadOnly:       false,
		RetryAttempts:  3,
		RetryDelay:     100 * time.Millisecond,
	}
}

// TransactionResult 事务执行结果
type TransactionResult struct {
	Success     bool
	Error       error
	RollbackErr error
	Retries     int
	Duration    time.Duration
}

// WithTransaction 在事务中执行函数
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(*gorm.DB) error, opts ...TransactionOptions) (*TransactionResult, error) {
	options := DefaultTransactionOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	startTime := time.Now()
	result := &TransactionResult{
		Success:  false,
		Retries:  0,
		Duration: 0,
	}

	// 设置超时上下文
	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel()
	}

	// 重试逻辑
	var lastErr error
	for i := 0; i <= options.RetryAttempts; i++ {
		if i > 0 {
			result.Retries = i
			// 等待一段时间再重试
			time.Sleep(options.RetryDelay * time.Duration(i))
			logger.DefaultStructuredLogger().
				WithField("attempt", i).
				WithField("error", lastErr.Error()).
				Warn("Retrying database transaction")
		}

		// 开始事务
		tx := tm.db.WithContext(ctx)
		if options.IsolationLevel != 0 {
			// GORM的Begin方法不支持直接传递TxOptions，需要使用Begin()然后设置隔离级别
			tx = tx.Begin()
			tx.Exec(fmt.Sprintf("SET TRANSACTION ISOLATION LEVEL %s", getIsolationLevelString(options.IsolationLevel)))
		} else {
			tx = tx.Begin()
		}

		// 检查事务是否成功开始
		if tx.Error != nil {
			lastErr = tx.Error
			logger.DefaultStructuredLogger().
				WithError(tx.Error).
				Error("Failed to begin database transaction")
			continue
		}

		// 执行事务函数
		err := fn(tx)
		if err != nil {
			// 发生错误，回滚事务
			rollbackErr := tx.Rollback().Error
			result.RollbackErr = rollbackErr
			lastErr = err

			logger.DefaultStructuredLogger().
				WithError(err).
				WithField("rollback_error", rollbackErr).
				Error("Database transaction failed and was rolled back")

			// 如果是死锁错误，可以重试
			if isDeadlockError(err) && i < options.RetryAttempts {
				continue
			}
			break
		}

		// 提交事务
		if commitErr := tx.Commit().Error; commitErr != nil {
			result.RollbackErr = commitErr
			lastErr = commitErr

			logger.DefaultStructuredLogger().
				WithError(commitErr).
				Error("Failed to commit database transaction")
			break
		}

		// 事务成功
		result.Success = true
		result.Duration = time.Since(startTime)

		logger.DefaultStructuredLogger().
			WithField("duration_ms", result.Duration.Milliseconds()).
			WithField("retries", result.Retries).
			Info("Database transaction completed successfully")
		return result, nil
	}

	// 所有重试都失败
	result.Error = lastErr
	result.Duration = time.Since(startTime)

	logger.DefaultStructuredLogger().
		WithError(lastErr).
		WithField("retries", result.Retries).
		WithField("duration_ms", result.Duration.Milliseconds()).
		Error("Database transaction failed after all retries")

	return result, lastErr
}

// WithReadOnlyTransaction 在只读事务中执行函数
func (tm *TransactionManager) WithReadOnlyTransaction(ctx context.Context, fn func(*gorm.DB) error, opts ...TransactionOptions) (*TransactionResult, error) {
	options := DefaultTransactionOptions()
	if len(opts) > 0 {
		options = opts[0]
	}
	options.ReadOnly = true

	return tm.WithTransaction(ctx, fn, options)
}

// WithNestedTransaction 在嵌套事务中执行函数
func (tm *TransactionManager) WithNestedTransaction(ctx context.Context, parentTx *gorm.DB, fn func(*gorm.DB) error, opts ...TransactionOptions) (*TransactionResult, error) {
	options := DefaultTransactionOptions()
	if len(opts) > 0 {
		options = opts[0]
	}

	startTime := time.Now()
	result := &TransactionResult{
		Success:  false,
		Retries:  0,
		Duration: 0,
	}

	// 设置超时上下文
	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel()
	}

	// 使用保存点实现嵌套事务
	savepointName := fmt.Sprintf("sp_%d", time.Now().UnixNano())

	// 创建保存点
	if err := parentTx.Exec(fmt.Sprintf("SAVEPOINT %s", savepointName)).Error; err != nil {
		result.Error = err
		result.Duration = time.Since(startTime)

		logger.DefaultStructuredLogger().
			WithError(err).
			Error("Failed to create savepoint for nested transaction")
		return result, err
	}

	// 执行函数
	err := fn(parentTx)
	if err != nil {
		// 回滚到保存点
		rollbackErr := parentTx.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", savepointName)).Error
		result.RollbackErr = rollbackErr
		result.Error = err
		result.Duration = time.Since(startTime)

		logger.DefaultStructuredLogger().
			WithError(err).
			WithField("rollback_error", rollbackErr).
			Error("Nested transaction failed and was rolled back to savepoint")
		return result, err
	}

	// 提交嵌套事务（释放保存点）
	if commitErr := parentTx.Exec(fmt.Sprintf("RELEASE SAVEPOINT %s", savepointName)).Error; commitErr != nil {
		result.Error = commitErr
		result.Duration = time.Since(startTime)

		logger.DefaultStructuredLogger().
			WithError(commitErr).
			Error("Failed to release savepoint for nested transaction")
		return result, commitErr
	}

	// 嵌套事务成功
	result.Success = true
	result.Duration = time.Since(startTime)

	logger.DefaultStructuredLogger().
		WithField("duration_ms", result.Duration.Milliseconds()).
		WithField("savepoint", savepointName).
		Info("Nested transaction completed successfully")
	return result, nil
}

// isDeadlockError 检查是否是死锁错误
func isDeadlockError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	// MySQL deadlock error
	if strings.Contains(errStr, "Deadlock") {
		return true
	}
	// PostgreSQL deadlock error
	if strings.Contains(errStr, "deadlock detected") {
		return true
	}
	// SQL Server deadlock error
	if strings.Contains(errStr, "deadlock") && strings.Contains(errStr, "victim") {
		return true
	}

	return false
}

// getIsolationLevelString 将隔离级别转换为字符串
func getIsolationLevelString(level sql.IsolationLevel) string {
	switch level {
	case sql.LevelReadUncommitted:
		return "READ UNCOMMITTED"
	case sql.LevelReadCommitted:
		return "READ COMMITTED"
	case sql.LevelRepeatableRead:
		return "REPEATABLE READ"
	case sql.LevelSerializable:
		return "SERIALIZABLE"
	default:
		return "READ COMMITTED"
	}
}