package pkg

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 全局日志实例
type Logger struct {
	*zap.Logger
}

// 全局变量
var (
	globalLogger *Logger
	once         sync.Once
)

// 日志级别
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

// Options 日志配置选项
type Options struct {
	Level       string
	OutputPath  string
	ErrorPath   string
	Development bool
}

// DefaultOptions 返回默认配置
func DefaultOptions() Options {
	return Options{
		Level:       InfoLevel,
		OutputPath:  "stdout",
		ErrorPath:   "stderr",
		Development: false,
	}
}

// InitLogger 初始化日志系统
func InitLogger(options Options) error {
	var err error
	once.Do(func() {
		globalLogger, err = newLogger(options)
	})
	return err
}

// newLogger 创建新的日志实例
func newLogger(options Options) (*Logger, error) {
	// 设置日志级别
	level := zap.InfoLevel
	switch options.Level {
	case DebugLevel:
		level = zap.DebugLevel
	case WarnLevel:
		level = zap.WarnLevel
	case ErrorLevel:
		level = zap.ErrorLevel
	case FatalLevel:
		level = zap.FatalLevel
	}

	// 创建编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 创建编码器
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	if options.Development {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 创建写入器
	var writers []zapcore.WriteSyncer

	// 标准输出
	stdoutWriter := zapcore.Lock(os.Stdout)
	writers = append(writers, stdoutWriter)

	// 错误输出
	if options.ErrorPath != "" && options.ErrorPath != "stderr" {
		errWriter, _, err := zap.Open(options.ErrorPath)
		if err != nil {
			return nil, err
		}
		writers = append(writers, errWriter)
	}

	// 创建核心
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writers...),
		level,
	)

	// 创建日志器
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return &Logger{zapLogger}, nil
}

// GetLogger 获取全局日志实例
func GetLogger() *Logger {
	if globalLogger == nil {
		// 如果没有初始化，使用默认配置
		options := DefaultOptions()
		InitLogger(options)
	}
	return globalLogger
}

// Debug 记录调试信息
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info 记录信息
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn 记录警告信息
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error 记录错误信息
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal 记录致命错误并退出
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// With 添加字段到日志
func With(fields ...zap.Field) *Logger {
	return &Logger{GetLogger().With(fields...)}
}

// WithError 添加错误字段到日志
func WithError(err error) *Logger {
	return &Logger{GetLogger().With(zap.Error(err))}
}

// Sync 刷新日志缓冲区
func Sync() error {
	return GetLogger().Sync()
}
