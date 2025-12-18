package logger

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// AsyncLogger 异步日志处理器
type AsyncLogger struct {
	logger    *zap.Logger
	logCh     chan LogEntry
	wg        sync.WaitGroup
	closeOnce sync.Once
	closed    chan struct{}
	config    AsyncLogConfig
}

// LogEntry 日志条目结构
type LogEntry struct {
	Level   zapcore.Level
	Message string
	Fields  []zap.Field
	Time    time.Time
}

// AsyncLogConfig 异步日志配置
type AsyncLogConfig struct {
	Enabled            bool          `mapstructure:"async_enabled"`
	BufferSize         int           `mapstructure:"buffer_size"`
	FlushInterval      time.Duration `mapstructure:"flush_interval"`
	SamplingFirst      int64         `mapstructure:"sampling_first"`
	SamplingThereAfter int64         `mapstructure:"sampling_there_after"`
}

// DefaultAsyncLogConfig 默认异步日志配置
func DefaultAsyncLogConfig() AsyncLogConfig {
	return AsyncLogConfig{
		Enabled:            true,
		BufferSize:         1000,
		FlushInterval:      time.Second,
		SamplingFirst:      100,
		SamplingThereAfter: 1000,
	}
}

// NewAsyncLogger 创建新的异步日志处理器
func NewAsyncLogger(logger *zap.Logger, config AsyncLogConfig) *AsyncLogger {
	if !config.Enabled {
		return &AsyncLogger{logger: logger}
	}

	asyncLogger := &AsyncLogger{
		logger: logger,
		logCh:  make(chan LogEntry, config.BufferSize),
		closed: make(chan struct{}),
		config: config,
	}

	// 启动后台处理协程
	asyncLogger.wg.Add(1)
	go asyncLogger.processLogs()

	// 启动定期刷新协程
	if config.FlushInterval > 0 {
		asyncLogger.wg.Add(1)
		go asyncLogger.flushPeriodically()
	}

	return asyncLogger
}

// processLogs 处理日志条目
func (al *AsyncLogger) processLogs() {
	defer al.wg.Done()

	sampler := NewSampler(al.config.SamplingFirst, al.config.SamplingThereAfter)

	for {
		select {
		case entry, ok := <-al.logCh:
			if !ok {
				return
			}

			// 应用采样
			if sampler.Sample(entry.Message) {
				al.writeLog(entry)
			}
		case <-al.closed:
			// 处理剩余的日志条目
			for {
				select {
				case entry, ok := <-al.logCh:
					if !ok {
						return
					}
					al.writeLog(entry)
				default:
					return
				}
			}
		}
	}
}

// flushPeriodically 定期刷新日志
func (al *AsyncLogger) flushPeriodically() {
	defer al.wg.Done()

	ticker := time.NewTicker(al.config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if al.logger != nil {
				al.logger.Sync()
			}
		case <-al.closed:
			return
		}
	}
}

// writeLog 写入日志
func (al *AsyncLogger) writeLog(entry LogEntry) {
	if al.logger == nil {
		return
	}

	switch entry.Level {
	case zapcore.DebugLevel:
		al.logger.Debug(entry.Message, entry.Fields...)
	case zapcore.InfoLevel:
		al.logger.Info(entry.Message, entry.Fields...)
	case zapcore.WarnLevel:
		al.logger.Warn(entry.Message, entry.Fields...)
	case zapcore.ErrorLevel:
		al.logger.Error(entry.Message, entry.Fields...)
	case zapcore.FatalLevel:
		al.logger.Fatal(entry.Message, entry.Fields...)
	}
}

// Log 记录日志
func (al *AsyncLogger) Log(level zapcore.Level, msg string, fields ...zap.Field) {
	if !al.config.Enabled {
		// 如果异步日志未启用，直接写入
		switch level {
		case zapcore.DebugLevel:
			al.logger.Debug(msg, fields...)
		case zapcore.InfoLevel:
			al.logger.Info(msg, fields...)
		case zapcore.WarnLevel:
			al.logger.Warn(msg, fields...)
		case zapcore.ErrorLevel:
			al.logger.Error(msg, fields...)
		case zapcore.FatalLevel:
			al.logger.Fatal(msg, fields...)
		}
		return
	}

	entry := LogEntry{
		Level:   level,
		Message: msg,
		Fields:  fields,
		Time:    time.Now(),
	}

	select {
	case al.logCh <- entry:
	default:
		// 如果通道已满，降级为同步写入并发出警告
		al.writeLog(entry)
		if al.logger != nil {
			al.logger.Warn("Async log buffer is full, falling back to sync logging")
		}
	}
}

// Close 关闭异步日志处理器
func (al *AsyncLogger) Close() error {
	al.closeOnce.Do(func() {
		close(al.closed)
		if al.config.Enabled {
			close(al.logCh)
		}
		al.wg.Wait()
	})
	return nil
}

// Sampler 日志采样器
type Sampler struct {
	first      int64
	thereAfter int64
	counter    map[string]*samplingCounter
	mutex      sync.Mutex
}

type samplingCounter struct {
	count int64
}

// NewSampler 创建新的日志采样器
func NewSampler(first, thereAfter int64) *Sampler {
	return &Sampler{
		first:      first,
		thereAfter: thereAfter,
		counter:    make(map[string]*samplingCounter),
	}
}

// Sample 判断是否应该采样记录该消息
func (s *Sampler) Sample(msg string) bool {
	// 如果没有设置采样参数，则记录所有消息
	if s.first <= 0 && s.thereAfter <= 0 {
		return true
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	counter, exists := s.counter[msg]
	if !exists {
		counter = &samplingCounter{}
		s.counter[msg] = counter
	}

	counter.count++

	// 前first条全部记录
	if counter.count <= s.first {
		return true
	}

	// 之后每隔thereAfter条记录一次
	if s.thereAfter > 0 && (counter.count-s.first)%s.thereAfter == 0 {
		return true
	}

	return false
}
