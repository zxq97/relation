package relationtask

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/relationsvc"
	"github.com/zxq97/relation/internal/store"
)

func relation(ctx context.Context, kfkmsg *kafka.KafkaMessage) {
	follow := &relationsvc.FollowRequest{}
	err := proto.Unmarshal(kfkmsg.Message, follow)
	if err != nil {
		env.ExcLogger.Printf("")
		return
	}
	switch kfkmsg.EventType {
	case kafka.EventTypeCreate:
		err = store.Follow(ctx, follow.Uid, follow.ToUid)
		if err != nil {
			env.ExcLogger.Println()
			return
		}

	case kafka.EventTypeDelete:

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

	}
}
