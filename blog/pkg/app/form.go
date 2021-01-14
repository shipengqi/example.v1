package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/shipengqi/example.v1/blog/pkg/errno"
)

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, form interface{}) error {
	err := c.Bind(form)
	if err != nil {
		return errno.ErrBadRequest
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return errno.ErrInternalServer
	}
	if !check {
		MarkErrors(valid.Errors)
		return errno.ErrBadRequest
	}

	return nil
}
