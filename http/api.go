package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	xerrors "github.com/snail-plus/gopkg/errors"
	"github.com/snail-plus/gopkg/http/model"
	"github.com/snail-plus/gopkg/log"
	"net/http"
)

type Api struct{}

func (r Api) Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, model.Response[any]{Data: data, Code: model.SystemErrCodeSuccess})
}

func (r Api) Failure(c *gin.Context, err error) {
	log.Infof("Failure: %v", err)

	// 判断错误类型是不是 SystemError
	var systemErr *xerrors.SystemError
	if errors.As(err, &systemErr) {
		code := systemErr.Code
		c.JSON(http.StatusOK, model.Response[any]{Code: model.SystemErrCode(code), Msg: systemErr.Message})
	} else {
		c.JSON(http.StatusOK, model.Response[any]{
			Code: model.SystemErrCodeFailure,
			Msg:  err.Error(),
		})
	}

}

func (r Api) Result(c *gin.Context, data any, err error) {
	if err != nil {
		r.Failure(c, err)
	} else {
		r.Success(c, data)
	}
}
