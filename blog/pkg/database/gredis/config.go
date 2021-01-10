package gredis

import "time"

type Config struct {
	host        string
	password    string
	maxIdle     int
	maxActive   int
	idleTimeout time.Duration
}

func (c *Config) Host() string {
	return c.host
}

func (c *Config) SetHost(host string) {
	c.host = host
}

func (c *Config) Password() string {
	return c.password
}

func (c *Config) SetPassword(password string) {
	c.password = password
}

func (c *Config) MaxIdle() int {
	return c.maxIdle
}

func (c *Config) SetMaxIdle(maxIdle int) {
	c.maxIdle = maxIdle
}

func (c *Config) MaxActive() int {
	return c.maxActive
}

func (c *Config) SetMaxActive(maxActive int) {
	c.maxActive = maxActive
}

func (c *Config) IdleTimeout() time.Duration {
	return c.idleTimeout
}

func (c *Config) SetIdleTimeout(idleTimeout time.Duration) {
	c.idleTimeout = idleTimeout
}