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

	"github.com/zxq97/gokit/pkg/cache/xmemcache"
	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/gokit/pkg/config"
	"github.com/zxq97/gokit/pkg/database/xmysql"
	"github.com/zxq97/gokit/pkg/mq"
	"github.com/zxq97/gokit/pkg/mq/kafka"
	server2 "github.com/zxq97/gokit/pkg/server"
	"github.com/zxq97/gokit/pkg/server/consumer"
	"github.com/zxq97/relation/app/relationship/job/internal/biz"
	"github.com/zxq97/relation/app/relationship/job/internal/conf"
	"github.com/zxq97/relation/app/relationship/job/internal/data"
	"github.com/zxq97/relation/app/relationship/job/internal/service"
	"github.com/zxq97/relation/app/relationship/pkg/dal/cache"
	"github.com/zxq97/relation/app/relationship/pkg/dal/query"
	"github.com/zxq97/relation/app/relationship/pkg/message"
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
	producer, err := kafka.NewProducer(appConf.Kafka)
	if err != nil {
		panic(err)
	}
	dbCli, err := xmysql.NewMysqlDB(appConf.Mysql)
	if err != nil {
		panic(err)
	}
	redisCli := xredis.NewXRedis(appConf.Redis)
	memcacheCli := xmemcache.NewXMemcache(appConf.Memcache)
	repo := data.NewRelationshipRepo(producer, cache.Use(redisCli, memcacheCli), query.Use(dbCli))
	useCase := biz.NewRelationshipUseCase(repo)
	server := service.NewRelationshipJob(useCase)
	followConsumer, err := kafka.NewConsumer(appConf.Kafka, []string{message.TopicRelationFollow}, "relationship_job_follow", server.Relation, mq.WithProcTimeout(time.Second*3))
	if err != nil {
		panic(err)
	}
	rebuildConsumer, err := kafka.NewConsumer(appConf.Kafka, []string{message.TopicRelationCacheRebuild}, "relationship_job_rebuild", server.Rebuild, mq.WithProcTimeout(time.Second*3))
	if err != nil {
		panic(err)
	}
	s, err := consumer.NewServer([]mq.Consumer{followConsumer, rebuildConsumer}, server2.WithStartTimeout(time.Second), server2.WithStopTimeout(time.Second))
	if err != nil {
		panic(err)
	}
	if err = s.Start(context.Background()); err != nil {
		panic(err)
	}

	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		errCh <- http.ListenAndServe(appConf.Server.HttpBind, nil)
	}()

	select {
	case err = <-errCh:
		serr := s.Stop(context.Background())
		log.Println("relationship job stop err", err, serr)
	case sig := <-sigCh:
		serr := s.Stop(context.Background())
		log.Println("relationship job stop sign", sig, serr)
	}
}
