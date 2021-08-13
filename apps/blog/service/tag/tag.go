package tag

import (
	"fmt"
	"github.com/shipengqi/example.v1/apps/blog/dao"
	"github.com/shipengqi/example.v1/apps/blog/model"
	e2 "github.com/shipengqi/example.v1/apps/blog/pkg/e"
	"github.com/shipengqi/example.v1/apps/blog/pkg/export"
	log "github.com/shipengqi/example.v1/apps/blog/pkg/logger"
	"github.com/shipengqi/example.v1/apps/blog/pkg/setting"
	"github.com/shipengqi/example.v1/apps/blog/pkg/utils"
	"github.com/shipengqi/example.v1/apps/blog/service/cache"
	"io"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tealeg/xlsx"
)

const EXPORT_EXT = "xlsx"

type Interface interface {
	GetTags(maps map[string]interface{}, page int) ([]model.Tag, error)
	AddTag(name, createdBy string) error
	EditTag(id int, name, modifiedBy string) (data map[string]interface{}, err error)
	DeleteTag(id int) (err error)
	Export(name string) (string, error)
	Import(r io.Reader) error
}

type tag struct {
	dao dao.Interface
}

func New(d dao.Interface) Interface {
	return &tag{dao: d}
}

func (t *tag) GetTags(maps map[string]interface{}, page int) ([]model.Tag, error) {

	var name string
	if n, ok := maps["name"].(string); ok {
		name = n
	}

	pageNum := t.getPage(page)
	c := cache.Tag{
		Name:     name,
		State:    0,
		PageNum:  pageNum,
		PageSize: setting.AppSettings().PageSize,
	}
	key := c.GetTagsCacheKey()
	tagsCache, err := t.dao.GetTagsCache(key)
	if err != nil {
		return nil, e2.Wrap(err, "get tags cache")
	}

	if tagsCache != nil && len(tagsCache) != 0 {
		log.Info().Msgf("get cache with key: %s", key)
		return tagsCache, nil
	}

	list, err := t.dao.GetTags(pageNum, setting.AppSettings().PageSize, maps)
	if err != nil {
		return nil, e2.Wrap(err, "get tags")
	}

	err = t.dao.SetTagsCache(key, list, 3600)
	if err == nil {
		log.Debug().Msgf("set cache with key: %s", key)
	}
	return list, nil
}

func (t *tag) AddTag(name, createdBy string) error {
	exists, err := t.dao.ExistTagByName(name)
	if err != nil {
		return e2.Wrap(err, "exist tag")
	}

	if exists {
		return e2.ErrExistTag
	}

	return t.dao.AddTag(name, createdBy)
}

func (t *tag) EditTag(id int, name, modifiedBy string) (data map[string]interface{}, err error) {
	exists, err := t.dao.ExistTagByID(id)
	if err != nil {
		return nil, e2.Wrap(err, "exist tag")
	}
	if !exists {
		return nil, e2.ErrNotExistTag
	}
	data = make(map[string]interface{})
	data["modified_by"] = modifiedBy
	if name != "" {
		data["name"] = name
	}

	err = t.dao.EditTag(id, data)
	if err != nil {
		return nil, e2.Wrap(err, "edit tag")
	}
	return data, nil
}

func (t *tag) DeleteTag(id int) (err error) {
	exists, err := t.dao.ExistTagByID(id)
	if err != nil {
		return e2.Wrap(err, "exist tag")
	}
	if !exists {
		return e2.ErrNotExistTag
	}

	return t.dao.DeleteTag(id)
}

func (t *tag) Export(name string) (string, error) {
	maps := t.getMaps(name)
	tags, err := t.GetTags(maps, 0)
	if err != nil {
		return "", err
	}

	// f := excelize.NewFile()
	// Create a new sheet.
	// index := f.NewSheet("Sheet2")
	// Set value of a cell.
	// f.SetCellValue("Sheet2", "A2", "Hello world.")
	// f.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	// f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	// if err := f.SaveAs("Book1.xlsx"); err != nil {
	// 	fmt.Println(err)
	// }

	xlsFile := xlsx.NewFile()
	sheet, err := xlsFile.AddSheet("Tag Information")
	if err != nil {
		return "", err
	}

	titles := []string{"ID", "Name", "Created By", "Created At", "Updated By", "Updated At"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range tags {
		values := []string{
			strconv.Itoa(int(v.ID)),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(int(v.CreatedAt)),
			v.ModifiedBy,
			strconv.Itoa(int(v.UpdatedAt)),
		}

		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	genTime := time.Now().Unix()
	filename := fmt.Sprintf("tags-%d.%s", genTime, EXPORT_EXT)

	dirFullPath := export.GetExcelFullPath()
	err = utils.IsNotExistMkDir(dirFullPath)
	if err != nil {
		return "", err
	}

	err = xlsFile.Save(fmt.Sprintf("%s/%s", dirFullPath, filename))
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (t *tag) Import(r io.Reader) error {
	file, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows := file.GetRows("Tag Information")
	// remove duplicate name from file and database
	for i, row := range rows {
		if i > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}
			if len(data) == 0 || len(data) < 3 {
				continue
			}
			err := t.dao.AddTag(data[1], data[2])
			if err != nil {
				log.Warn().Err(err).Msgf("add cell: %s", data[1])
			}
		}
	}

	return nil
}


// getPage get page parameters
func (t *tag) getPage(page int) int {
	result := 0
	if page > 0 {
		result = (page - 1) * setting.AppSettings().PageSize
	}

	return result
}

func (t *tag) getMaps(name string) map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_at"] = false

	if len(name) > 0 {
		maps["name"] = name
	}

	return maps
}
