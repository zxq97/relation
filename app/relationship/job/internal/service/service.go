package service

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/zxq97/gokit/pkg/mq"
	"github.com/zxq97/relation/api/relationship/service/v1"
	"github.com/zxq97/relation/app/relationship/job/internal/biz"
)

type RelationshipJob struct {
	uc *biz.RelationshipUseCase
}

func NewRelationshipJob(uc *biz.RelationshipUseCase) *RelationshipJob {
	return &RelationshipJob{uc: uc}
}

func (s *RelationshipJob) Relation(ctx context.Context, msg *mq.MqMessage) error {
	follow := &v1.FollowRequest{}
	err := proto.Unmarshal(msg.Message, follow)
	if err != nil {
		return err
	}
	switch msg.Tag {
	case mq.TagCreate:
		err = s.uc.Follow(ctx, follow.Uid, follow.ToUid)
	case mq.TagDelete:
		err = s.uc.Unfollow(ctx, follow.Uid, follow.ToUid)
	}
	return err
}

func (s *RelationshipJob) Rebuild(ctx context.Context, msg *mq.MqMessage) error {
	list := &v1.ListRequest{}
	err := proto.Unmarshal(msg.Message, list)
	if err != nil {
		return err
	}
	switch msg.Tag {
	case mq.TagListMissed:
		err = s.uc.FollowerCacheRebuild(ctx, list.Uid, list.LastId)
	}
	return err
}
