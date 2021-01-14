package dao

import (
	"encoding/json"

	"gorm.io/gorm"

	"github.com/shipengqi/example.v1/blog/model"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
)

func (d *Dao) SetTagsCache(key string, data interface{}, exp int) error {

	if err := d.redis.Set(key, data, exp); err != nil {
		log.Error().Err(err).Msgf("set cache with key: %s", key)
		return err
	}

	return nil
}

func (d *Dao) GetTagsCache(key string) ([]model.Tag, error) {
	var (
		tags []model.Tag
		err  error
	)

	if d.redis.Exists(key) {
		data, err := d.redis.Get(key)
		if err == nil {
			err = json.Unmarshal(data, &tags)
			if err != nil {
				log.Error().Err(err).Msgf("unmarshal cache with key: %s", key)
				return nil, err
			}
			return tags, nil
		}
		log.Warn().Msgf("no cache with key: %s", key)
	}

	return nil, err
}

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

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

func (d *Dao) GetTagTotal(maps interface{}) (int64, error) {
	var count int64
	if err := d.db.Model(&model.Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}



func (d *Dao) AddTag(name string, state int, createdBy string) error {
	tag := &model.Tag{
		Name:      name,
		CreatedBy: createdBy,
		State:     state,
	}
	if err := d.db.Create(tag).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) DeleteTag(id int) error {
	if err := d.db.Where("id = ?", id).Delete(&model.Tag{}).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) EditTag(id int, data interface{}) error {
	if err := d.db.Model(&model.Tag{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (d *Dao) ExistTagByName(name string) (bool, error) {
	var tag model.Tag
	err := d.db.Select("id").Where("name = ?", name).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

func (d *Dao) ExistTagByID(id int) (bool, error) {
	var tag model.Tag
	err := d.db.Select("id").Where("id = ?", id).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}


