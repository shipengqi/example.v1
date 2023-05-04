package get_started

import (
	"fmt"

	"go.uber.org/dig"
	"gopkg.in/ini.v1"
)

type Option struct {
	ConfigFile string `short:"c" long:"config" description:"Name of config file."`
}

func NewOption() (*Option, error) {
	return &Option{ConfigFile: "config.ini"}, nil
}

func NewConfig(opt *Option) (*ini.File, error) {
	cfg, err := ini.Load(opt.ConfigFile)
	return cfg, err
}

func BuildContainer() *dig.Container {
	container := dig.New()

	_ = container.Provide(NewConfig)
	_ = container.Provide(NewOption)

	return container
}

func ExampleDigContainer() {
	c := BuildContainer()
	_ = c.Invoke(func(cfg *ini.File) {
		fmt.Println("App Name:", cfg.Section("").Key("app_name").String())
		fmt.Println("Log Level:", cfg.Section("").Key("log_level").String())
	})

	// Output:
	// App Name: awesome digs
	// Log Level: DEBUG
}
