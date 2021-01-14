package service

import (
	"github.com/shipengqi/example.v1/blog/dao"
	"github.com/shipengqi/example.v1/blog/pkg/jwt"
	"github.com/shipengqi/example.v1/blog/pkg/setting"
)

type Service struct {
	dao *dao.Dao
	jwt *jwt.Jwt
}

func New(c *setting.Setting) (s *Service) {
	s = &Service{
		dao: dao.New(c),
		jwt: jwt.New(c.App.SingingKey),
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
