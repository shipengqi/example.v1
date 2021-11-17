package router

import (
	"github.com/shipengqi/example.v1/apps/blog/pkg/app"
	"github.com/shipengqi/example.v1/apps/blog/pkg/e"
	"github.com/shipengqi/example.v1/apps/blog/pkg/export"
	log "github.com/shipengqi/example.v1/apps/blog/pkg/logger"
	"github.com/shipengqi/example.v1/apps/blog/pkg/qrcode"
	"github.com/shipengqi/example.v1/apps/blog/pkg/setting"
	"github.com/shipengqi/example.v1/apps/blog/pkg/upload"
	"github.com/shipengqi/example.v1/apps/blog/router/api"
	v12 "github.com/shipengqi/example.v1/apps/blog/router/api/v1"
	middleware2 "github.com/shipengqi/example.v1/apps/blog/router/middleware"
	"github.com/shipengqi/example.v1/apps/blog/service"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/shipengqi/example.v1/apps/blog/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
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
	r.Use(middleware2.Logging())

	// 404 Handler.
	r.NoRoute(func(c *gin.Context) {
		app.SendResponse(c, e.ErrNothingFound, "incorrect route")
	})

	gin.SetMode(setting.Settings().RunMode)

	// static file server
	r.StaticFS("/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/login", api.Login)

	v1 := r.Group("/api/v1")
	v1.Use(middleware2.Authenticate(s))
	{
		v1.GET("/tags", v12.GetTags)
		v1.POST("/tags", v12.AddTag)
		v1.PUT("/tags/:id", v12.EditTag)
		v1.DELETE("/tags/:id", v12.DeleteTag)
		v1.POST("/tags/export", v12.ExportTag)
		v1.POST("/tags/import", v12.ImportTag)

		v1.GET("/articles", v12.GetArticles)
		v1.GET("/articles/:id", v12.GetArticle)
		v1.POST("/articles", v12.AddArticle)
		v1.PUT("/articles/:id", v12.EditArticle)
		v1.DELETE("/articles/:id", v12.DeleteArticle)
		v1.POST("/articles/poster/generate", v12.GenerateArticlePoster)

		v1.GET("/users", v12.GetUsers)
		v1.POST("/users", v12.AddUser)
		v1.PUT("/users/:id", v12.EditUser)
		v1.DELETE("/users/:id", v12.DeleteUser)

		v1.POST("/images", v12.UploadImage)
	}

	return r
}
