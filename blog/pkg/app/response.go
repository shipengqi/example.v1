package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shipengqi/example.v1/blog/pkg/errno"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	e := errno.Cause(err)

	c.JSON(http.StatusOK, Response{
		Code: e.Code(),
		Msg:  e.Message(),
		Data: data,
	})

	return
}
