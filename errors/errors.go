// Copyright 2024 eve.  All rights reserved.

package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

type SystemError struct {
	Message string
	Code    int
	err     error
}

// Error 方法返回错误信息的字符串表示，满足 error 接口
func (e *SystemError) Error() string {
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

// NewSystemError 是一个创建并返回 *MyCustomError 的函数
func NewSystemError(code int, message string) *SystemError {
	cause := errors.New(message)
	err := errors.WithStack(cause)
	return &SystemError{
		Message: message,
		Code:    code,
		err:     err,
	}
}
