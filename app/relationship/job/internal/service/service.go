package service

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/google/wire"
	"github.com/zxq97/gokit/pkg/mq"
	"github.com/zxq97/relation/app/relationship/job/internal/biz"
	"github.com/zxq97/relation/app/relationship/pkg/message"
)

var ProviderSet = wire.NewSet(NewRelationshipJob)

type RelationshipJob struct {
	uc *biz.RelationshipUseCase
}

func NewRelationshipJob(uc *biz.RelationshipUseCase) *RelationshipJob {
	return &RelationshipJob{uc: uc}
}

func (s *RelationshipJob) Relation(ctx context.Context, msg *mq.MqMessage) error {
	follow := &message.AsyncFollow{}
	err := proto.Unmarshal(msg.Message, follow)
	if err != nil {
		return err
	}
	switch msg.Tag {
	case message.TagCreate:
		err = s.uc.Follow(ctx, follow.Uid, follow.ToUid)
	case message.TagDelete:
		err = s.uc.Unfollow(ctx, follow.Uid, follow.ToUid)
	}
	return err
}

func (s *RelationshipJob) Rebuild(ctx context.Context, msg *mq.MqMessage) error {
	list := &message.CacheRebuild{}
	err := proto.Unmarshal(msg.Message, list)
	if err != nil {
		return err
	}
	switch msg.Tag {
	case message.TagListMissed:
		err = s.uc.FollowerCacheRebuild(ctx, list.Uid, list.LastId)
	}
	return err
}
