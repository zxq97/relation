package conf

import (
	"github.com/zxq97/gokit/pkg/etcd"
	"github.com/zxq97/gokit/pkg/rpc"
)

type Config struct {
	Server *rpc.Config  `yaml:"server"`
	Etcd   *etcd.Config `yaml:"etcd"`
}
