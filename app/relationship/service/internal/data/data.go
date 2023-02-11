package data

import (
	"context"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/gokit/pkg/cast"
	"github.com/zxq97/gokit/pkg/mq"
	"github.com/zxq97/gokit/pkg/mq/kafka"
	"github.com/zxq97/relation/app/relationship/service/internal/biz"
	"gorm.io/gorm"
)

var _ biz.RelationshipRepo = (*relationshipRepo)(nil)

type relationshipRepo struct {
	producer *kafka.Producer
	mc       *memcache.Client
	redis    *xredis.XRedis
	db       *gorm.DB
}

func NewRelationshipRepo(producer *kafka.Producer, mc *memcache.Client, redis *xredis.XRedis, db *gorm.DB) biz.RelationshipRepo {
	return &relationshipRepo{
		producer: producer,
		mc:       mc,
		redis:    redis,
		db:       db,
	}
}

func (r *relationshipRepo) Follow(ctx context.Context, uid, touid int64) error {
	return r.producer.SendMessage(ctx, kafka.TopicRelationFollow, cast.FormatInt(uid), mq.TagCreate, &biz.FollowKafka{Uid: uid, ToUid: touid, CreateTime: time.Now().UnixMilli()})
}

func (r *relationshipRepo) Unfollow(ctx context.Context, uid, touid int64) error {
	return r.producer.SendMessage(ctx, kafka.TopicRelationFollow, cast.FormatInt(uid), mq.TagDelete, &biz.FollowKafka{Uid: uid, ToUid: touid, CreateTime: time.Now().UnixMilli()})
}

func (r *relationshipRepo) GetRelationCount(ctx context.Context, uids []int64) (map[int64]*biz.RelationCount, error) {
	m, missed, err := r.cacheGetRelationCount(ctx, uids)
	if err != nil || len(missed) != 0 {
		dbm, err := r.getRelationCount(ctx, missed)
		if err != nil {
			return nil, err
		}
		for k, v := range dbm {
			m[k] = v
		}
		_ = r.setRelationCount(ctx, dbm)
	}
	return m, nil
}

func (r *relationshipRepo) GetUsersFollow(ctx context.Context, uids []int64) (map[int64][]*biz.FollowItem, error) {
	m, missed, err := r.getFollowsCacheL1(ctx, uids)
	if err != nil || len(missed) != 0 {
		m2, missed2, err := r.getFollowsCacheL2(ctx, missed)
		if err != nil || len(missed2) != 0 {
			dbm, err := r.sfGetUsersFollow(ctx, missed2)
			if err != nil {
				return nil, err
			}
			for k, v := range dbm {
				m2[k] = v
			}
			_ = r.setFollowsCacheL2(ctx, dbm)
		}
		for k, v := range m2 {
			m[k] = v
			_ = r.setFollowCacheL1(ctx, k, v)
		}
	}
	return m, nil
}

func (r *relationshipRepo) GetUserFollower(ctx context.Context, uid int64, lastid int64) ([]*biz.FollowItem, error) {
	list, err := r.getFollowerList(ctx, uid, lastid, 20)
	if err != nil {
		list, err = r.sfGetUserFollower(ctx, uid, lastid)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (r *relationshipRepo) GetIsFollowMap(ctx context.Context, uid int64, uids []int64) (map[int64]int64, error) {
	followMap, err := r.isFollows(ctx, uid, uids)
	if err != nil {
		list, err := r.sfGetUserFollow(ctx, uid)
		if err != nil {
			return nil, err
		}
		m := make(map[int64]int64, len(list))
		for _, v := range list {
			m[v.ToUid] = v.CreateTime
		}
		for _, v := range uids {
			if t, ok := m[v]; ok {
				followMap[v] = t
			}
		}
	}
	return followMap, nil
}

func (r *relationshipRepo) GetIsFollowerMap(ctx context.Context, uid int64, uids []int64) (map[int64]int64, error) {
	followerMap, missed, err := r.isFollowers(ctx, uid, uids)
	if err != nil || len(missed) != 0 {
		dbm, err := r.sfGetUsersFollow(ctx, missed)
		if err != nil {
			return nil, err
		}
		for k, v := range dbm {
			for _, f := range v {
				if f.ToUid == uid {
					followerMap[k] = f.CreateTime
					break
				}
			}
		}
	}
	return followerMap, nil
}
