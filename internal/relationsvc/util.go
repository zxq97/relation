package relationsvc

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/zxq97/gotool/concurrent"
	"github.com/zxq97/gotool/constant"
	"github.com/zxq97/gotool/generate"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/relation/internal/cache"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/model"
	"github.com/zxq97/relation/internal/store"
)

func packKafkaMsg(ctx context.Context, req proto.Message, eventtype int32) ([]byte, error) {
	trace, ok := ctx.Value(constant.TraceIDKey).(string)
	if !ok {
		trace = generate.UUIDStr()
	}
	bs, err := proto.Marshal(req)
	if err != nil {
		env.ExcLogger.Printf("ctx %v follow Marshal req %#v err %v", ctx, req, err)
		return nil, err
	}
	kfkmsg := &kafka.KafkaMessage{
		TraceId:   trace,
		EventType: eventtype,
		Message:   bs,
	}
	bs, err = proto.Marshal(kfkmsg)
	if err != nil {
		env.ExcLogger.Printf("ctx %v follow Marshal kfkmsg %#v err %v", ctx, kfkmsg, err)
	}
	return bs, err
}

func getUserFollow(ctx context.Context, uid int64) ([]*model.FollowItem, error) {
	val, err, _ := sfg.Do(fmt.Sprintf(sfKeyGetFollowList, uid), func() (interface{}, error) {
		list, err := store.GetAllUserFollow(ctx, uid)
		if err != nil {
			env.ExcLogger.Println()
			return nil, err
		}
		concurrent.Go(func() {
			cache.SetFollowList(context.TODO(), uid, list)
		})
		return list, nil
	})
	list, ok := val.([]*model.FollowItem)
	if err != nil || !ok {
		return nil, err
	}
	return list, nil
}

func fcDAO2DTO(m map[int64]*model.UserFollowCount) map[int64]*RelationCount {
	rm := make(map[int64]*RelationCount, len(m))
	for k, v := range m {
		rm[k] = &RelationCount{
			FollowCount:   int64(v.FollowCount),
			FollowerCount: int64(v.FollowerCount),
		}
	}
	return rm
}
