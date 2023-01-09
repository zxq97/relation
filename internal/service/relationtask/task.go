package relationtask

import (
	"github.com/zxq97/gotool/config"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/relation/internal/biz"
	"github.com/zxq97/relation/internal/env"
)

var (
	consumers = []*kafka.Consumer{}
)

type RelationTask struct {
	biz *biz.RelationTaskBIZ
}

func InitRelationTask(conf *RelationTaskConfig) (*RelationTask, error) {
	err := env.InitLog(conf.LogPath)
	if err != nil {
		return nil, err
	}
	rtb, err := biz.NewRelationTaskBIZ(conf.Redis["redis"], conf.MC["mc"], conf.Mysql["relation"])
	if err != nil {
		return nil, err
	}
	return &RelationTask{
		biz: rtb,
	}, nil
}

func InitConsumer(conf *config.KafkaConf, task *RelationTask) error {
	relationConsumer, err := kafka.InitConsumer(conf.Addr, []string{kafka.TopicRelationFollow}, "relation_task_follow", task.relation, env.ApiLogger, env.ExcLogger)
	if err != nil {
		return err
	}
	rebuildConsumer, err := kafka.InitConsumer(conf.Addr, []string{kafka.TopicRelationCacheRebuild}, "relation_task_rebuild", task.rebuild, env.ApiLogger, env.ExcLogger)
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
