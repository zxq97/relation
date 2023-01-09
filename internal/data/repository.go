package data

import (
	"github.com/patrickmn/go-cache"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/gotool/memcachex"
	"github.com/zxq97/gotool/redisx"
	"upper.io/db.v3/lib/sqlbuilder"
)

//type RelationBFFRepo struct {
//	client relationsvc.RelationSvcClient
//}

type RelationSvcRepoImpl struct {
	redis    *redisx.RedisX
	mc       *memcachex.MemcacheX
	sess     sqlbuilder.Database
	producer *kafka.Producer
}

type RelationTaskRepoImpl struct {
	redis *redisx.RedisX
	mc    *memcachex.MemcacheX
	sess  sqlbuilder.Database
	cache *cache.Cache
}
