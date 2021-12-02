package upload

import (
	"bytes"
	"github.com/shipengqi/example.v1/apps/blog/pkg/setting"
	utils2 "github.com/shipengqi/example.v1/apps/blog/pkg/utils"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
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
	fileName = utils2.EncodeMD5(fileName)

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
	ext := utils2.GetExt(fileName)
	for _, allowExt := range setting.AppSettings().ImageAllowExt {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

// CheckImageSize check image size
func CheckImageSize(f multipart.File) bool {
	size, err := utils2.GetSize(f)
	if err != nil {
		return false
	}

	return size <= setting.AppSettings().ImageMaxSize
}

// CheckImage check if the file exists
func CheckImage(src string) error {

	err := utils2.IsNotExistMkDir(src)
	if err != nil {
		return errors.Wrap(err, "utils.IsNotExistMkDir")
	}

	perm := utils2.CheckPermission(src)
	if perm == true {
		return errors.Wrap(err, "utils.CheckPermission")
	}

	return nil
}
