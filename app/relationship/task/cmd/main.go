package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zxq97/gokit/pkg/cache/xmemcache"
	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/gokit/pkg/config"
	"github.com/zxq97/gokit/pkg/database/xmysql"
	"github.com/zxq97/gokit/pkg/mq"
	"github.com/zxq97/gokit/pkg/mq/kafka"
	server2 "github.com/zxq97/gokit/pkg/server"
	"github.com/zxq97/gokit/pkg/server/consumer"
	"github.com/zxq97/relation/app/relationship/pkg/dal/cache"
	"github.com/zxq97/relation/app/relationship/pkg/dal/query"
	"github.com/zxq97/relation/app/relationship/pkg/message"
	"github.com/zxq97/relation/app/relationship/task/internal/biz"
	"github.com/zxq97/relation/app/relationship/task/internal/conf"
	"github.com/zxq97/relation/app/relationship/task/internal/data"
	"github.com/zxq97/relation/app/relationship/task/internal/service"
)

var (
	flagConf string
	appConf  conf.Config
)

func init() {
	flag.StringVar(&flagConf, "conf", "app/relationship/task/configs/relationship_task.yaml", "config path, eg: -conf config.yaml")
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
	redisCli := xredis.NewXRedis(appConf.Redis)
	memcacheCli := xmemcache.NewXMemcache(appConf.Memcache)
	repo := data.NewRelationshipRepo(cache.Use(redisCli, memcacheCli), query.Use(dbCli))
	useCase := biz.NewRelationshipUseCase(repo)
	server := service.NewRelationshipTask(useCase)
	syncConsumer, err := kafka.NewConsumer(appConf.Kafka, []string{message.TopicRelationSyncCount}, "relationship_task_sync", server.SyncRecordByUID, mq.WithProcTimeout(time.Second*3))
	if err != nil {
		panic(err)
	}
	s, err := consumer.NewServer([]mq.Consumer{syncConsumer}, server2.WithStartTimeout(time.Second), server2.WithStopTimeout(time.Second))
	if err != nil {
		panic(err)
	}
	c := cron.New()
	if _, err = c.AddFunc("*/5 * * * *", func() {
		if err = server.CronTaskSyncRecord(context.Background(), 10000); err != nil {
			log.Println(err)
		}
	}); err != nil {
		panic(err)
	}
	if err = s.Start(context.Background()); err != nil {
		panic(err)
	}
	c.Start()

	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		errCh <- http.ListenAndServe(appConf.Server.HttpBind, nil)
	}()

	select {
	case err = <-errCh:
		c.Stop()
		serr := s.Stop(context.Background())
		log.Println("relationship job stop err", err, serr)
	case sig := <-sigCh:
		c.Stop()
		serr := s.Stop(context.Background())
		log.Println("relationship job stop sign", sig, serr)
	}
}
