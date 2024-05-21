package model

type SystemErrCode int

const (
	SystemErrCodeSuccess SystemErrCode = 0
	SystemErrCodeFailure SystemErrCode = 90000
)

type Response[T any] struct {
	Code SystemErrCode `json:"code"`
	Msg  string        `json:"msg"`
	Data T             `json:"data"`
}
