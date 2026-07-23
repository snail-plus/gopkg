package log

import (
	"sync"
	"time"

	"github.com/DeRuina/timberjack"
)

// 自定义按天分割的日志写入器
type dailyRotator struct {
	mu sync.Mutex         // 保证线程安全
	lj *timberjack.Logger // 底层timberjack写入器
}

// NewDailyRotator 创建按天分割的日志写入器
func NewDailyRotator(logPath string, maxSize, maxBackups, maxAge int) *dailyRotator {
	lj := &timberjack.Logger{
		Filename:         logPath,
		MaxSize:          maxSize,
		MaxBackups:       maxBackups,
		MaxAge:           maxAge,
		Compress:         true,
		LocalTime:        true,
		RotationInterval: 24 * time.Hour, // 按天轮转
	}

	return &dailyRotator{
		lj: lj,
	}
}

// Write 实现zapcore.WriteSyncer的Write方法
func (d *dailyRotator) Write(p []byte) (n int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.lj.Write(p)
}

// Sync 实现zapcore.WriteSyncer的Sync方法（确保日志刷新到磁盘）
func (d *dailyRotator) Sync() error {
	return nil
}
