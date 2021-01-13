package gredis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Config struct {
	Host        string        `ini:"HOST"`
	Password    string        `ini:"PASSWORD"`
	MaxIdle     int           `ini:"MAX_IDLE"`
	MaxActive   int           `ini:"MAX_ACTIVE"`
	IdleTimeout time.Duration `ini:"IDLE_TIMEOUT"`
}

type Pool struct {
	conf  *Config
	coons *redis.Pool
}

func New(c *Config) *Pool {
	pool := &redis.Pool{
		MaxIdle:     c.MaxIdle,
		MaxActive:   c.MaxActive,
		IdleTimeout: c.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", c.Host)
			if err != nil {
				return nil, err
			}
			if c.Password != "" {
				if _, err := conn.Do("AUTH", c.Password); err != nil {
					_ = conn.Close()
					return nil, err
				}
			}
			return conn, nil
		},
		// health check func, PINGs connections that have been idle more than a minute
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return &Pool{
		conf:  c,
		coons: pool,
	}
}

func (p *Pool) SetString(key string, data string, exp int) error {
	conn := p.coons.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, data, "EX", exp)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pool) Set(key string, data interface{}, exp int) error {
	conn := p.coons.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value, "EX", exp)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pool) Exists(key string) bool {
	conn := p.coons.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func (p *Pool) Get(key string) ([]byte, error) {
	conn := p.coons.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (p *Pool) Delete(key string) (bool, error) {
	conn := p.coons.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func (p *Pool) LikeDeletes(key string) error {
	conn := p.coons.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = p.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Pool) Ping() (err error) {
	conn := p.coons.Get()
	defer conn.Close()

	// _, err = conn.Do("SET", "PING", "PONG")
	_, err = conn.Do("PING")
	return
}

func (p *Pool) Close() (err error) {
	return p.coons.Close()
}
