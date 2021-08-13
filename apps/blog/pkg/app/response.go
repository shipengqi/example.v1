package app

import (
	"github.com/shipengqi/example.v1/apps/blog/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	errno := e.Cause(err)

	c.JSON(http.StatusOK, Response{
		Code: errno.Code(),
		Msg:  errno.Message(),
		Data: data,
	})

	return
}
