package orm

import (
	"context"
	log "github.com/shipengqi/example.v1/apps/blog/pkg/logger"
	"time"

	// database driver
	// _ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	ol "gorm.io/gorm/logger"
)

type Config struct {
	DbType        string        `ini:"TYPE"`
	DSN           string        `ini:"DSN"`            // data source name.
	Active        int           `ini:"ACTIVE"`         // pool
	Idle          int           `ini:"IDLE"`           // pool
	IdleTimeout   time.Duration `ini:"IDLE_TIMEOUT"`   // connect max life time.
	SlowThreshold time.Duration `ini:"SLOW_THRESHOLD"` // slow log threshold
}

// ----------------------------------------------------------------------------
// Implements orm logger interface
type ormLogger struct {
	slowThreshold time.Duration
}

func (o *ormLogger) LogMode(ol.LogLevel) ol.Interface {
	return o
}

func (o *ormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	log.Info().Msgf(str, args)
}

func (o *ormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	log.Warn().Msgf(str, args)
}

func (o *ormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	log.Error().Msgf(str, args)
}

func (o *ormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	isSlow := false
	if elapsed > o.slowThreshold {
		isSlow = true
		log.Warn().Msgf("SLOW: %t, Elapsed: %.f ms, SQL: %s", isSlow, float64(elapsed.Nanoseconds())/1e6, sql)
	} else {
		log.Trace().Msgf("SLOW: %t, Elapsed: %.f ms, SQL: %s", isSlow, float64(elapsed.Nanoseconds())/1e6, sql)
	}
}

// NewMySQL new MySQL db and retry connection when has error.
func New(c *Config) (db *gorm.DB) {
	var err error
	if c.SlowThreshold == 0 {
		c.SlowThreshold = 200 * time.Millisecond
	}
	ormConf := &gorm.Config{
		Logger: &ormLogger{slowThreshold: c.SlowThreshold},
	}
	if c.DbType == "postgre" {
		db, err = gorm.Open(postgres.Open(c.DSN), ormConf)
	} else {
		db, err = gorm.Open(mysql.Open(c.DSN), ormConf)
	}

	if err != nil {
		log.Error().Msgf("db dsn(%s) error: %v", c.DSN, err)
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Msgf("db.DB() error: %v", err)
		panic(err)
	}

	// sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(c.Idle)
	// sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(c.Active)
	// sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(c.IdleTimeout)

	return
}
