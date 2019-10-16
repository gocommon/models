package models

import (
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gocommon/models/errors"
)

// RedisService RedisService
type RedisService struct {
	Enable bool   `dsn:"query.enable"`
	Addr   string `dsn:"address"`
	Passwd string `dsn:"password"`

	MaxIdle     int `dsn:"query.maxidle"`
	IdleTimeout int `dsn:"query.idletimeout"`
}

// SRedisPools SRedisPools
var SRedisPools = &SRedisPool{}

// SRedisPool SRedisPool
type SRedisPool struct {
	rw    sync.RWMutex
	pools map[string]*redis.Pool
}

// Get Get
func (s *SRedisPool) Get(name string) *redis.Pool {
	s.rw.RLock()
	defer s.rw.RUnlock()

	k := "default"
	if len(name) > 0 {
		k = name
	}

	if pool, ok := s.pools[k]; ok {
		return pool

	}

	return nil
}

// Reload Reload
func (s *SRedisPool) Reload(confs map[string]RedisService) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	for _, v := range s.pools {
		go func(pool *redis.Pool) {
			pool.Close()
		}(v)
	}

	pools := make(map[string]*redis.Pool)
	for name, conf := range confs {
		pools[name] = newRedis(conf)
	}

	s.pools = pools
}

// HasRedis HasRedis
func (s *SRedisPool) HasRedis(name string) bool {
	_, ok := s.pools[name]

	return ok
}

// RedisPools RedisPools
// var RedisPools map[string]*redis.Pool

// Redis Redis
func Redis(name ...string) *redis.Pool {
	k := "default"
	if len(name) > 0 {
		k = name[0]
	}

	if pool := SRedisPools.Get(k); pool != nil {
		return pool
	}

	// if pool, ok := RedisPools[k]; ok {
	// 	return pool

	// }

	panic(errors.New("unkonw redis %s", k))

}

// HasRedis HasRedis
func HasRedis(name string) bool {
	// _, ok := RedisPools[name]
	// return ok
	return SRedisPools.HasRedis(name)
}

// InitRedis InitRedis
func InitRedis(confs map[string]RedisService) {
	SRedisPools.Reload(confs)

	// RedisPools = make(map[string]*redis.Pool)
	// for name, conf := range confs {
	// 	RedisPools[name] = newRedis(conf)
	// }
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
