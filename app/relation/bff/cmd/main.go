package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/zxq97/gokit/pkg/config"
	"github.com/zxq97/gokit/pkg/etcd"
	"github.com/zxq97/gokit/pkg/rpc"
	blackv1 "github.com/zxq97/relation/api/black/service/v1"
	v1 "github.com/zxq97/relation/api/relation/bff/v1"
	relationshipv1 "github.com/zxq97/relation/api/relationship/service/v1"
	"github.com/zxq97/relation/app/relation/bff/internal/biz"
	"github.com/zxq97/relation/app/relation/bff/internal/conf"
	"github.com/zxq97/relation/app/relation/bff/internal/data"
	"github.com/zxq97/relation/app/relation/bff/internal/service"
)

var (
	flagConf string
	appConf  conf.Config
)

func init() {
	flag.StringVar(&flagConf, "conf", "app/relation/bff/configs/relation_bff.yaml", "config path, eg: -conf config.yaml")
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc, err := rpc.NewGrpcServer(ctx, appConf.Server, etcdCli)
	if err != nil {
		panic(err)
	}
	relationshipConn, err := rpc.NewGrpcConn(ctx, "relationship_service", etcdCli)
	if err != nil {
		panic(err)
	}
	blackConn, err := rpc.NewGrpcConn(ctx, "black_service", etcdCli)
	if err != nil {
		panic(err)
	}
	relationshipClient := relationshipv1.NewRelationSvcClient(relationshipConn)
	blackClient := blackv1.NewBlackSvcClient(blackConn)
	relationRepo := data.NewRelationRepo(relationshipClient, blackClient)
	relationUseCase := biz.NewRelationUseCase(relationRepo)
	blackRepo := data.NewBlackRepo(relationshipClient, blackClient)
	blackUseCase := biz.NewBlackUseCase(blackRepo)
	server := service.NewRelationService(relationUseCase, blackUseCase)
	v1.RegisterRelationBFFServer(svc, server)
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
		log.Println("relation bff stop err", err)
	case sig := <-sigCh:
		cancel()
		log.Println("relation bff stop sign", sig)
	}
}
