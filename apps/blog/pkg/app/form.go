package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/shipengqi/example.v1/apps/blog/pkg/e"
)

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, form interface{}) error {
	err := c.Bind(form)
	if err != nil {
		return e.ErrBadRequest
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return e.ErrInternalServer
	}
	if !check {
		MarkErrors(valid.Errors)
		return e.ErrBadRequest
	}

	return nil
}
