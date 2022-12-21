package relationtask

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/relation/internal/cache"
	"github.com/zxq97/relation/internal/constant"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/model"
	"github.com/zxq97/relation/internal/service/relationsvc"
	"github.com/zxq97/relation/internal/store"
)

const (
	localCacheKeyFollower = "loc_fo_%d_%d" // uid lastid
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
		key := fmt.Sprintf(localCacheKeyFollower, list.Uid, list.LastId)
		// xxx 运用 local cache 防止上层多个cache missed 多次穿透db
		_, ok := localCache.Get(key)
		if !ok {
			item, err := store.GetUserFollower(ctx, list.Uid, list.LastId, constant.ListRebuildSize)
			if err != nil {
				env.ExcLogger.Println()
				return
			}
			err = cache.AddUserFollower(ctx, list.Uid, item)
			if err != nil {
				env.ExcLogger.Println()
				return
			}
			localCache.Set(key, nil, time.Minute*5)
		}
	}
}
