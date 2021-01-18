package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/unknwon/com"

	"github.com/shipengqi/example.v1/blog/pkg/app"
	"github.com/shipengqi/example.v1/blog/pkg/e"
)

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

type EditTagForm struct {
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	ID         int    `form:"id" valid:"Required;Min(1)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Get multiple article tags
// @Produce application/json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}
	var state = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	data, err := svc.TagSvc.GetTags(maps)
	if err != nil {
		app.SendResponse(c, err, data)
		return
	}
	app.SendResponse(c, e.OK, data)
}

// @Summary Add article tag
// @Produce application/json
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	var form AddTagForm

	err := app.BindAndValid(c, &form)
	if err != nil {
		app.SendResponse(c, err, nil)
		return
	}
	err = svc.TagSvc.AddTag(form.Name, form.CreatedBy, form.State)
	if err != nil {
		app.SendResponse(c, err, nil)
		return
	}
	app.SendResponse(c, e.OK, nil)
}

// @Summary Update article tag
// @Produce application/json
// @Param id path int true "ID"
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param modified_by body string true "ModifiedBy"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	form := EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}
	err := app.BindAndValid(c, &form)
	if err != nil {
		app.SendResponse(c, err, nil)
		return
	}

	data, err := svc.TagSvc.EditTag(form.ID, form.State, form.Name, form.ModifiedBy)
	if err != nil {
		app.SendResponse(c, err, data)
		return
	}

	app.SendResponse(c, e.OK, data)
}

// @Summary Delete article tag
// @Produce application/json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID must greater than 0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.SendResponse(c, e.ErrBadRequest, nil)
		return
	}

	err := svc.TagSvc.DeleteTag(id)
	if err != nil {
		app.SendResponse(c, err, nil)
		return
	}
	app.SendResponse(c, e.OK, nil)
}
