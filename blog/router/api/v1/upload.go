package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/shipengqi/example.v1/blog/pkg/app"
	"github.com/shipengqi/example.v1/blog/pkg/e"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
	"github.com/shipengqi/example.v1/blog/pkg/upload"
)

// @Summary Upload Image
// @Produce  json
// @Param image formData file true "Image File"
// @Success 200 {object} app.Response
// @Router /api/v1/images [post]
func UploadImage(c *gin.Context) {
	// upload multi files
	// form, _ := c.MultipartForm()
	// files := form.File["images"]
	//
	// for _, file := range files {
	// 	c.SaveUploadedFile(file, dst)
	// }

	// upload one file
	file, image, err := c.Request.FormFile("images")
	if err != nil {
		app.SendResponse(c, e.ErrMultiFormErr, nil)
		return
	}
	log.Debug().Msgf("upload image: %s", image.Filename)

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fmt.Sprintf("%s/%s", fullPath, imageName)

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		app.SendResponse(c, e.ErrCheckImage, nil)
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		log.Warn().Err(err)
		app.SendResponse(c, e.ErrCheckImage, nil)
		return
	}

	// check image exists
	// 上传文件至指定目录
	err = c.SaveUploadedFile(image, src)
	if err != nil {
		log.Warn().Err(err)
		app.SendResponse(c, e.ErrUploadImage, nil)
		return
	}

	log.Debug().Msgf("image: %s uploaded!", image.Filename)
	app.SendResponse(c, e.OK, map[string]string{
		"image_url":      upload.GetImageFullUrl(imageName),
		"image_save_url": fmt.Sprintf("%s/%s", savePath, imageName),
	})
}

