// Copyright 2024 eve.  All rights reserved.

// Package log is a log package used by onex project.
//
//nolint:interfacebloat
package log

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

type Field = zapcore.Field

// Logger 定义了 onex 项目的日志接口. 该接口只包含了支持的日志记录方法.
type Logger interface {
	Debugf(format string, args ...any)
	Debugw(msg string, keyvals ...any)
	Infof(format string, args ...any)
	Infow(msg string, keyvals ...any)
	Warnf(format string, args ...any)
	Warnw(msg string, keyvals ...any)
	Errorf(format string, args ...any)
	Errorw(err error, msg string, keyvals ...any)
	Panicf(format string, args ...any)
	Panicw(msg string, keyvals ...any)
	Fatalf(format string, args ...any)
	Fatalw(msg string, keyvals ...any)
	With(fields ...Field) Logger
	AddCallerSkip(skip int) Logger
	Sync()

	gormlogger.Interface
}

// zapLogger 是 Logger 接口的具体实现. 它底层封装了 zap.Logger.
type zapLogger struct {
	z         *zap.Logger
	s         *zap.SugaredLogger
	sugarOnce sync.Once
	opts      *Options
}

// sugar 返回缓存的 SugaredLogger，避免每次调用都重新创建.
func (l *zapLogger) sugar() *zap.SugaredLogger {
	l.sugarOnce.Do(func() { l.s = l.z.Sugar() })
	return l.s
}

// 确保 zapLogger 实现了 Logger 接口. 以下变量赋值，可以使错误在编译期被发现.
var _ Logger = (*zapLogger)(nil)

var (
	// std 定义了默认的全局 Logger.
	std = newLogger(NewOptions())
)

var once sync.Once

// Init 使用指定的选项初始化 Logger.
func Init(opts *Options) {
	once.Do(func() {
		std = newLogger(opts)
		Infof("init logger done")
	})
}

// newLogger 根据传入的 opts 创建 Logger.
func newLogger(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	// 将文本格式的日志级别，例如 info 转换为 zapcore.Level 类型以供后面使用
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		// 如果指定了非法的日志级别，则默认使用 info 级别
		zapLevel = zapcore.InfoLevel
	}

	// 创建一个默认的 encoder 配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 自定义 MessageKey 为 message，message 语义更明确
	encoderConfig.MessageKey = "message"
	// 自定义 TimeKey 为 timestamp，timestamp 语义更明确
	encoderConfig.TimeKey = "timestamp"
	// 指定时间序列化函数，将时间序列化为 `2006-01-02 15:04:05.000` 格式，更易读
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// 指定 time.Duration 序列化函数，将 time.Duration 序列化为经过的毫秒数的浮点数
	// 毫秒数比默认的秒数更精确
	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}
	// when output to local path, with color is forbidden
	if opts.Format == "console" && opts.EnableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 使用 cfg 创建 *zap.Logger 对象
	var err error
	var zapLog *zap.Logger

	if len(opts.OutputPaths) > 0 && opts.OutputPaths[0] != "stdout" {
		logFileName := opts.OutputPaths[0]
		// 保持历史行为：日志统一写入 ./logs 目录，文件名取自配置路径的最后一段（去掉扩展名），
		// 避免 `logs/app.log` 这类路径被拼接成 `./logs/logs/app.log`
		base := filepath.Base(logFileName)
		logPathTemplate := filepath.Join("./logs", strings.TrimSuffix(base, filepath.Ext(base))+".log")

		// 确保日志目录存在
		_ = os.MkdirAll("./logs", 0755)

		// 创建按天分割的写入器
		dailyWriter := NewDailyRotator(logPathTemplate, opts.MaxSize, opts.MaxBackups, opts.MaxAge)

		zapLog = zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(dailyWriter),
			zap.NewAtomicLevelAt(zapLevel)))
	} else {
		// 创建构建 zap.Logger 需要的配置
		cfg := &zap.Config{
			// 是否在日志中显示调用日志所在的文件和行号，例如：`"caller":"onex/onex.go:75"`
			DisableCaller: opts.DisableCaller,
			// 是否禁止在 panic 及以上级别打印堆栈信息
			DisableStacktrace: opts.DisableStacktrace,
			// 指定日志级别
			Level: zap.NewAtomicLevelAt(zapLevel),
			// 指定日志显示格式，可选值：console, json
			Encoding:      opts.Format,
			EncoderConfig: encoderConfig,
			// 指定日志输出位置
			OutputPaths: []string{"stdout"},
			// 设置 zap 内部错误输出位置
			ErrorOutputPaths: []string{"stderr"},
		}
		zapLog, err = cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(2))
	}

	if err != nil {
		panic(err)
	}

	logger := &zapLogger{z: zapLog, opts: opts}
	// 把标准库的 log.Logger 的 info 级别的输出重定向到 zap.Logger
	zap.RedirectStdLog(zapLog)

	return logger
}

func Default() Logger {
	return std
}

func (l *zapLogger) Options() *Options {
	return l.opts
}

// Sync 调用底层 zap.Logger 的 Sync 方法，将缓存中的日志刷新到磁盘文件中. 主程序需要在退出前调用 Sync.
func Sync() { std.Sync() }

func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}

// Debugf 输出 debug 级别的日志.
func Debugf(format string, args ...any) {
	std.Debugf(format, args...)
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.sugar().Debugf(format, args...)
}

// Debugw 输出 debug 级别的日志.
func Debugw(msg string, keyvals ...any) {
	std.Debugw(msg, keyvals...)
}

func (l *zapLogger) Debugw(msg string, keyvals ...any) {
	l.sugar().Debugw(msg, keyvals...)
}

// Infof 输出 info 级别的日志.
func Infof(format string, args ...any) {
	std.Infof(format, args...)
}

// ErrorS 以结构化方式记录错误。err 为 nil 时按普通消息记录，不会 panic。
// keyvals 为 klog 风格的 key/value 对。
func ErrorS(err error, msg string, keyvals ...any) {
	fields := fieldsFromKeyvals(keyvals...)
	if err != nil {
		fields = append([]Field{zap.Error(err)}, fields...)
	}
	std.z.Error(msg, fields...)
}

// InfoS 以结构化方式记录消息，keyvals 为 klog 风格的 key/value 对。
func InfoS(msg string, keyvals ...any) {
	std.z.Info(msg, fieldsFromKeyvals(keyvals...)...)
}

func (l *zapLogger) Infof(msg string, keyvals ...any) {
	l.sugar().Infof(msg, keyvals...)
}

// Infow 输出 info 级别的日志.
func Infow(msg string, keyvals ...any) {
	std.Infow(msg, keyvals...)
}

func (l *zapLogger) Infow(msg string, keyvals ...any) {
	l.sugar().Infow(msg, keyvals...)
}

// Warnf 输出 warning 级别的日志.
func Warnf(format string, args ...any) {
	std.Warnf(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...any) {
	l.sugar().Warnf(format, args...)
}

// Warnw 输出 warning 级别的日志.
func Warnw(msg string, keyvals ...any) {
	std.Warnw(msg, keyvals...)
}

func (l *zapLogger) Warnw(msg string, keyvals ...any) {
	l.sugar().Warnw(msg, keyvals...)
}

// Errorf 输出 error 级别的日志.
func Errorf(format string, args ...any) {
	std.Errorf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...any) {
	l.sugar().Errorf(format, args...)
}

// Errorw 输出 error 级别的日志.
func Errorw(err error, msg string, keyvals ...any) {
	std.Errorw(err, msg, keyvals...)
}

func (l *zapLogger) Errorw(err error, msg string, keyvals ...any) {
	l.sugar().Errorw(msg, append(keyvals, "err", err)...)
}

// Panicf 输出 panic 级别的日志.
func Panicf(format string, args ...any) {
	std.Panicf(format, args...)
}

func (l *zapLogger) Panicf(format string, args ...any) {
	l.sugar().Panicf(format, args...)
}

// Panicw 输出 panic 级别的日志.
func Panicw(msg string, keyvals ...any) {
	std.Panicw(msg, keyvals...)
}

func (l *zapLogger) Panicw(msg string, keyvals ...any) {
	l.sugar().Panicw(msg, keyvals...)
}

// Fatalf 输出 fatal 级别的日志.
func Fatalf(format string, args ...any) {
	std.Fatalf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...any) {
	l.sugar().Fatalf(format, args...)
}

// Fatalw 输出 fatal 级别的日志.
func Fatalw(msg string, keyvals ...any) {
	std.Fatalw(msg, keyvals...)
}

func (l *zapLogger) Fatalw(msg string, keyvals ...any) {
	l.sugar().Fatalw(msg, keyvals...)
}

func With(fields ...Field) Logger {
	return std.With(fields...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *zapLogger) With(fields ...Field) Logger {
	if len(fields) == 0 {
		return l
	}

	lc := l.clone()
	lc.z = lc.z.With(fields...)
	return lc
}

func AddCallerSkip(skip int) Logger {
	return std.AddCallerSkip(skip)
}

// AddCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the AddCaller option). When building wrappers around the
// Logger and SugaredLogger, supplying this Option prevents zap from always
// reporting the wrapper code as the caller.
func (l *zapLogger) AddCallerSkip(skip int) Logger {
	lc := l.clone()
	lc.z = lc.z.WithOptions(zap.AddCallerSkip(skip))
	return lc
}

// clone 深度拷贝 zapLogger.
// 子 logger 需要基于新的 z 重新构建 SugaredLogger，否则 With/AddCallerSkip 等变更不会生效.
// 注意逐字段赋值，避免拷贝内嵌的 sync.Once（vet 会报 copylocks）.
func (l *zapLogger) clone() *zapLogger {
	return &zapLogger{
		z:    l.z,
		opts: l.opts,
	}
}
