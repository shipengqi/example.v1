package tag

import (
	"github.com/shipengqi/example.v1/blog/dao"
	"github.com/shipengqi/example.v1/blog/model"
	"github.com/shipengqi/example.v1/blog/pkg/e"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
	"github.com/shipengqi/example.v1/blog/service/cache"
)

type Interface interface {
	GetTags(maps map[string]interface{}) ([]model.Tag, error)
	AddTag(name, createdBy string, state int) error
	EditTag(id, state int, name, modifiedBy string) (data map[string]interface{}, err error)
	DeleteTag(id int) (err error)
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

func (t *tag) AddTag(name, createdBy string, state int) error {
	exists, err := t.dao.ExistTagByName(name)
	if err != nil {
		return e.Wrap(err, "exist tag")
	}

	if !exists {
		return e.ErrExistTag
	}

	return t.dao.AddTag(name, state, createdBy)
}

func (t *tag) EditTag(id, state int, name, modifiedBy string) (data map[string]interface{}, err error) {
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
	if state != -1 {
		data["state"] = state
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
