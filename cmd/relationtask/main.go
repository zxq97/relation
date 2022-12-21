package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zxq97/gotool/config"
	relationtask2 "github.com/zxq97/relation/internal/service/relationtask"
)

var (
	confPath = flag.String("conf", "", "configuration file")
	conf     relationtask2.RelationTaskConfig
)

func main() {
	flag.Parse()
	err := config.LoadYaml(*confPath, &conf)
	if err != nil {
		panic(err)
	}
	conf.Initialize()

	err = relationtask2.InitRelationTask(&conf)
	if err != nil {
		panic(err)
	}
	err = relationtask2.InitConsumer(conf.Kafka["kafka"])
	if err != nil {
		panic(err)
	}

	http.Handle("/metrics", promhttp.Handler())

	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		errCh <- http.ListenAndServe(conf.HTTPBind, nil)
	}()

	select {
	case err = <-errCh:
		relationtask2.StopConsumer()
		log.Println("relationtask stop err", errCh)
	case sign := <-sigCh:
		relationtask2.StopConsumer()
		log.Println("relationtask stop sign", sign)
	}
}
