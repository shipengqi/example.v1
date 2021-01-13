package dao

import (
	"gorm.io/gorm"

	"github.com/shipengqi/example.v1/blog/pkg/cache/gredis"
	"github.com/shipengqi/example.v1/blog/pkg/database/orm"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
	"github.com/shipengqi/example.v1/blog/pkg/setting"
)

// Dao data access object
type Dao struct {
	db    *gorm.DB
	redis *gredis.Pool
}

// New create instance of Dao
func New(c *setting.Setting) (d *Dao) {
	d = &Dao{
		db:    orm.New(c.DB),
		redis: gredis.New(c.Redis),
	}
	return
}

// Ping dao.
func (d *Dao) Ping() (err error) {
	sqlDB, err := d.db.DB()
	if err != nil {
		return
	}
	if err = sqlDB.Ping(); err != nil {
		return
	}
	if err = d.redis.Ping(); err != nil {
		return
	}
	return
}

// Close dao.
func (d *Dao) Close() {
	if d.db != nil {
		sqlDB, err := d.db.DB()
		if err != nil {
			log.Warn().Msgf("db.DB() err: %v", err)
			return
		}
		_ = sqlDB.Close()
	}

	if d.redis != nil {
		_ = d.redis.Close()
	}
}
