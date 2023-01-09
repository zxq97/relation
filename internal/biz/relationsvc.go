package biz

import (
	"context"
	"sort"

	"github.com/zxq97/gotool/cast"
	"github.com/zxq97/gotool/config"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/relation/internal/data"
	"github.com/zxq97/relation/internal/env"
)

type RelationSvcRepo interface {
	GetRelationCount(context.Context, []int64) (map[int64]*data.UserFollowCount, error)
	GetUserFollow(context.Context, int64) ([]*data.FollowItem, error)
	GetUsersFollow(context.Context, []int64) (map[int64][]*data.FollowItem, error)
	GetUserFollower(context.Context, int64, int64) ([]*data.FollowItem, error)
	GetIsFollowMap(context.Context, int64, []int64) (map[int64]int64, error)
	GetIsFollowerMap(context.Context, int64, []int64) (map[int64]int64, error)
}

type RelationSvcBIZ struct {
	producer *kafka.Producer
	repo     RelationSvcRepo
}

func NewRelationSvcBIZ(redisConf *config.RedisConf, mcConf *config.MCConf, mysqlConf *config.MysqlConf, addr []string) (*RelationSvcBIZ, error) {
	producer, err := kafka.InitKafkaProducer(addr, env.ApiLogger, env.ExcLogger)
	if err != nil {
		return nil, err
	}
	repo, err := data.NewRelationSvcRepoImpl(redisConf, mcConf, mysqlConf, producer)
	if err != nil {
		return nil, err
	}
	return &RelationSvcBIZ{
		producer: producer,
		repo:     repo,
	}, nil
}

func (rsb *RelationSvcBIZ) Follow(ctx context.Context, uid, touid int64) error {
	return rsb.producer.SendKafkaMsg(ctx, kafka.TopicRelationFollow, cast.FormatInt(uid), &data.FollowKafka{Uid: uid, ToUid: touid}, kafka.EventTypeCreate)
}

func (rsb *RelationSvcBIZ) Unfollow(ctx context.Context, uid, touid int64) error {
	return rsb.producer.SendKafkaMsg(ctx, kafka.TopicRelationFollow, cast.FormatInt(uid), &data.FollowKafka{Uid: uid, ToUid: touid}, kafka.EventTypeDelete)
}

func (rsb *RelationSvcBIZ) GetFollowList(ctx context.Context, uid, lastid int64) ([]*FollowItem, error) {
	list, err := rsb.repo.GetUserFollow(ctx, uid)
	if err != nil {
		return nil, err
	}
	idx := sort.Search(len(list), func(i int) bool {
		return list[i].ToUid == lastid
	})
	if idx != len(list) {
		right := idx + 20
		if right > len(list) {
			right = len(list)
		}
		list = list[idx:right]
	}
	return listPO2DO(list), nil
}

func (rsb *RelationSvcBIZ) GetFollowerList(ctx context.Context, uid, lastid int64) ([]*FollowItem, error) {
	list, err := rsb.repo.GetUserFollower(ctx, uid, lastid)
	if err != nil {
		return nil, err
	}
	return listPO2DO(list), nil
}

func (rsb *RelationSvcBIZ) GetRelation(ctx context.Context, uid int64, uids []int64) (map[int64]*UserRelation, error) {
	followMap, err := rsb.repo.GetIsFollowMap(ctx, uid, uids)
	if err != nil {
		return nil, err
	}
	followerMap, err := rsb.repo.GetIsFollowerMap(ctx, uid, uids)
	if err != nil {
		return nil, err
	}
	relationMap := make(map[int64]*UserRelation, len(uids))
	for k, v := range followMap {
		relationMap[k].Relation |= 1
		relationMap[k].FollowTime = v
	}
	for k, v := range followerMap {
		relationMap[k].Relation |= 2
		relationMap[k].FollowedTime = v
	}
	return relationMap, nil
}

func (rsb *RelationSvcBIZ) GetRelationCount(ctx context.Context, uids []int64) (map[int64]*UserFollowCount, error) {
	m, err := rsb.repo.GetRelationCount(ctx, uids)
	if err != nil {
		return nil, err
	}
	return cmPO2DO(m), nil
}

func (rsb *RelationSvcBIZ) GetUsersFollow(ctx context.Context, uids []int64) (map[int64][]*FollowItem, error) {
	m, err := rsb.repo.GetUsersFollow(ctx, uids)
	if err != nil {
		return nil, err
	}
	return imPO2DO(m), nil
}
