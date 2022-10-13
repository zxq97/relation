package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/zxq97/gotool/config"
	"github.com/zxq97/gotool/rpc"
	"github.com/zxq97/relation/internal/relationsvc"
	"google.golang.org/grpc/reflection"
)

var (
	confPath = flag.String("conf", "", "configuration file")
	conf     relationsvc.RelationSvcConfig
)

func main() {
	flag.Parse()
	err := config.LoadYaml(*confPath, &conf)
	if err != nil {
		panic(err)
	}
	conf.Initialize()
	err = relationsvc.InitRelationSvc(&conf)
	if err != nil {
		panic(err)
	}
	etcdCli, err := conf.Etcd["etcd"].InitEtcd()
	if err != nil {
		panic(err)
	}
	svc, er := rpc.NewGrpcServer(etcdCli, conf.Svc.Name+"_"+conf.Svc.Bind, conf.Svc.Bind)
	_, err = er.Register()
	if err != nil {
		panic(err)
	}
	relationsvc.RegisterRelationSvcServer(svc, relationsvc.RelationSvc{})

	lis, err := net.Listen("tcp", conf.Svc.Bind)
	if err != nil {
		panic(err)
	}

	reflection.Register(svc)
	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		errCh <- svc.Serve(lis)
	}()
	go func() {
		errCh <- http.ListenAndServe(conf.Svc.HttpBind, nil)
	}()

	select {
	case err = <-errCh:
		er.Stop()
		log.Println("relationsvc stop err", err)
	case sig := <-sigCh:
		er.Stop()
		log.Println("relationsvc stop sign", sig)
	}
}
