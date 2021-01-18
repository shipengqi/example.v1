package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shipengqi/example.v1/blog/pkg/e"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	e := e.Cause(err)

	c.JSON(http.StatusOK, Response{
		Code: e.Code(),
		Msg:  e.Message(),
		Data: data,
	})

	return
}
