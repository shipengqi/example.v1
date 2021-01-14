package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/shipengqi/example.v1/blog/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/shipengqi/example.v1/blog/middleware"
	"github.com/shipengqi/example.v1/blog/pkg/app"
	"github.com/shipengqi/example.v1/blog/pkg/errno"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
	"github.com/shipengqi/example.v1/blog/pkg/setting"
	"github.com/shipengqi/example.v1/blog/router/api"
	apiv1 "github.com/shipengqi/example.v1/blog/router/api/v1"
	"github.com/shipengqi/example.v1/blog/service"
)

func Init(s *service.Service) *gin.Engine {
	err := s.Ping()
	if err != nil {
		log.Fatal().Msgf("Ping err: %s", err)
	}
	r := gin.New()
	api.Init(s)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Logging())

	// 404 Handler.
	r.NoRoute(func(c *gin.Context) {
		app.SendResponse(c, errno.ErrNothingFound, "incorrect route")
	})

	gin.SetMode(setting.Settings().RunMode)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("login", api.Login)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/tags", apiv1.GetTags)
		v1.POST("/tags", apiv1.AddTag)
		v1.PUT("/tags/:id", apiv1.EditTag)
		v1.DELETE("/tags/:id", apiv1.DeleteTag)

		v1.GET("/articles", apiv1.GetArticles)
		v1.GET("/articles/:id", apiv1.GetArticle)
		v1.POST("/articles", apiv1.AddArticle)
		v1.PUT("/articles/:id", apiv1.EditArticle)
		v1.DELETE("/articles/:id", apiv1.DeleteArticle)
		v1.POST("/articles/poster/generate", apiv1.GenerateArticlePoster)
	}

	return r
}
