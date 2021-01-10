package setting

import (
	"log"

	"github.com/go-ini/ini"
	"github.com/pkg/errors"

	"github.com/shipengqi/example.v1/blog/pkg/database/gredis"
	"github.com/shipengqi/example.v1/blog/pkg/database/orm"
	"github.com/shipengqi/example.v1/blog/pkg/logger"
)

var settings = &Setting{}

type Setting struct {
	RunMode      string
	App          *app
	Server       *server
	DB           *orm.Config
	Redis        *gredis.Config
	Log          *logger.Config
}

func Settings() *Setting {
	return settings
}

func Init() (*Setting, error) {
	var err error
	_, err = ini.Load("conf/app.debug.ini")
	if err != nil {
		return nil, errors.Wrap(err, "setting.Init, fail to parse 'conf/app.debug.ini'")
	}

	return settings, nil
}

// mapTo map section
func mapTo(cfg *ini.File, section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}