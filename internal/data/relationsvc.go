package data

import (
	"context"
	"sort"

	"github.com/zxq97/gotool/cast"
	"github.com/zxq97/gotool/concurrent"
	"github.com/zxq97/gotool/config"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/gotool/memcachex"
	"github.com/zxq97/gotool/redisx"
	"github.com/zxq97/relation/internal/env"
)

func NewRelationSVCRepo(redisConf *config.RedisConf, mcConf *config.MCConf, conf *config.MysqlConf, addr []string) (*relationSVCRepo, error) {
	repo := &relationSVCRepo{}
	sess, err := conf.InitDB()
	if err != nil {
		return nil, err
	}
	producer, err := kafka.InitKafkaProducer(addr, env.ApiLogger, env.ExcLogger)
	if err != nil {
		return nil, err
	}
	repo.sess = sess
	repo.producer = producer
	repo.redis = redisx.NewRedisX(redisConf)
	repo.mc = memcachex.NewMemcacheX(mcConf)
	return repo, nil
}

func (repo *relationSVCRepo) Follow(ctx context.Context, uid, touid int64) error {
	return sendKafkaMsg(ctx, repo.producer, kafka.TopicRelationFollow, cast.FormatInt(uid), &FollowKafka{Uid: uid, ToUid: touid}, kafka.EventTypeCreate)
}

func (repo *relationSVCRepo) Unfollow(ctx context.Context, uid, touid int64) error {
	return sendKafkaMsg(ctx, repo.producer, kafka.TopicRelationFollow, cast.FormatInt(uid), &FollowKafka{Uid: uid, ToUid: touid}, kafka.EventTypeDelete)
}

func (repo *relationSVCRepo) GetFollowList(ctx context.Context, uid, lastid int64) ([]*FollowItem, error) {
	list, err := getFollowCacheL1(ctx, repo.mc, uid)
	if err != nil {
		list, err = getFollowCacheL2(ctx, repo.redis, uid)
		if err != nil {
			list, err = sfGetUserFollow(ctx, repo.sess, uid)
			if err != nil {
				return nil, err
			}
			concurrent.Go(func() {
				_ = setFollowCacheL1(ctx, repo.mc, uid, list)
			})
		}
		concurrent.Go(func() {
			_ = setFollowCacheL2(ctx, repo.redis, uid, list)
		})
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
	return list, nil
}

func (repo *relationSVCRepo) GetFollowerList(ctx context.Context, uid, lastid int64) ([]*FollowItem, error) {
	list, err := getFollowerList(ctx, repo.redis, uid, lastid, 20)
	if err != nil {
		list, err = sfGetUserFollower(ctx, repo.sess, repo.producer, uid, lastid)
		if err != nil {
			return nil, err
		}
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
	return list, nil
}

func (repo *relationSVCRepo) GetRelation(ctx context.Context, uid int64, uids []int64) (map[int64]*UserRelation, error) {
	followMap, err := isFollows(ctx, repo.redis, uid, uids)
	if err != nil {
		list, err := sfGetUserFollow(ctx, repo.sess, uid)
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
	followerMap, missed, err := isFollowers(ctx, repo.redis, uid, uids)
	if err != nil || len(missed) != 0 {
		dbm, err := sfGetUsersFollow(ctx, repo.sess, missed)
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

func (repo *relationSVCRepo) GetRelationCount(ctx context.Context, uids []int64) (map[int64]*UserFollowCount, error) {
	m, missed, err := getRelationCount(ctx, repo.redis, uids)
	if err != nil || len(missed) != 0 {
		dbm, err := getFollowCount(ctx, repo.sess, missed)
		if err != nil {
			return nil, err
		}
		for k, v := range dbm {
			m[k] = v
		}
		_ = setRelationCount(ctx, repo.redis, dbm)
	}
	return m, nil
}

func (repo *relationSVCRepo) GetUsersFollow(ctx context.Context, uids []int64) (map[int64][]*FollowItem, error) {
	m, missed, err := getFollowsCacheL1(ctx, repo.mc, uids)
	if err != nil || len(missed) != 0 {
		dbm, err := sfGetUsersFollow(ctx, repo.sess, missed)
		if err != nil {
			return nil, err
		}
		for k, v := range dbm {
			m[k] = v
		}
	}
	return m, nil
}
