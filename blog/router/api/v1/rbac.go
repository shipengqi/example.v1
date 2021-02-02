package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/shipengqi/example.v1/blog/pkg/app"
	"github.com/shipengqi/example.v1/blog/pkg/e"
)

type AddUserForm struct {
	Username string `form:"username" valid:"Required"`
	Password string `form:"password" valid:"Required"`
	Phone    string `form:"phone" valid:"Required"`
	Email    string `form:"email"`
}

func GetUsers(c *gin.Context) {

}

func AddUser(c *gin.Context) {
	var form AddUserForm

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.SendResponse(c, err, nil)
		return
	}
	err = svc.RBAC.AddUser(form.Username, form.Password, form.Phone, form.Email)
	if err != nil {
		app.SendResponse(c, err, nil)
		return
	}
	app.SendResponse(c, e.OK, nil)
}

func EditUser(c *gin.Context) {

}

func DeleteUser(c *gin.Context) {

}
