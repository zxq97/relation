package data

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/zxq97/gotool/concurrent"
	"github.com/zxq97/gotool/config"
	"github.com/zxq97/gotool/memcachex"
	"github.com/zxq97/gotool/redisx"
)

func NewRelationTaskRepoImpl(redisConf *config.RedisConf, mcConf *config.MCConf, mysqlConf *config.MysqlConf) (*RelationTaskRepoImpl, error) {
	repo := &RelationTaskRepoImpl{}
	sess, err := mysqlConf.InitDB()
	if err != nil {
		return nil, err
	}
	repo.sess = sess
	repo.redis = redisx.NewRedisX(redisConf)
	repo.mc = memcachex.NewMemcacheX(mcConf)
	repo.cache = cache.New(5*time.Minute, 15*time.Minute)
	return repo, nil
}

func (repo *RelationTaskRepoImpl) Follow(ctx context.Context, uid, touid int64) error {
	err := follow(ctx, repo.sess, uid, touid)
	if err != nil {
		return err
	}
	item := &FollowItem{ToUid: touid, CreateTime: time.Now().UnixMilli()}
	list := []*FollowItem{item}
	eg := concurrent.NewErrGroup(ctx)
	eg.Go(func() error {
		return addFollowCacheL1(ctx, repo.mc, uid, list)
	})
	eg.Go(func() error {
		return addFollowCacheL2(ctx, repo.redis, uid, list)
	})
	eg.Go(func() error {
		return addFollower(ctx, repo.redis, uid, item)
	})
	eg.Go(func() error {
		addRelationCount(ctx, repo.redis, uid, touid, 1)
		return nil
	})
	return eg.Wait()
}

func (repo *RelationTaskRepoImpl) Unfollow(ctx context.Context, uid, touid int64) error {
	err := unfollow(ctx, repo.sess, uid, touid)
	if err != nil {
		return err
	}
	eg := concurrent.NewErrGroup(ctx)
	eg.Go(func() error {
		return delFollowCacheL1(ctx, repo.mc, uid, touid)
	})
	eg.Go(func() error {
		return delFollowCacheL2(ctx, repo.redis, uid, touid)
	})
	eg.Go(func() error {
		return delFollower(ctx, repo.redis, uid, touid)
	})
	eg.Go(func() error {
		addRelationCount(ctx, repo.redis, uid, touid, -1)
		return nil
	})
	return eg.Wait()
}

func (repo *RelationTaskRepoImpl) FollowerCacheRebuild(ctx context.Context, uid, lastid int64) error {
	key := fmt.Sprintf(lcKeyFollower, uid, lastid)
	_, ok := lcGet(repo.cache, key)
	if ok {
		return nil
	}
	list, err := getUserFollower(ctx, repo.sess, uid, lastid, 100)
	if err != nil {
		return err
	}
	err = addUserFollower(ctx, repo.redis, uid, list)
	if err != nil {
		return err
	}
	lcSet(repo.cache, key, nil, time.Minute*5)
	return nil
}
