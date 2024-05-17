package http

import (
	"gitee.com/eve_3/gopkg/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response[T any] struct {
	Code SystemErrCode `json:"code"`
	Msg  string        `json:"msg"`
	Data T             `json:"data"`
}

type Api struct{}

func (r Api) Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response[any]{Data: data, Code: SystemErrCodeSuccess})
}

func (r Api) Failure(c *gin.Context, err error) {
	log.Infof("Failure: %v", err)
	c.JSON(http.StatusOK, Response[any]{
		Code: SystemErrCodeFailure,
		Msg:  err.Error(),
	})
}

func (r Api) Result(c *gin.Context, data any, err error) {
	if err != nil {
		r.Failure(c, err)
	} else {
		r.Success(c, data)
	}
}

type SystemErrCode int

const (
	SystemErrCodeSuccess SystemErrCode = 0
	SystemErrCodeFailure SystemErrCode = 90000
)
