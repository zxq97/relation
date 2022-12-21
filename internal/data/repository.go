package data

import (
	"github.com/patrickmn/go-cache"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/gotool/memcachex"
	"github.com/zxq97/gotool/redisx"
	"github.com/zxq97/relation/internal/service/relationsvc"
	"upper.io/db.v3/lib/sqlbuilder"
)

type relationBFFRepo struct {
	client relationsvc.RelationSvcClient
}

type relationSVCRepo struct {
	redis    *redisx.RedisX
	mc       *memcachex.MemcacheX
	sess     sqlbuilder.Database
	producer *kafka.Producer
}

type relationTaskRepo struct {
	redis *redisx.RedisX
	mc    *memcachex.MemcacheX
	sess  sqlbuilder.Database
	cache *cache.Cache
}
