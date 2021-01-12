package gredis

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Host        string        `ini:"HOST"`
	Password    string        `ini:"PASSWORD"`
	MaxIdle     int           `ini:"MAX_IDLE"`
	MaxActive   int           `ini:"MAX_ACTIVE"`
	IdleTimeout time.Duration `ini:"IDLE_TIMEOUT"`
}

var redisConn *redis.Pool

func New(c *Config) *redis.Pool {
	redisConn = &redis.Pool{
		MaxIdle:         c.MaxIdle,
		MaxActive:       c.MaxActive,
		IdleTimeout:     c.IdleTimeout,
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
		// health check func
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return redisConn
}

func SetString(key string, data string, exp int) error {
	conn := redisConn.Get()
	defer conn.Close()


	_, err := conn.Do("SET", key, data, "EX", exp)
	if err != nil {
		return err
	}

	return nil
}

func Set(key string, data interface{}, exp int) error {
	conn := redisConn.Get()
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

func Exists(key string) bool {
	conn := redisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func Get(key string) ([]byte, error) {
	conn := redisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := redisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := redisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
