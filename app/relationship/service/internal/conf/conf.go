package conf

import (
	"github.com/zxq97/gokit/pkg/cache/xmemcache"
	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/gokit/pkg/database/xmysql"
	"github.com/zxq97/gokit/pkg/etcd"
	"github.com/zxq97/gokit/pkg/mq/kafka"
	"github.com/zxq97/gokit/pkg/rpc"
)

type Config struct {
	Server   *rpc.Config       `yaml:"server"`
	Etcd     *etcd.Config      `yaml:"etcd"`
	Redis    *xredis.Config    `yaml:"redis"`
	Mysql    *xmysql.Config    `yaml:"mysql"`
	Memcache *xmemcache.Config `yaml:"memcache"`
	Kafka    *kafka.Config     `yaml:"kafka"`
}
