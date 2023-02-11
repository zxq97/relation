package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/zxq97/gokit/pkg/mq/kafka"
	"github.com/zxq97/relation/api/relationship/service/v1"
	"github.com/zxq97/relation/app/relationship/job/internal/biz"
)

var (
	once sync.Once
)

type RelationshipJob struct {
	uc        *biz.RelationshipUseCase
	consumers []*kafka.Consumer
	done      []<-chan struct{}
}

func NewRelationshipJob(uc *biz.RelationshipUseCase) *RelationshipJob {
	return &RelationshipJob{uc: uc}
}

func (s *RelationshipJob) AddConsumer(conf *kafka.Config) error {
	var err error
	once.Do(func() {
		s.consumers = make([]*kafka.Consumer, 2)
		s.done = make([]<-chan struct{}, 2)
		s.consumers[0], s.done[0], err = kafka.NewConsumer(conf, []string{kafka.TopicRelationFollow}, "relationship_job_follow", "relationship_job", s.relation, 1, 1, time.Second*5)
		if err != nil {
			return
		}
		s.consumers[1], s.done[1], err = kafka.NewConsumer(conf, []string{kafka.TopicRelationCacheRebuild}, "relationship_job_rebuild", "relationship_job", s.rebuild, 1, 1, time.Second*5)
		if err != nil {
			return
		}
	})
	return err
}

func (s *RelationshipJob) Start() {
	for _, v := range s.consumers {
		v.Start()
	}
}

func (s *RelationshipJob) Stop() {
	wg := sync.WaitGroup{}
	for k, v := range s.consumers {
		if err := v.Stop(); err != nil {
			log.Println("consumer stop", err)
		}
		idx := k
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-s.done[idx]
		}()
	}
	wg.Wait()
}

func (s *RelationshipJob) relation(ctx context.Context, msg *kafka.KafkaMessage) error {
	follow := &v1.FollowRequest{}
	err := proto.Unmarshal(msg.Message, follow)
	if err != nil {
		return err
	}
	switch msg.EventType {
	case kafka.EventTypeCreate:
		err = s.uc.Follow(ctx, follow.Uid, follow.ToUid)
	case kafka.EventTypeDelete:
		err = s.uc.Unfollow(ctx, follow.Uid, follow.ToUid)
	}
	return err
}

func (s *RelationshipJob) rebuild(ctx context.Context, msg *kafka.KafkaMessage) error {
	list := &v1.ListRequest{}
	err := proto.Unmarshal(msg.Message, list)
	if err != nil {
		return err
	}
	switch msg.EventType {
	case kafka.EventTypeListMissed:
		err = s.uc.FollowerCacheRebuild(ctx, list.Uid, list.LastId)
	}
	return err
}
