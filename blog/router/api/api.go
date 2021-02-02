package api

import (
	"github.com/gin-gonic/gin"
	"github.com/shipengqi/example.v1/blog/pkg/e"

	"github.com/shipengqi/example.v1/blog/pkg/app"
	apiv1 "github.com/shipengqi/example.v1/blog/router/api/v1"
	"github.com/shipengqi/example.v1/blog/service"
)

type LoginForm struct {
	Username string `form:"username" valid:"Required;"`
	Password string `form:"password" valid:"Required;"`
}

var svc *service.Service

func Init(s *service.Service) {
	svc = s
	apiv1.Init(s)
}

// @Summary Login
// @Produce application/json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /login [post]
func Login(c *gin.Context) {
	var form LoginForm
	data := make(map[string]interface{})
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.SendResponse(c, err, nil)
		return
	}

	token, rbac, err := svc.AuthSvc.Login(form.Username, form.Password)
	if err != nil {
		app.SendResponse(c, err, nil)
		return
	}
	data["token"] = token
	data["groups"] = rbac.Groups
	data["roles"] = rbac.Roles
	c.SetCookie("X-AUTH-TOKEN", token, 3600, "", "", true, true)
	app.SendResponse(c, e.OK, data)
}
