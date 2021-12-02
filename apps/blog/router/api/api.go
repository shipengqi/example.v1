package api

import (
	"github.com/gin-gonic/gin"
	model2 "github.com/shipengqi/example.v1/apps/blog/model"
	app2 "github.com/shipengqi/example.v1/apps/blog/pkg/app"
	"github.com/shipengqi/example.v1/apps/blog/pkg/e"
	apiv1 "github.com/shipengqi/example.v1/apps/blog/router/api/v1"
	"github.com/shipengqi/example.v1/apps/blog/service"
)

type LoginRequest struct {
	Username string `form:"username" valid:"Required;"`
	Password string `form:"password" valid:"Required;"`
}

type LoginResponse struct {
	model2.User

	Groups []model2.Group `json:"groups"`
	Roles  []model2.Role  `json:"roles"`
	Token  string         `json:"token"`
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
	var form LoginRequest
	err := app2.BindAndValid(c, &form)
	if err != nil {
		app2.SendResponse(c, err, nil)
		return
	}

	token, rbac, err := svc.AuthSvc.Login(form.Username, form.Password)
	if err != nil {
		app2.SendResponse(c, err, nil)
		return
	}

	c.SetCookie("X-AUTH-TOKEN", token, 3600, "", "", true, true)
	app2.SendResponse(c, e.OK, LoginResponse{
		User:   rbac.U,
		Groups: rbac.Groups,
		Roles:  rbac.Roles,
		Token:  token,
	})
}
