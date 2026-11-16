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
	if err == nil {
		r.Success(c, nil)
		return
	}

	log.Errorf("Failure: %v", err)

	// 业务错误使用显式定义的状态码与消息，可安全返回给客户端。
	var systemErr *xerrors.SystemError
	if errors.As(err, &systemErr) {
		c.JSON(http.StatusOK, model.Response[any]{
			Code: model.SystemErrCode(systemErr.Code),
			Msg:  systemErr.Message,
		})
		return
	}

	// 未知错误不暴露内部细节，并使用 500 状态码。
	c.JSON(http.StatusInternalServerError, model.Response[any]{
		Code: model.SystemErrCodeFailure,
		Msg:  "internal server error",
	})
}

func (r Api) Result(c *gin.Context, data any, err error) {
	if err != nil {
		r.Failure(c, err)
	} else {
		r.Success(c, data)
	}
}
