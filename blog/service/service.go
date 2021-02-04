package service

import (
	"github.com/robfig/cron"

	"github.com/shipengqi/example.v1/blog/dao"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
	"github.com/shipengqi/example.v1/blog/pkg/setting"
	"github.com/shipengqi/example.v1/blog/service/identity"
	"github.com/shipengqi/example.v1/blog/service/rbac"
	"github.com/shipengqi/example.v1/blog/service/tag"
)

type Service struct {
	dao     dao.Interface
	cron    *cron.Cron
	AuthSvc identity.Interface
	TagSvc  tag.Interface
	RBAC    rbac.Interface
}

func New(c *setting.Setting) (s *Service) {
	d := dao.New(c)
	s = &Service{
		dao:     d,
		cron:    cron.New(),
		AuthSvc: identity.New(c.App.SingingKey, d),
		TagSvc:  tag.New(d),
		RBAC:    rbac.New(d),
	}

	if err := s.cron.AddFunc(c.App.PingCron, s.cronPing); err != nil {
		panic(err)
	}

	s.cron.Start()

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

// cronPing check dao health, just a cron example.
func (s *Service) cronPing() {
	log.Trace().Msg("cron ping")
	err := s.dao.Ping()
	if err != nil {
		log.Warn().Err(err).Msg("cron ping")
	}
}