package log

import (
	"fmt"
	"sync"
	"time"

	"github.com/DeRuina/timberjack"
)

// 自定义按天分割的日志写入器
type dailyRotator struct {
	mu         sync.Mutex         // 保证线程安全
	lj         *timberjack.Logger // 底层lumberjack写入器
	currentDay string             // 当前日志文件对应的日期（格式：2006-01-02）
	logPath    string             // 日志文件路径模板（如：./logs/app-%s.log）
	nextCheck  int64
}

// NewDailyRotator 创建按天分割的日志写入器
func NewDailyRotator(logPath string, maxSize, maxBackups, maxAge int) *dailyRotator {
	now := time.Now()
	currentDay := now.Format(time.DateOnly)
	nextCheck := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour).Unix()
	lj := &timberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   true,
		LocalTime:  true,
	}

	return &dailyRotator{
		lj:         lj,
		currentDay: currentDay,
		logPath:    logPath,
		nextCheck:  nextCheck,
	}
}

// Write 实现zapcore.WriteSyncer的Write方法
func (d *dailyRotator) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 检查是否跨天
	now := time.Now()
	if now.Unix() >= d.nextCheck {
		// 跨天：关闭旧文件，创建新文件
		_ = d.lj.Close()
		// 更新当前日期
		d.currentDay = now.Format(time.DateOnly)
		d.lj.Filename = fmt.Sprintf(d.logPath, d.currentDay)
		// 滚动日志
		_ = d.lj.Rotate()

		// 更新下一次检查时间（明天0点）
		d.nextCheck = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour).Unix()
	}

	// 写入日志内容
	return d.lj.Write(p)
}

// Sync 实现zapcore.WriteSyncer的Sync方法（确保日志刷新到磁盘）
func (d *dailyRotator) Sync() error {
	return nil
}
