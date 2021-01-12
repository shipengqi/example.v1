package dao

import (
	"gorm.io/gorm"

	"github.com/shipengqi/example.v1/blog/model"
)

func (d *Dao) GetTags(pageNum int, pageSize int, maps interface{}) ([]model.Tag, error) {
	var (
		tags []model.Tag
		err  error
	)

	if pageSize > 0 && pageNum > 0 {
		err = d.db.Where(maps).Find(&tags).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = d.db.Where(maps).Find(&tags).Error
	}

	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

func (d *Dao) GetTagTotal(maps interface{}) (count int64) {
	d.db.Model(&model.Tag{}).Where(maps).Count(&count)

	return
}

func (d *Dao) ExistTagByName(name string) bool {
	var tag model.Tag
	d.db.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func (d *Dao) AddTag(name string, state int, createdBy string) bool {
	d.db.Create(&model.Tag{
		Name:      name,
		CreatedBy: createdBy,
		State:     state,
	})

	return true
}

func (d *Dao) ExistTagByID(id int) bool {
	var tag model.Tag
	d.db.Select("id").Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func (d *Dao) DeleteTag(id int) bool {
	d.db.Where("id = ?", id).Delete(&model.Tag{})

	return true
}

func (d *Dao) EditTag(id int, data interface{}) bool {
	d.db.Model(&model.Tag{}).Where("id = ?", id).Updates(data)

	return true
}
