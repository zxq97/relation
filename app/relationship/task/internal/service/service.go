package service

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/google/wire"
	"github.com/zxq97/gokit/pkg/mq"
	"github.com/zxq97/relation/app/relationship/pkg/message"
	"github.com/zxq97/relation/app/relationship/task/internal/biz"
)

var ProviderSet = wire.NewSet(NewRelationshipTask)

type RelationshipTask struct {
	uc *biz.RelationshipUseCase
}

func NewRelationshipTask(uc *biz.RelationshipUseCase) *RelationshipTask {
	return &RelationshipTask{uc: uc}
}

func (s *RelationshipTask) CronTaskSyncRecord(ctx context.Context, limit int64) error {
	return s.uc.CronTaskSyncRecord(ctx, limit)
}

func (s *RelationshipTask) SyncRecordByUID(ctx context.Context, msg *mq.MqMessage) error {
	record := &message.SyncCount{}
	err := proto.Unmarshal(msg.Message, record)
	if err != nil {
		return err
	}
	switch msg.Tag {
	case message.TagSync:
		err = s.uc.SyncRecordByUID(ctx, record.Uid)
	}
	return err
}
