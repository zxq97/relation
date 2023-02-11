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
	"github.com/zxq97/relation/api/black/service/v1"
	"github.com/zxq97/relation/app/black/service/internal/conf"
	"github.com/zxq97/relation/app/black/service/internal/service"
)

var (
	flagConf string
	appConf  conf.Config
)

func init() {
	flag.StringVar(&flagConf, "conf", "app/black/service/configs/black_service.yaml", "config path, eg: -conf config.yaml")
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
	server := service.NewBlackService()
	v1.RegisterBlackSvcServer(svc, server)
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
		log.Println("black service stop err", err)
	case sig := <-sigCh:
		cancel()
		log.Println("black service stop sign", sig)
	}
}
