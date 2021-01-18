package setting

import (
	"time"

	"github.com/go-ini/ini"
	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/blog/pkg/cache/gredis"
	"github.com/shipengqi/example.v1/blog/pkg/database/orm"
	"github.com/shipengqi/example.v1/blog/pkg/logger"
)

var settings = New()

type App struct {
	SingingKey   string `ini:"SIGNING_KEY"`
	PageSize     int    `ini:"PAGE_SIZE"`
	IsPrintStack bool   `ini:"IS_PRINT_STACK"`
}

type Server struct {
	HttpPort     int           `ini:"HTTP_PORT"`
	HttpsPort    int           `ini:"HTTPS_PORT"`
	ReadTimeout  time.Duration `ini:"READ_TIMEOUT"`
	WriteTimeout time.Duration `ini:"WRITE_TIMEOUT"`
}

type Setting struct {
	RunMode string         `ini:"RUN_MODE"`
	App     *App           `ini:"app"`
	Server  *Server        `ini:"server"`
	DB      *orm.Config    `ini:"database"`
	Redis   *gredis.Config `ini:"redis"`
	Log     *logger.Config `ini:"log"`
}

func New() *Setting {
	return &Setting{
		RunMode: "",
		App:     &App{},
		Server:  &Server{},
		DB:      &orm.Config{},
		Redis:   &gredis.Config{},
		Log:     &logger.Config{},
	}
}

func Settings() *Setting {
	return settings
}

func ServerSettings() *Server {
	return settings.Server
}

func AppSettings() *App {
	return settings.App
}

func LogSettings() *logger.Config {
	return settings.Log
}

func DBSettings() *orm.Config {
	return settings.DB
}

func RedisSettings() *gredis.Config {
	return settings.Redis
}

func Init(filename string) (*Setting, error) {
	cfg, err := ini.Load(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "setting.Init, fail to parse '%s'", filename)
	}
	cfg.BlockMode = false

	err = cfg.MapTo(settings)
	if err != nil {
		return nil, err
	}

	settings.Server.ReadTimeout = settings.Server.ReadTimeout * time.Second
	settings.Server.WriteTimeout = settings.Server.WriteTimeout * time.Second
	settings.DB.IdleTimeout = settings.DB.IdleTimeout * time.Second
	settings.DB.SlowThreshold = settings.DB.IdleTimeout * time.Millisecond
	settings.Redis.IdleTimeout = settings.Redis.IdleTimeout * time.Second

	return settings, nil
}
