package orm

import (
	"time"

	// database driver
	// _ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

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

// New new db and retry connection when has error.
func New(c *Config) (db *gorm.DB) {
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{})
	if err != nil {
		log.Error().Msgf("db dsn(%s) error: %v", c.DSN, err)
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Msgf("db.DB() error: %v", err)
		panic(err)
	}
	if sqlDB == nil {
		log.Warn().Msg("db.DB() get nil")
		return
	}
	// sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(c.Idle)
	// sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(c.Active)
	// sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(c.IdleTimeout)

	return
}
