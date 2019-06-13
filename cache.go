package models

import (
	"github.com/gocommon/cache"
	"github.com/gocommon/models/errors"
)

// Caches Caches
var Caches = map[string]cache.Cacher{}

// InitCache InitCache
func InitCache(confs map[string]cache.Options) {
	for name := range confs {
		Caches[name] = cache.New(confs[name])
	}
}

// Cache Cache
func Cache(name ...string) cache.Cacher {
	cname := "default"
	if len(name) > 0 {
		cname = name[0]
	}

	c, has := Caches[cname]
	if !has {
		return cache.NewErrCacher(errors.New("cache not found %s", name))
	}

	return c
}
