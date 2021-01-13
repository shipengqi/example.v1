package service

import (
	"github.com/shipengqi/example.v1/blog/pkg/errno"
)

func (s *Service) GetTags(maps map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	list, err := s.dao.GetTags(0, 10, maps)
	if err != nil {
		return nil, errno.Wrap(err, "get tags")
	}
	total, err := s.dao.GetTagTotal(maps)
	if err != nil {
		return nil, errno.Wrap(err, "get total")
	}
	data["lists"] = list
	data["total"] = total
	return data, nil
}

func (s *Service) AddTag(name, createdBy string, state int) error {
	exists, err := s.dao.ExistTagByName(name)
	if err != nil {
		return errno.Wrap(err, "exist tag")
	}

	if !exists {
		return errno.ErrExistTag
	}

	return s.dao.AddTag(name, state, createdBy)
}

func (s *Service) EditTag(id, state int, name, modifiedBy string) (data map[string]interface{}, err error) {
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

func (s *Service) DeleteTag(id int) (err error) {
	exists, err := s.dao.ExistTagByID(id)
	if err != nil {
		return errno.Wrap(err, "exist tag")
	}
	if !exists {
		return errno.ErrNotExistTag
	}

	return s.dao.DeleteTag(id)
}
