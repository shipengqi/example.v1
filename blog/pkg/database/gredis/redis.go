package gredis

import "time"

type Config struct {
	Host        string        `ini:"HOST"`
	Password    string        `ini:"PASSWORD"`
	MaxIdle     int           `ini:"MAX_IDLE"`
	MaxActive   int           `ini:"MAX_ACTIVE"`
	IdleTimeout time.Duration `ini:"IDLE_TIMEOUT"`
}

func New() {

}
