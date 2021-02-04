package export

import (
	"bytes"

	"github.com/shipengqi/example.v1/blog/pkg/setting"
)

const EXT = ".xlsx"

// GetExcelFullUrl get the full access path of the Excel file
func GetExcelFullUrl(name string) string {
	var buffer bytes.Buffer
	buffer.WriteString(setting.AppSettings().RootEndpoint)
	buffer.WriteString(GetExcelPath())
	buffer.WriteString("/")
	buffer.WriteString(name)
	return buffer.String()
}

// GetExcelPath get the relative save path of the Excel file
func GetExcelPath() string {
	return setting.AppSettings().ExportSavePath
}

// GetExcelFullPath Get the full save path of the Excel file
func GetExcelFullPath() string {
	return setting.AppSettings().FileRootPath + GetExcelPath()
}
