package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/unknwon/com"

	"github.com/shipengqi/example.v1/blog/model"
	"github.com/shipengqi/example.v1/blog/pkg/app"
	"github.com/shipengqi/example.v1/blog/pkg/errno"
	"github.com/shipengqi/example.v1/blog/service"
)

var s *service.Service

// @Summary Get multiple article tags
// @Produce  json
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

	data := s.GetTags(maps)
	app.SendResponse(c, errno.OK, data)
}

// @Summary Add article tag
// @Produce  json
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.SendResponse(c, errno.ErrBadRequest, nil)
		return
	}

	if model.ExistTagByName(name) {
		app.SendResponse(c, errno.ErrExistTag, nil)
		return
	}

	model.AddTag(name, state, createdBy)

	app.SendResponse(c, errno.OK, nil)
}

// @Summary Update article tag
// @Produce  json
// @Param id path int true "ID"
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param modified_by body string true "ModifiedBy"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.SendResponse(c, errno.ErrBadRequest, nil)
		return
	}

	if !model.ExistTagByID(id) {
		app.SendResponse(c, errno.ErrNotExistTag, nil)
		return
	}

	data := make(map[string]interface{})
	data["modified_by"] = modifiedBy
	if name != "" {
		data["name"] = name
	}
	if state != -1 {
		data["state"] = state
	}

	model.EditTag(id, data)

	app.SendResponse(c, errno.OK, data)
}

// @Summary Delete article tag
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 200 {object} app.Response
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.SendResponse(c, errno.ErrBadRequest, nil)
		return
	}

	if !model.ExistTagByID(id) {
		app.SendResponse(c, errno.ErrNotExistTag, nil)
		return
	}

	model.DeleteTag(id)

	app.SendResponse(c, errno.OK, nil)
}
