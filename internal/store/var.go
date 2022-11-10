package store

import (
	"github.com/zxq97/gotool/config"
	"upper.io/db.v3/lib/sqlbuilder"
)

var (
	sess sqlbuilder.Database
)

func InitStore(conf *config.MysqlConf) error {
	var err error
	sess, err = conf.InitDB()
	return err
}
