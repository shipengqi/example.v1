package v1

import (
	"fmt"

	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/shipengqi/example.v1/blog/service/article"

	"github.com/shipengqi/example.v1/blog/pkg/app"
	"github.com/shipengqi/example.v1/blog/pkg/e"
	"github.com/shipengqi/example.v1/blog/pkg/qrcode"
)

const (
	QRCODE_URL = "https://github.com/shipengqi/example.v1"
)

type GenArticleRequest struct {
	Name string `form:"name" valid:"MaxSize(100)"`
}

type GenArticleResponse struct {
	PosterUrl     string `json:"poster_url"`
	PosterSaveUrl string `json:"poster_save_url"`
}

func GetArticles(c *gin.Context) {

}

func GetArticle(c *gin.Context) {

}

func AddArticle(c *gin.Context) {

}

func EditArticle(c *gin.Context) {

}

func DeleteArticle(c *gin.Context) {

}

func GenerateArticlePoster(c *gin.Context) {
	qrc := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	posterName := fmt.Sprintf("%s-%s.%s", article.POSTER_FLAG, qrc.Name, qrc.GetQrCodeExt())
	bg := article.NewPosterBg(
		"bg.jpg",
		&article.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&article.Pt{
			X: 125,
			Y: 298,
		},
	)
	poster := article.NewPoster(posterName, qrc, bg)

	// path := qrcode.GetQrCodeFullPath()
	// name, err := qrc.Encode(path)
	// if err != nil {
	// 	app.SendResponse(c, err, nil)
	// 	return
	// }

	_, err := poster.Generate()
	if err != nil {
		app.SendResponse(c, err, nil)
		return
	}

	app.SendResponse(c, e.OK, GenArticleResponse{
		PosterUrl:     qrcode.GetQrCodeFullUrl(posterName),
		PosterSaveUrl: fmt.Sprintf("%s/%s", qrcode.GetQrCodePath(), posterName),
	})
	return
}
