package router

import (
	"github.com/gin-gonic/gin"

	"github.com/shipengqi/example.v1/blog/pkg/setting"
	apiv1 "github.com/shipengqi/example.v1/blog/router/api/v1"
)

func Init() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.Settings().RunMode)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/tags", apiv1.GetTags)
		v1.POST("/tags", apiv1.AddTag)
		v1.PUT("/tags/:id", apiv1.EditTag)
		v1.DELETE("/tags/:id", apiv1.DeleteTag)
	}

	return r
}
