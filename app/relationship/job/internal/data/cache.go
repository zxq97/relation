package data

import "time"

const (
	lcKeyFollower = "loc_fo_%d" // uid
)

func (r *relationshipRepo) localCacheGet(key string) (interface{}, bool) {
	return r.cache.Get(key)
}

func (r *relationshipRepo) localCacheSet(key string, val interface{}, d time.Duration) {
	r.cache.Set(key, val, d)
}
