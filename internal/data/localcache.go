package data

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	lcKeyFollower = "loc_fo_%d_%d" // uid lastid
)

func lcGet(lc *cache.Cache, key string) (interface{}, bool) {
	return lc.Get(key)
}

func lcSet(lc *cache.Cache, key string, val interface{}, d time.Duration) {
	lc.Set(key, val, d)
}
