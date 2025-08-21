package jsonx

import (
	"encoding/json"
	"fmt"
	"time"
)

func MustUnmarshal[T any](data string) *T {
	var r T
	err := json.Unmarshal([]byte(data), &r)
	if err != nil {
		panic(err)
	}

	return &r
}

func MustMarshal(data any) string {
	buf, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(buf)
}

// CustomTime 自定义时间类型，支持通过 struct tag 指定格式
type CustomTime struct {
	time.Time
	format string
}

// NewCustomTime 创建一个新的 CustomTime 实例
func NewCustomTime(t time.Time) CustomTime {
	return CustomTime{
		Time:   t,
		format: time.DateTime,
	}
}

// NewCustomTimeWithFormat 创建一个指定格式的 CustomTime 实例
func NewCustomTimeWithFormat(t time.Time, format string) CustomTime {
	return CustomTime{
		Time:   t,
		format: format,
	}
}

// MarshalJSON 实现 json.Marshaler 接口
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	format := ct.format
	if format == "" {
		format = time.RFC3339 // 默认格式
	}
	return []byte(ct.Time.Format(format)), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// 去掉引号
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// 尝试不同的格式解析时间
	formats := []string{
		ct.format,             // 指定的格式
		time.DateTime,         // 时间戳格式
		time.RFC3339,          // RFC3339 标准格式
		"2006-01-02T15:04:05", // ISO8601 格式
		"2006-01-02",          // 日期格式
		time.RFC822,           // RFC822 格式
		time.RFC850,           // RFC850 格式
		time.Kitchen,          // Kitchen 格式
	}

	var err error
	for _, format := range formats {
		if format == "" {
			continue
		}
		ct.Time, err = time.Parse(format, str)
		if err == nil {
			break
		}
	}

	// 如果所有格式都失败，返回错误
	if err != nil {
		return fmt.Errorf("无法解析时间字符串: %s", str)
	}

	return nil
}

// String 实现 Stringer 接口
func (ct CustomTime) String() string {
	if ct.format == "" {
		ct.format = time.RFC3339 // 默认格式
	}
	return ct.Time.Format(ct.format)
}

// SetFormat 设置时间格式
func (ct *CustomTime) SetFormat(format string) {
	ct.format = format
}

// GetFormat 获取时间格式
func (ct *CustomTime) GetFormat() string {
	return ct.format
}
