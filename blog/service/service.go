package service

import (
	"github.com/shipengqi/example.v1/blog/dao"
	"github.com/shipengqi/example.v1/blog/pkg/setting"
	"github.com/shipengqi/example.v1/blog/service/identity"
	"github.com/shipengqi/example.v1/blog/service/tag"
)

type Service struct {
	dao     dao.Dao
	AuthSvc *identity.Svc
	TagSvc  *tag.Svc
}

func New(c *setting.Setting) (s *Service) {
	d := dao.New(c)
	s = &Service{
		dao:     d,
		AuthSvc: identity.New(c.App.SingingKey),
		TagSvc:  tag.New(d),
	}
	return
}

// Ping check dao health.
func (s *Service) Ping() (err error) {
	return s.dao.Ping()
}

// Close close all dao.
func (s *Service) Close() {
	s.dao.Close()
}
