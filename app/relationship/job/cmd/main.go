package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/zxq97/gokit/pkg/cache/xmemcache"
	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/gokit/pkg/config"
	"github.com/zxq97/gokit/pkg/database/xmysql"
	"github.com/zxq97/relation/app/relationship/job/internal/biz"
	"github.com/zxq97/relation/app/relationship/job/internal/conf"
	"github.com/zxq97/relation/app/relationship/job/internal/data"
	"github.com/zxq97/relation/app/relationship/job/internal/service"
)

var (
	flagConf string
	appConf  conf.Config
)

func init() {
	flag.StringVar(&flagConf, "conf", "app/relationship/job/configs/relationship_job.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	err := config.LoadYaml(flagConf, &appConf)
	if err != nil {
		panic(err)
	}
	dbCli, err := xmysql.NewMysqlDB(appConf.Mysql)
	if err != nil {
		panic(err)
	}
	redisCli := xredis.NewRedis(appConf.Redis)
	memcacheCli := xmemcache.NewMemcache(appConf.Memcache)
	repo := data.NewRelationshipRepo(memcacheCli, redisCli, dbCli)
	userCase := biz.NewRelationshipUseCase(repo)
	server := service.NewRelationshipJob(userCase)
	err = server.AddConsumer(appConf.Kafka)
	if err != nil {
		panic(err)
	}
	server.Start()
	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		errCh <- http.ListenAndServe(appConf.Server.HttpBind, nil)
	}()

	select {
	case err = <-errCh:
		server.Stop()
		log.Println("relationship job stop err", err)
	case sig := <-sigCh:
		server.Stop()
		log.Println("relationship job stop sign", sig)
	}
}
