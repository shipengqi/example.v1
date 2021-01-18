package tag

import (
	"github.com/shipengqi/example.v1/blog/dao"
	"github.com/shipengqi/example.v1/blog/model"
	"github.com/shipengqi/example.v1/blog/pkg/errno"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
	"github.com/shipengqi/example.v1/blog/service/cache"
)

type Svc struct {
	dao dao.Dao
}

func New(d dao.Dao) *Svc {
	return &Svc{dao: d}
}

func (s *Svc) GetTags(maps map[string]interface{}) ([]model.Tag, error) {

	c := cache.Tag{
		Name:     "",
		State:    0,
		PageNum:  0,
		PageSize: 0,
	}
	key := c.GetTagsCacheKey()
	tagsCache, err := s.dao.GetTagsCache(key)
	if err != nil {
		return nil, errno.Wrap(err, "get tags cache")
	}

	if tagsCache != nil {
		log.Info().Msgf("get cache with key: %s", key)
		return tagsCache, nil
	}

	list, err := s.dao.GetTags(0, 10, maps)
	if err != nil {
		return nil, errno.Wrap(err, "get tags")
	}

	err = s.dao.SetTagsCache(key, list, 3600)
	if err == nil {
		log.Debug().Msgf("set cache with key: %s", key)
	}
	return list, nil
}

func (s *Svc) AddTag(name, createdBy string, state int) error {
	exists, err := s.dao.ExistTagByName(name)
	if err != nil {
		return errno.Wrap(err, "exist tag")
	}

	if !exists {
		return errno.ErrExistTag
	}

	return s.dao.AddTag(name, state, createdBy)
}

func (s *Svc) EditTag(id, state int, name, modifiedBy string) (data map[string]interface{}, err error) {
	exists, err := s.dao.ExistTagByID(id)
	if err != nil {
		return nil, errno.Wrap(err, "exist tag")
	}
	if !exists {
		return nil, errno.ErrNotExistTag
	}
	data = make(map[string]interface{})
	data["modified_by"] = modifiedBy
	if name != "" {
		data["name"] = name
	}
	if state != -1 {
		data["state"] = state
	}

	err = s.dao.EditTag(id, data)
	if err != nil {
		return nil, errno.Wrap(err, "edit tag")
	}
	return data, nil
}

func (s *Svc) DeleteTag(id int) (err error) {
	exists, err := s.dao.ExistTagByID(id)
	if err != nil {
		return errno.Wrap(err, "exist tag")
	}
	if !exists {
		return errno.ErrNotExistTag
	}

	return s.dao.DeleteTag(id)
}
