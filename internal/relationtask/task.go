package relationtask

import (
	"github.com/zxq97/gotool/config"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/relation/internal/cache"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/store"
)

var (
	consumers = []*kafka.Consumer{}
)

func InitCommentTask(conf *RelationTaskConfig) error {
	err := env.InitLog(conf.LogPath)
	if err != nil {
		return err
	}
	cache.InitCache(conf.Redis["redis"], conf.MC["mc"])
	err = store.InitStore(conf.Mysql["relation"])
	return err
}

func InitConsumer(conf *config.KafkaConf) error {
	relationConsumer, err := kafka.InitConsumer(conf.Addr, []string{kafka.TopicRelationFollow}, "relation_task_follow", relation, env.ApiLogger, env.ExcLogger)
	if err != nil {
		return err
	}
	rebuildConsumer, err := kafka.InitConsumer(conf.Addr, []string{kafka.TopicRelationCacheRebuild}, "relation_task_rebuild", rebuild, env.ApiLogger, env.ExcLogger)
	if err != nil {
		return err
	}
	consumers = append(consumers, relationConsumer, rebuildConsumer)
	StartConsumer(consumers)
	return nil
}

func StartConsumer(consumers []*kafka.Consumer) {
	for _, v := range consumers {
		v.Start()
	}
}

func StopConsumer() {
	for _, v := range consumers {
		err := v.Stop()
		if err != nil {
			env.ExcLogger.Println("StopConsumer err", err)
		}
	}
}
