package data

import (
	"context"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/patrickmn/go-cache"
	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/gokit/pkg/concurrent"
	"github.com/zxq97/relation/app/relationship/job/internal/biz"
	"gorm.io/gorm"
)

var _ biz.RelationshipRepo = (*relationshipRepo)(nil)

type relationshipRepo struct {
	cache *cache.Cache
	mc    *memcache.Client
	redis *xredis.XRedis
	db    *gorm.DB
}

func NewRelationshipRepo(mc *memcache.Client, redis *xredis.XRedis, db *gorm.DB) *relationshipRepo {
	return &relationshipRepo{
		cache: cache.New(time.Minute*5, time.Minute*15),
		mc:    mc,
		redis: redis,
		db:    db,
	}
}

func (r *relationshipRepo) Follow(ctx context.Context, uid, touid int64) error {
	err := r.follow(ctx, uid, touid)
	if err != nil {
		return err
	}
	now := time.Now().UnixMilli()
	list := []*biz.FollowItem{{ToUid: touid, CreateTime: now}}
	eg := concurrent.NewErrGroup(ctx)
	eg.Go(func() error {
		return r.addFollowCacheL1(ctx, uid, list)
	})
	eg.Go(func() error {
		return r.addFollowCacheL2(ctx, uid, list)
	})
	eg.Go(func() error {
		return r.addFollower(ctx, touid, &biz.FollowItem{ToUid: uid, CreateTime: now})
	})
	eg.Go(func() error {
		return r.addRelationCount(ctx, uid, touid, 1)
	})
	return eg.Wait()
}

func (r *relationshipRepo) Unfollow(ctx context.Context, uid, touid int64) error {
	err := r.unfollow(ctx, uid, touid)
	if err != nil {
		return err
	}
	eg := concurrent.NewErrGroup(ctx)
	eg.Go(func() error {
		return r.delFollowCacheL1(ctx, uid, touid)
	})
	eg.Go(func() error {
		return r.delFollowCacheL2(ctx, uid, touid)
	})
	eg.Go(func() error {
		return r.delFollower(ctx, uid, touid)
	})
	eg.Go(func() error {
		return r.addRelationCount(ctx, uid, touid, -1)
	})
	return eg.Wait()
}

func (r *relationshipRepo) FollowerCacheRebuild(ctx context.Context, uid, lastid int64) error {
	key := fmt.Sprintf(lcKeyFollower, uid)
	_, ok := r.localCacheGet(key)
	if ok {
		return nil
	}
	list, err := r.getUserFollower(ctx, uid, lastid)
	if err != nil {
		return err
	}
	err = r.appendFollower(ctx, uid, list)
	if err != nil {
		return err
	}
	r.localCacheSet(key, struct{}{}, time.Minute)
	return nil
}
