package setting

import (
	"testing"

	"github.com/shipengqi/example.v1/blog/pkg/database/gredis"
	"github.com/shipengqi/example.v1/blog/pkg/database/orm"
	"github.com/shipengqi/example.v1/blog/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestSettings(t *testing.T) {
	s := Settings()
	expected := New()
	assert.Equal(t, expected, s)
}

func TestInit(t *testing.T) {
	s, err := Init("C:\\code\\example.v1\\blog\\conf\\app.debug.ini")
	if err != nil {
		t.Fatalf("An error occurs when setting.Init: %#v", err)
	}
	expected := New()
	expected.RunMode = "debug"
	expected.Server = &Server{
		HttpPort:     8080,
		HttpsPort:    8443,
		ReadTimeout:  60000000000,
		WriteTimeout: 60000000000,
	}
	expected.App = &App{
		SingingKey:   "3633423$412342199",
		PageSize:     10,
		IsPrintStack: true,
	}
	expected.Log = &logger.Config{
		Level:  "trace",
		Output: "/var/run/example.v1",
		Prefix: "example.v1",
	}
	expected.DB = &orm.Config{
		DbType:      "mysql",
		User:        "root",
		Password:    "123456",
		Host:        "16.155.194.49:33061",
		Name:        "blog",
		TablePrefix: "blog_",
	}
	expected.Redis = &gredis.Config{
		Host:        "16.155.194.49:6379",
		Password:    "",
		MaxIdle:     30,
		MaxActive:   30,
		IdleTimeout: 200000000000,
	}
	assert.Equal(t, expected, s)
}
