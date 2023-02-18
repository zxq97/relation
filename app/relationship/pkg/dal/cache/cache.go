package cache

import (
	"context"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/patrickmn/go-cache"
	"github.com/zxq97/gokit/pkg/cache/xredis"
	"github.com/zxq97/gokit/pkg/concurrent"
	"github.com/zxq97/relation/app/relationship/pkg/bizdata"
	"github.com/zxq97/relation/app/relationship/pkg/message"
)

type Cache struct {
	cache *cache.Cache
	redis *xredis.XRedis
	mc    *memcache.Client
}

func Use(redis *xredis.XRedis, mc *memcache.Client) *Cache {
	return &Cache{
		cache: cache.New(time.Second*5, time.Minute*10),
		redis: redis,
		mc:    mc,
	}
}

func (c *Cache) Follow(ctx context.Context, uid, touid int64) error {
	now := time.Now().UnixMilli()
	list := []*bizdata.FollowItem{{ToUid: touid, CreateTime: now}}
	eg := concurrent.NewErrGroup(ctx)
	eg.Go(func() error {
		return c.addFollowCacheL1(ctx, uid, list)
	})
	eg.Go(func() error {
		return c.addFollowCacheL2(ctx, uid, list)
	})
	eg.Go(func() error {
		return c.addFollower(ctx, touid, &bizdata.FollowItem{ToUid: uid, CreateTime: now})
	})
	eg.Go(func() error {
		return c.addRelationCount(ctx, uid, touid, 1)
	})
	return eg.Wait()
}

func (c *Cache) Unfollow(ctx context.Context, uid, touid int64) error {
	eg := concurrent.NewErrGroup(ctx)
	eg.Go(func() error {
		return c.delFollowCacheL1(ctx, uid, touid)
	})
	eg.Go(func() error {
		return c.delFollowCacheL2(ctx, uid, touid)
	})
	eg.Go(func() error {
		return c.delFollower(ctx, uid, touid)
	})
	eg.Go(func() error {
		return c.addRelationCount(ctx, uid, touid, -1)
	})
	return eg.Wait()
}

func (c *Cache) AppendFollower(ctx context.Context, uid int64, list []*bizdata.FollowItem) error {
	return c.appendFollower(ctx, uid, list)
}

func (c *Cache) BatchSyncCount(ctx context.Context, uid int64) bool {
	return c.batchSyncCount(ctx, uid)
}

func (c *Cache) IsFollow(ctx context.Context, uid int64, uids []int64) (map[int64]int64, error) {
	// fixme batch size?
	if len(uids) > 100 {
		m, _, _ := c.GetUsersFollow(ctx, []int64{uid})
		l, ok := m[uid]
		if ok {
			fm := make(map[int64]int64, len(uids))
			s := make(map[int64]int64, len(l))
			for _, v := range l {
				s[v.ToUid] = v.CreateTime
			}
			for _, v := range uids {
				if x, ok := s[v]; ok {
					fm[v] = x
				}
			}
		}
	}
	return c.isFollows(ctx, uid, uids)
}

func (c *Cache) IsFollower(ctx context.Context, uid int64, uids []int64) (map[int64]int64, []int64, error) {
	return c.isFollowers(ctx, uid, uids)
}

func (c *Cache) GetUsersFollow(ctx context.Context, uids []int64) (map[int64][]*bizdata.FollowItem, []int64, error) {
	var (
		m       map[int64][]*bizdata.FollowItem
		m2      map[int64][]*bizdata.FollowItem
		missed  []int64
		missed2 []int64
		err     error
	)
	m, missed, err = c.getFollowsCacheL1(ctx, uids)
	if err != nil || len(missed) != 0 {
		m2, missed2, _ = c.getFollowsCacheL2(ctx, missed)
		for k, v := range m2 {
			m[k] = v
			_ = c.setFollowCacheL1(ctx, k, v)
		}
	}
	return m, missed2, nil
}

func (c *Cache) SetUsersFollow(ctx context.Context, m map[int64][]*bizdata.FollowItem) error {
	_ = c.setFollowsCacheL2(ctx, m)
	for k, v := range m {
		_ = c.setFollowCacheL1(ctx, k, v)
	}
	return nil
}

func (c *Cache) GetFollowerList(ctx context.Context, uid, lastid int64) ([]*bizdata.FollowItem, error) {
	return c.getFollowerList(ctx, uid, lastid, message.ListBatchSize)
}

func (c *Cache) GetRelationCount(ctx context.Context, uids []int64) (map[int64]*bizdata.RelationCount, []int64, error) {
	return c.getRelationCount(ctx, uids)
}

func (c *Cache) SetRelationCount(ctx context.Context, m map[int64]*bizdata.RelationCount) error {
	return c.setRelationCount(ctx, m)
}

func (c *Cache) IncrFollowerCount(ctx context.Context, uid, cnt int64) error {
	return c.incrFollowerCount(ctx, uid, cnt)
}
