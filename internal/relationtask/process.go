package relationtask

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/relation/internal/cache"
	"github.com/zxq97/relation/internal/constant"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/model"
	"github.com/zxq97/relation/internal/relationsvc"
	"github.com/zxq97/relation/internal/store"
)

func relation(ctx context.Context, kfkmsg *kafka.KafkaMessage) {
	follow := &relationsvc.FollowRequest{}
	err := proto.Unmarshal(kfkmsg.Message, follow)
	if err != nil {
		env.ExcLogger.Println()
		return
	}
	switch kfkmsg.EventType {
	case kafka.EventTypeCreate:
		err = store.Follow(ctx, follow.Uid, follow.ToUid)
		if err != nil {
			env.ExcLogger.Println()
			return
		}
		cache.AddRelation(ctx, follow.Uid, &model.FollowItem{ToUid: follow.ToUid, CreateTime: time.Now().UnixMilli()})
	case kafka.EventTypeDelete:
		err = store.Unfollow(ctx, follow.Uid, follow.ToUid)
		if err != nil {
			env.ExcLogger.Println()
			return
		}
		cache.DelRelation(ctx, follow.Uid, follow.ToUid)
	}
}

func rebuild(ctx context.Context, kfkmsg *kafka.KafkaMessage) {
	list := &relationsvc.ListRequest{}
	err := proto.Unmarshal(kfkmsg.Message, list)
	if err != nil {
		env.ExcLogger.Println()
		return
	}
	switch kfkmsg.EventType {
	case kafka.EventTypeListMissed:
		item, err := store.GetUserFollower(ctx, list.Uid, list.LastId, constant.ListRebuildSize)
		if err != nil {
			env.ExcLogger.Println()
			return
		}
		cache.AddUserFollower(ctx, list.Uid, item)
	}
}
