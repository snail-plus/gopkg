package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.evebatterycloud.com/infra/gopkg/http/model"
	"gitlab.evebatterycloud.com/infra/gopkg/log"
	"net/http"
)

type Api struct{}

func (r Api) Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, model.Response[any]{Data: data, Code: model.SystemErrCodeSuccess})
}

func (r Api) Failure(c *gin.Context, err error) {
	log.Infof("Failure: %v", err)
	c.JSON(http.StatusOK, model.Response[any]{
		Code: model.SystemErrCodeFailure,
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
