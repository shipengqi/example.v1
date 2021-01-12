package dao

import (
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"

	"github.com/shipengqi/example.v1/blog/pkg/database/gredis"
	"github.com/shipengqi/example.v1/blog/pkg/database/orm"
	log "github.com/shipengqi/example.v1/blog/pkg/logger"
	"github.com/shipengqi/example.v1/blog/pkg/setting"
)

// Dao data access object
type Dao struct {
	db    *gorm.DB
	redis *redis.Pool
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
	if err = d.pingRedis(); err != nil {
		return
	}
	return
}

// Close dao.
func (d *Dao) Close() {
	if d.db != nil {
		sqlDB, err := d.db.DB()
		if err != nil {
			return
		}
		log.Warn().Msgf("db.DB() err: %v", err)
		_ = sqlDB.Close()
	}

	if d.redis != nil {
		_ = d.redis.Close()
	}
}

func (d *Dao) pingRedis() (err error) {
	conn := d.redis.Get()
	defer conn.Close()
	_, err = conn.Do("SET", "PING", "PONG")
	return
}