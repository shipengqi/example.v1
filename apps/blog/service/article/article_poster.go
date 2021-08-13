package article

import (
	"fmt"
	"github.com/shipengqi/example.v1/apps/blog/pkg/qrcode"
	"github.com/shipengqi/example.v1/apps/blog/pkg/setting"
	"github.com/shipengqi/example.v1/apps/blog/pkg/utils"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
)

const (
	POSTER_FLAG = "poster"
)

type Poster interface {
	Generate() (string, error)
}

type poster struct {
	name string
	qr   *qrcode.QrCode
	bg   *PosterBg
}

type PosterBg struct {
	Name string
	*Rect
	*Pt
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

type DrawText struct {
	JPG    draw.Image
	Merged *os.File

	Title string
	X0    int
	Y0    int
	Size0 float64

	SubTitle string
	X1       int
	Y1       int
	Size1    float64
}

func NewPosterBg(name string, rect *Rect, pt *Pt) *PosterBg {
	return &PosterBg{
		Name:          name,
		Rect:          rect,
		Pt:            pt,
	}
}

func NewPoster(name string, qr *qrcode.QrCode, bg *PosterBg) Poster {
	return &poster{
		name: name,
		qr:   qr,
		bg:   bg,
	}
}

func (p *poster) checkMergedImage(path string) bool {
	if utils.IsExist(fmt.Sprintf("%s/%s", path, p.name)) {
		return false
	}

	return true
}

func (p *poster) openMergedImage(path string) (*os.File, error) {
	f, err := utils.MustOpen(path, p.name)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (p *poster) drawText(d *DrawText, fontName string) error {
	fontSource := fmt.Sprintf("%s%s/%s",
		setting.AppSettings().FileRootPath,
		setting.AppSettings().FontSavePath,
		fontName)
	fontSourceBytes, err := ioutil.ReadFile(fontSource)
	if err != nil {
		return err
	}

	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)
	if err != nil {
		return err
	}

	fc := freetype.NewContext()
	// 设置屏幕每英寸的分辨率
	fc.SetDPI(72)
	// 设置用于绘制文本的字体
	fc.SetFont(trueTypeFont)
	// 以磅为单位设置字体大小
	fc.SetFontSize(d.Size0)
	// 设置剪裁矩形以进行绘制
	fc.SetClip(d.JPG.Bounds())
	// 设置目标图像
	fc.SetDst(d.JPG)
	// 设置绘制操作的源图像，通常为 image.Uniform
	fc.SetSrc(image.Black)

	// 根据 Pt 的坐标值绘制给定的文本内容
	pt := freetype.Pt(d.X0, d.Y0)
	_, err = fc.DrawString(d.Title, pt)
	if err != nil {
		return err
	}

	fc.SetFontSize(d.Size1)
	_, err = fc.DrawString(d.SubTitle, freetype.Pt(d.X1, d.Y1))
	if err != nil {
		return err
	}

	err = jpeg.Encode(d.Merged, d.JPG, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *poster) Generate() (string, error) {
	// 获取二维码存储路径
	fullPath := qrcode.GetQrCodeFullPath()
	// 生成二维码
	qrName, err := p.qr.Encode(fullPath)
	if err != nil {
		return "", err
	}

	// 检查合并后图像（指的是存放合并后的海报）是否存在
	if !p.checkMergedImage(fullPath) { // 若不存在，则生成
		mergedF, err := p.openMergedImage(fullPath)
		if err != nil {
			return "", err
		}
		defer mergedF.Close()

		// 打开背景图
		bgF, err := utils.MustOpen(fullPath, p.bg.Name)
		if err != nil {
			return "", err
		}
		defer bgF.Close()

		// 打开二维码图像
		qrF, err := utils.MustOpen(fullPath, qrName)
		if err != nil {
			return "", err
		}
		defer qrF.Close()

		// 解码 bgF 和 qrF 返回 image.Image
		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			return "", err
		}
		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			return "", err
		}

		// 创建一个新的 RGBA 图像
		jpg := image.NewRGBA(image.Rect(p.bg.Rect.X0, p.bg.Rect.Y0, p.bg.Rect.X1, p.bg.Rect.Y1))

		// 在 RGBA 图像上绘制 背景图 (bgF)
		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		// 在已绘制背景图的 RGBA 图像上，在指定 Point 上绘制二维码图像（qrF）
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(p.bg.Pt.X, p.bg.Pt.Y)), draw.Over)

		err = p.drawText(&DrawText{
			JPG:    jpg,
			Merged: mergedF,

			Title: "Golang blog example",
			X0:    80,
			Y0:    160,
			Size0: 42,

			SubTitle: "--- Pooky",
			X1:       320,
			Y1:       220,
			Size1:    36,
		}, "msyhbd.ttc")

		err = jpeg.Encode(mergedF, jpg, nil)
		if err != nil {
			return "", err
		}
	}

	return qrName, nil
}
