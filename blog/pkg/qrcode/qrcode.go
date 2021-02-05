package qrcode

import (
	"bytes"
	"fmt"
	"image/jpeg"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/shipengqi/example.v1/blog/pkg/utils"

	"github.com/shipengqi/example.v1/blog/pkg/setting"
)

type QrCode struct {
	Name   string
	URL    string
	Width  int
	Height int
	Ext    string
	Level  qr.ErrorCorrectionLevel
	Mode   qr.Encoding
}

const (
	EXT_JPG = "jpg"
)

func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		Name:   GetQrCodeFileName(url),
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPG,
	}
}

func GetQrCodePath() string {
	return setting.AppSettings().QrCodeSavePath
}

func GetQrCodeFullPath() string {
	return setting.AppSettings().FileRootPath + setting.AppSettings().QrCodeSavePath
}

func GetQrCodeFullUrl(name string) string {
	var buffer bytes.Buffer
	buffer.WriteString(setting.AppSettings().RootEndpoint)
	buffer.WriteString(GetQrCodePath())
	buffer.WriteString("/")
	buffer.WriteString(name)
	return buffer.String()
}

// GetQrCodeFileName get qr file name
func GetQrCodeFileName(value string) string {
	return utils.EncodeMD5(value)
}

func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

// Encode generate QR code
func (q *QrCode) Encode(path string) (string, error) {
	name := fmt.Sprintf("%s.%s", q.Name, q.GetQrCodeExt())
	src := fmt.Sprintf("%s/%s", path, name)
	if utils.IsExist(src) {
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", err
		}

		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", err
		}

		f, err := utils.MustOpen(path, name)
		if err != nil {
			return "", err
		}
		defer f.Close()

		// 第三个参数可设置其图像质量，默认值为 75
		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", err
		}
	}

	return name, nil
}
