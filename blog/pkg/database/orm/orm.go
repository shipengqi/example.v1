package orm

import (
	"strings"
	"time"

	// database driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	log "github.com/shipengqi/example.v1/blog/pkg/logger"
)

type Config struct {
	TablePrefix string        `ini:"TABLE_PREFIX"`
	DbType      string        `ini:"TYPE"`
	DSN         string        `ini:"DSN"`          // data source name.
	Active      int           `ini:"ACTIVE"`       // pool
	Idle        int           `ini:"IDLE"`         // pool
	IdleTimeout time.Duration `ini:"IDLE_TIMEOUT"` // connect max life time.
}

type ormLog struct{}

func (l ormLog) Print(v ...interface{}) {
	log.Info().Msgf(strings.Repeat("%v ", len(v)), v...)
}

// New new db and retry connection when has error.
func New(c *Config) (db *gorm.DB) {
	db, err := gorm.Open(c.DbType, c.DSN)
	if err != nil {
		log.Error().Msgf("db dsn(%s) error: %v", c.DSN, err)
		panic(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return c.TablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(c.Idle)
	db.DB().SetMaxOpenConns(c.Active)
	db.DB().SetConnMaxLifetime(c.IdleTimeout)
	// db.SetLogger(ormLog{})
	return
}
