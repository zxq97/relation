package relationtask

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/service/relationsvc"
)

func (task *RelationTask) relation(ctx context.Context, kfkmsg *kafka.KafkaMessage) {
	follow := &relationsvc.FollowRequest{}
	err := proto.Unmarshal(kfkmsg.Message, follow)
	if err != nil {
		env.ExcLogger.Println()
		return
	}
	switch kfkmsg.EventType {
	case kafka.EventTypeCreate:
		err = task.biz.Follow(ctx, follow.Uid, follow.ToUid)
	case kafka.EventTypeDelete:
		err = task.biz.Unfollow(ctx, follow.Uid, follow.ToUid)
	}
	if err != nil {
		env.ExcLogger.Println()
	}
}

func (task *RelationTask) rebuild(ctx context.Context, kfkmsg *kafka.KafkaMessage) {
	list := &relationsvc.ListRequest{}
	err := proto.Unmarshal(kfkmsg.Message, list)
	if err != nil {
		env.ExcLogger.Println()
		return
	}
	switch kfkmsg.EventType {
	case kafka.EventTypeListMissed:
		err = task.biz.FollowerCacheRebuild(ctx, list.Uid, list.LastId)
	}
	if err != nil {
		env.ExcLogger.Println()
	}
}
