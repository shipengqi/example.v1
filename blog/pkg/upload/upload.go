package upload

import (
	"bytes"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/blog/pkg/setting"
	"github.com/shipengqi/example.v1/blog/pkg/utils"
)

// GetImageFullUrl get the full access path
func GetImageFullUrl(name string) string {
	var buffer bytes.Buffer
	buffer.WriteString(setting.AppSettings().RootEndpoint)
	buffer.WriteString(GetImagePath())
	buffer.WriteString("/")
	buffer.WriteString(name)
	return buffer.String()
}

// GetImageName get image name
func GetImageName(name string) string {
	ext := filepath.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = utils.EncodeMD5(fileName)

	return fileName + ext
}

// GetImagePath get save path
func GetImagePath() string {
	return setting.AppSettings().ImageSavePath
}

// GetImageFullPath get full save path
func GetImageFullPath() string {
	return setting.AppSettings().FileRootPath + GetImagePath()
}

// CheckImageExt check image file ext
func CheckImageExt(fileName string) bool {
	ext := utils.GetExt(fileName)
	for _, allowExt := range setting.AppSettings().ImageAllowExt {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

// CheckImageSize check image size
func CheckImageSize(f multipart.File) bool {
	size, err := utils.GetSize(f)
	if err != nil {
		return false
	}

	return size <= setting.AppSettings().ImageMaxSize
}

// CheckImage check if the file exists
func CheckImage(src string) error {

	err := utils.IsNotExistMkDir(src)
	if err != nil {
		return errors.Wrap(err, "utils.IsNotExistMkDir")
	}

	perm := utils.CheckPermission(src)
	if perm == true {
		return errors.Wrap(err, "utils.CheckPermission")
	}

	return nil
}