package cache

import (
	"github.com/zxq97/gotool/config"
	"github.com/zxq97/gotool/memcachex"
	"github.com/zxq97/gotool/redisx"
)

var (
	mcx *memcachex.MemcacheX
	rdx *redisx.RedisX
)

func InitCache(redisConf *config.RedisConf, mcConf *config.MCConf) {
	rdx = redisx.NewRedisX(redisConf)
	mcx = memcachex.NewMemcacheX(mcConf)
}
