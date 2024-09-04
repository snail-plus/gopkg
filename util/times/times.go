package times

import (
	"fmt"
	"time"
)

func MillisecondsToTime(milliseconds int64) time.Time {
	seconds := milliseconds / 1000
	nanoseconds := (milliseconds % 1000) * int64(time.Millisecond)
	return time.Unix(seconds, nanoseconds)
}

const (
	DateTime = time.DateTime
	RFC3339  = time.RFC3339
)

var layouts = []string{
	time.DateTime,                 // RFC 3339 without microseconds
	"2006-01-02 15:04:05.999",     // RFC 3339 with milliseconds
	"2006-01-02T15:04:05",         // RFC 3339 without microseconds, with T separator
	"2006-01-02T15:04:05.999",     // RFC 3339 with milliseconds, with T separator
	"02/Jan/2006 15:04:05",        // Time.String() format
	"Mon Jan 2 15:04:05 MST 2006", // Reference time format for Go's time.Format
	time.RFC3339,                  // Time.RFC3339
	"2006-01-02T15:04Z",
}

// Parse 尝试多种常见格式解析时间字符串
func Parse(timeStr string) (time.Time, error) {
	// 定义一些常见的时间格式

	// 尝试每种格式
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, timeStr, time.Local)
		if err == nil {
			// 如果解析成功，返回时间
			return t, nil
		}
	}

	// 如果所有格式都失败，返回错误
	return time.Time{}, fmt.Errorf("无法解析时间字符串: %s", timeStr)
}

func MustParse(timeStr string) time.Time {
	t, err := Parse(timeStr)
	if err != nil {
		panic(err)
	}
	return t
}

func Format(tt time.Time, formatStr string) string {
	return tt.Format(formatStr)
}

// GetStartAndEndTimeOfDay 根据传入的时间获取这个时间的开始时间 结束时间
func GetStartAndEndTimeOfDay(t time.Time) (startTime, endTime time.Time) {
	// 获取这个时间的当天的开始时间
	startTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	// 加上一天的时间得到当天的结束时间
	endTime = startTime.AddDate(0, 0, 1).Add(-1 * time.Second)

	return startTime, endTime
}

func SubtractDays(t time.Time, days int) time.Time {
	// 计算往前推算的天数
	return t.AddDate(0, 0, -days)
}
