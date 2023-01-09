package data

import (
	"context"

	"github.com/zxq97/gotool/config"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/gotool/memcachex"
	"github.com/zxq97/gotool/redisx"
)

func NewRelationSvcRepoImpl(redisConf *config.RedisConf, mcConf *config.MCConf, mysqlConf *config.MysqlConf, producer *kafka.Producer) (*RelationSvcRepoImpl, error) {
	repo := &RelationSvcRepoImpl{}
	sess, err := mysqlConf.InitDB()
	if err != nil {
		return nil, err
	}
	repo.producer = producer
	repo.sess = sess
	repo.redis = redisx.NewRedisX(redisConf)
	repo.mc = memcachex.NewMemcacheX(mcConf)
	return repo, nil
}

func (repo *RelationSvcRepoImpl) GetUserFollow(ctx context.Context, uid int64) ([]*FollowItem, error) {
	list, err := getFollowCacheL1(ctx, repo.mc, uid)
	if err != nil {
		list, err = getFollowCacheL2(ctx, repo.redis, uid)
		if err != nil {
			list, err = sfGetUserFollow(ctx, repo.sess, uid)
			if err != nil {
				return nil, err
			}
			_ = setFollowCacheL1(ctx, repo.mc, uid, list)
		}
		_ = setFollowCacheL2(ctx, repo.redis, uid, list)
	}
	return list, nil
}

func (repo *RelationSvcRepoImpl) GetUserFollower(ctx context.Context, uid, lastid int64) ([]*FollowItem, error) {
	list, err := getFollowerList(ctx, repo.redis, uid, lastid, 20)
	if err != nil {
		list, err = sfGetUserFollower(ctx, repo.sess, repo.producer, uid, lastid)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (repo *RelationSvcRepoImpl) GetIsFollowMap(ctx context.Context, uid int64, uids []int64) (map[int64]int64, error) {
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
	return followMap, nil
}

func (repo *RelationSvcRepoImpl) GetIsFollowerMap(ctx context.Context, uid int64, uids []int64) (map[int64]int64, error) {
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
	return followerMap, nil
}

func (repo *RelationSvcRepoImpl) GetRelationCount(ctx context.Context, uids []int64) (map[int64]*UserFollowCount, error) {
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

func (repo *RelationSvcRepoImpl) GetUsersFollow(ctx context.Context, uids []int64) (map[int64][]*FollowItem, error) {
	m, missed, err := getFollowsCacheL1(ctx, repo.mc, uids)
	if err != nil || len(missed) != 0 {
		m2, missed2, err := getFollowsCacheL2(ctx, repo.redis, missed)
		if err != nil || len(missed2) != 0 {
			dbm, err := sfGetUsersFollow(ctx, repo.sess, missed2)
			if err != nil {
				return nil, err
			}
			for k, v := range dbm {
				m2[k] = v
			}
			_ = setFollowsCacheL2(ctx, repo.redis, dbm)
		}
		for k, v := range m2 {
			m[k] = v
			_ = setFollowCacheL1(ctx, repo.mc, k, v)
		}
	}
	return m, nil
}
