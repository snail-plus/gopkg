package log

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
	"time"
)

// 自定义按天分割的日志写入器
type dailyRotator struct {
	mu         sync.Mutex         // 保证线程安全
	lj         *lumberjack.Logger // 底层lumberjack写入器
	currentDay string             // 当前日志文件对应的日期（格式：2006-01-02）
	logPath    string             // 日志文件路径模板（如：./logs/app-%s.log）
}

// NewDailyRotator 创建按天分割的日志写入器
func NewDailyRotator(logPath string) *dailyRotator {
	// 初始化时获取当前日期
	currentDay := time.Now().Format("2006-01-02")
	// 构建初始日志文件路径
	initialFilename := fmt.Sprintf(logPath, currentDay)

	lj := &lumberjack.Logger{
		Filename:   initialFilename,
		MaxSize:    100,  // 单个文件最大大小（MB，即使未达阈值，也会按天分割）
		MaxBackups: 30,   // 保留旧文件数量
		MaxAge:     30,   // 保留旧文件天数
		Compress:   true, // 压缩旧文件
		LocalTime:  true, // 使用本地时间
	}

	return &dailyRotator{
		lj:         lj,
		currentDay: currentDay,
		logPath:    logPath,
	}
}

// Write 实现zapcore.WriteSyncer的Write方法
func (d *dailyRotator) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 检查是否跨天
	today := time.Now().Format("2006-01-02")
	if today != d.currentDay {
		// 跨天：关闭旧文件，创建新文件
		_ = d.lj.Close()
		d.currentDay = today
		d.lj.Filename = fmt.Sprintf(d.logPath, today)
		// 滚动日志
		_ = d.lj.Rotate()
	}

	// 写入日志内容
	return d.lj.Write(p)
}

// Sync 实现zapcore.WriteSyncer的Sync方法（确保日志刷新到磁盘）
func (d *dailyRotator) Sync() error {
	return nil
}
