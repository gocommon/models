package models

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gocommon/models/errors"
)

// RedisService RedisService
type RedisService struct {
	Enable bool
	Addr   string
	Passwd string

	MaxIdle     int
	IdleTimeout int
}

// RedisPools RedisPools
var RedisPools map[string]*redis.Pool

// Redis Redis
func Redis(name ...string) *redis.Pool {
	k := "default"
	if len(name) > 0 {
		k = name[0]
	}

	if pool, ok := RedisPools[k]; ok {
		return pool

	}

	panic(errors.New("unkonw redis %s", k))

	return nil
}

// HasRedis HasRedis
func HasRedis(name string) bool {
	_, ok := RedisPools[name]
	return ok
}

// InitRedis InitRedis
func InitRedis(confs map[string]RedisService) {
	RedisPools = make(map[string]*redis.Pool)
	for name, conf := range confs {
		RedisPools[name] = newRedis(conf)
	}
}

func newRedis(conf RedisService) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Addr)
			if err != nil {
				return nil, errors.Wrap(err, "redis Dial")
			}
			if len(conf.Passwd) > 0 {
				if _, err := c.Do("AUTH", conf.Passwd); err != nil {
					c.Close()
					return nil, errors.Wrap(err, "redis auth")
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return errors.Wrap(err, "redis ping")
		},
	}
}
