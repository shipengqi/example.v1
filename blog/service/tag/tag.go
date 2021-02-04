package tag

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/shipengqi/example.v1/blog/pkg/export"
	"github.com/shipengqi/example.v1/blog/pkg/utils"
	"github.com/tealeg/xlsx"

	"github.com/shipengqi/example.v1/blog/dao"
	"github.com/shipengqi/example.v1/blog/model"
	"github.com/shipengqi/example.v1/blog/pkg/e"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
	"github.com/shipengqi/example.v1/blog/service/cache"
)

const EXPORT_EXT = "xlsx"

type Interface interface {
	GetTags(maps map[string]interface{}) ([]model.Tag, error)
	AddTag(name, createdBy string) error
	EditTag(id int, name, modifiedBy string) (data map[string]interface{}, err error)
	DeleteTag(id int) (err error)
	Export() (string, error)
	Import(r io.Reader) error
}

type tag struct {
	dao dao.Interface
}

func New(d dao.Interface) Interface {
	return &tag{dao: d}
}

func (t *tag) GetTags(maps map[string]interface{}) ([]model.Tag, error) {

	c := cache.Tag{
		Name:     "",
		State:    0,
		PageNum:  0,
		PageSize: 0,
	}
	key := c.GetTagsCacheKey()
	tagsCache, err := t.dao.GetTagsCache(key)
	if err != nil {
		return nil, e.Wrap(err, "get tags cache")
	}

	if tagsCache != nil {
		log.Info().Msgf("get cache with key: %s", key)
		return tagsCache, nil
	}

	list, err := t.dao.GetTags(0, 10, maps)
	if err != nil {
		return nil, e.Wrap(err, "get tags")
	}

	err = t.dao.SetTagsCache(key, list, 3600)
	if err == nil {
		log.Debug().Msgf("set cache with key: %s", key)
	}
	return list, nil
}

func (t *tag) GetAll() ([]model.Tag, error) {

	c := cache.Tag{
		Name:     "",
		State:    0,
		PageNum:  0,
		PageSize: 0,
	}
	key := c.GetTagsCacheKey()
	tagsCache, err := t.dao.GetTagsCache(key)
	if err != nil {
		return nil, e.Wrap(err, "get tags cache")
	}

	if tagsCache != nil {
		log.Info().Msgf("get cache with key: %s", key)
		return tagsCache, nil
	}

	list, err := t.dao.GetTags(0, 10, map[string]interface{}{
		"deleted":  false,
		"disabled": false,
	})
	if err != nil {
		return nil, e.Wrap(err, "get tags")
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
		return e.Wrap(err, "exist tag")
	}

	if exists {
		return e.ErrExistTag
	}

	return t.dao.AddTag(name, createdBy)
}

func (t *tag) EditTag(id int, name, modifiedBy string) (data map[string]interface{}, err error) {
	exists, err := t.dao.ExistTagByID(id)
	if err != nil {
		return nil, e.Wrap(err, "exist tag")
	}
	if !exists {
		return nil, e.ErrNotExistTag
	}
	data = make(map[string]interface{})
	data["modified_by"] = modifiedBy
	if name != "" {
		data["name"] = name
	}

	err = t.dao.EditTag(id, data)
	if err != nil {
		return nil, e.Wrap(err, "edit tag")
	}
	return data, nil
}

func (t *tag) DeleteTag(id int) (err error) {
	exists, err := t.dao.ExistTagByID(id)
	if err != nil {
		return e.Wrap(err, "exist tag")
	}
	if !exists {
		return e.ErrNotExistTag
	}

	return t.dao.DeleteTag(id)
}

func (t *tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

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

	err = xlsFile.Save(dirFullPath + filename)
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

			}
		}
	}
}
