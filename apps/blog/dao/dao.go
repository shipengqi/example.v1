package dao

import (
	model2 "github.com/shipengqi/example.v1/apps/blog/model"
	"github.com/shipengqi/example.v1/apps/blog/pkg/cache/gredis"
	"github.com/shipengqi/example.v1/apps/blog/pkg/database/orm"
	log "github.com/shipengqi/example.v1/apps/blog/pkg/logger"
	"github.com/shipengqi/example.v1/apps/blog/pkg/setting"
	"gorm.io/gorm"
)

type Interface interface {
	Ping() (err error)
	Close()

	SetTagsCache(key string, data interface{}, exp int) error
	GetTagsCache(key string) ([]model2.Tag, error)
	GetTags(pageNum int, pageSize int, maps interface{}) ([]model2.Tag, error)
	GetTagTotal(maps interface{}) (int64, error)
	AddTag(name string, createdBy string) error
	DeleteTag(id int) error
	EditTag(id int, data interface{}) error
	ExistTagByName(name string) (bool, error)
	ExistTagByID(id int) (bool, error)

	GetUserRbac(userid uint) (*model2.UserRBAC, error)
	GetPermissionsWithRoles(roles []model2.Role) ([]model2.UserPermission, error)
	AddUser(name, pass, phone, email string) error
	GetUser(username string) (*model2.User, error)
}

// Dao data access object
type dao struct {
	db    *gorm.DB
	redis gredis.Pool
}

// New create instance of Dao
func New(c *setting.Setting) Interface {
	return &dao{
		db:    orm.New(c.DB),
		redis: gredis.New(c.Redis),
	}
}

// Ping dao.
func (d *dao) Ping() (err error) {
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
func (d *dao) Close() {
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
