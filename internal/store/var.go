package store

import (
	"github.com/zxq97/gotool/config"
	"upper.io/db.v3/lib/sqlbuilder"
)

var (
	dbCli sqlbuilder.Database
)

func InitStore(conf *config.MysqlConf) error {
	var err error
	dbCli, err = conf.InitDB()
	return err
}
