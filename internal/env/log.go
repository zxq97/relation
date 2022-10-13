package env

import (
	"log"

	"github.com/zxq97/gotool/config"
)

var (
	ApiLogger   *log.Logger
	ExcLogger   *log.Logger
	DebugLogger *log.Logger
)

func InitLog(conf *config.LogConf) error {
	var err error
	ApiLogger, err = config.InitLog(conf.Api)
	if err != nil {
		return err
	}
	ExcLogger, err = config.InitLog(conf.Exc)
	if err != nil {
		return err
	}
	DebugLogger, err = config.InitLog(conf.Debug)
	return err
}
