package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/zxq97/gokit/pkg/cache/xmemcache"
	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/gokit/pkg/config"
	"github.com/zxq97/gokit/pkg/database/xmysql"
	"github.com/zxq97/gokit/pkg/etcd"
	"github.com/zxq97/gokit/pkg/mq/kafka"
	"github.com/zxq97/gokit/pkg/rpc"
	"github.com/zxq97/relation/api/relationship/service/v1"
	"github.com/zxq97/relation/app/relationship/service/internal/biz"
	"github.com/zxq97/relation/app/relationship/service/internal/conf"
	"github.com/zxq97/relation/app/relationship/service/internal/data"
	"github.com/zxq97/relation/app/relationship/service/internal/service"
)

var (
	flagConf string
	appConf  conf.Config
)

func init() {
	flag.StringVar(&flagConf, "conf", "app/relationship/service/configs/relationship_service.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	err := config.LoadYaml(flagConf, &appConf)
	if err != nil {
		panic(err)
	}
	etcdCli, err := etcd.NewEtcd(appConf.Etcd)
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc, err := rpc.NewGrpcServer(ctx, appConf.Server, etcdCli)
	if err != nil {
		panic(err)
	}
	repo := data.NewRelationshipRepo(producer, memcacheCli, redisCli, dbCli)
	useCase := biz.NewRelationshipUseCase(repo)
	server := service.NewRelationshipService(useCase)
	v1.RegisterRelationSvcServer(svc, server)
	lis, err := net.Listen("tcp", appConf.Server.Bind)
	if err != nil {
		panic(err)
	}
	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		errCh <- svc.Serve(lis)
	}()
	go func() {
		errCh <- http.ListenAndServe(appConf.Server.HttpBind, nil)
	}()

	select {
	case err = <-errCh:
		cancel()
		log.Println("relationship service stop err", err)
	case sig := <-sigCh:
		cancel()
		log.Println("relationship service stop sign", sig)
	}
}
