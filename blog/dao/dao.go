package dao

import (
	"gorm.io/gorm"

	"github.com/shipengqi/example.v1/blog/pkg/database/orm"
)

// Dao data access object
type Dao struct {
	db        *gorm.DB
}

// New create instance of Dao
func New() (d *Dao) {
	d = &Dao{
		db:        orm.New(&orm.Config{}),
	}
	return
}

// Ping dao.
//func (d *Dao) Ping() (err error) {
//	return d.db.DB().Ping()
//}

// Close dao.
//func (d *Dao) Close() {
//	if d.db != nil {
//		_ = d.db.Close()
//	}
//}
