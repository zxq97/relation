package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/zxq97/gokit/pkg/cast"
	"github.com/zxq97/relation/app/relationship/pkg/bizdata"
)

const (
	relationCacheL2TTL     = time.Hour * 12
	redisKeySyncCount      = "rla_sync_%d" // uid
	redisKeyHUserFollow    = "rla_fow_%d"  // uid
	redisKeyHRelationCount = "lra_cnt_%d"  // uid
	redisKeyZUserFollower  = "rla_foe_%d"  // uid

	followField   = "follow"
	followerField = "follower"
)

func (c *Cache) batchSyncCount(ctx context.Context, uid int64) bool {
	key := fmt.Sprintf(redisKeySyncCount, uid)
	return c.redis.SetNX(ctx, key, "", time.Second).Val()
}

func (c *Cache) addRelationCount(ctx context.Context, uid, touid, incr int64) error {
	key := fmt.Sprintf(redisKeyHRelationCount, uid)
	err := c.redis.HIncrByXEX(ctx, key, followField, incr, relationCacheL2TTL)
	if err != nil {
		err = errors.WithMessage(err, "add follow count")
	}
	key = fmt.Sprintf(redisKeyHRelationCount, touid)
	err = c.redis.HIncrByXEX(ctx, key, followerField, incr, relationCacheL2TTL)
	return errors.WithMessage(err, "add follower count")
}

func (c *Cache) incrFollowerCount(ctx context.Context, uid, cnt int64) error {
	key := fmt.Sprintf(redisKeyHRelationCount, uid)
	return c.redis.HIncrByXEX(ctx, key, followerField, cnt, relationCacheL2TTL)
}

func (c *Cache) addFollower(ctx context.Context, uid int64, item *bizdata.FollowItem) error {
	key := fmt.Sprintf(redisKeyZUserFollower, item.ToUid)
	return c.redis.ZAddXEX(ctx, key, []*redis.Z{{Member: uid, Score: float64(item.CreateTime)}}, time.Hour)
}

func (c *Cache) appendFollower(ctx context.Context, uid int64, list []*bizdata.FollowItem) error {
	key := fmt.Sprintf(redisKeyZUserFollower, uid)
	zs := make([]*redis.Z, len(list))
	for k, v := range list {
		zs[k] = &redis.Z{
			Member: v.ToUid,
			Score:  float64(v.CreateTime),
		}
	}
	return c.redis.ZAddEX(ctx, key, zs, time.Hour)
}

func (c *Cache) addFollowCacheL2(ctx context.Context, uid int64, list []*bizdata.FollowItem) error {
	fieldMap := make(map[string]interface{}, len(list))
	for _, v := range list {
		fieldMap[cast.FormatInt(v.ToUid)] = v.CreateTime
	}
	return c.redis.HMSetXEX(ctx, fmt.Sprintf(redisKeyHUserFollow, uid), fieldMap, relationCacheL2TTL)
}

func (c *Cache) delFollowCacheL2(ctx context.Context, uid, touid int64) error {
	return c.redis.HDel(ctx, fmt.Sprintf(redisKeyHUserFollow, uid), cast.FormatInt(touid)).Err()
}

func (c *Cache) delFollower(ctx context.Context, uid, touid int64) error {
	key := fmt.Sprintf(redisKeyZUserFollower, touid)
	return c.redis.ZRem(ctx, key, uid).Err()
}

func (c *Cache) getFollowsCacheL2(ctx context.Context, uids []int64) (map[int64][]*bizdata.FollowItem, []int64, error) {
	cmdMap := make(map[int64]*redis.StringStringMapCmd, len(uids))
	pipe := c.redis.Pipeline()
	for _, v := range uids {
		key := fmt.Sprintf(redisKeyHUserFollow, v)
		cmdMap[v] = pipe.HGetAll(ctx, key)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cache get follow")
	}
	lm := make(map[int64][]*bizdata.FollowItem, len(uids))
	missed := make([]int64, 0, len(uids))
	for k, v := range cmdMap {
		val := v.Val()
		if len(val) == 0 {
			missed = append(missed, k)
			continue
		}
		list := make([]*bizdata.FollowItem, 0, len(val))
		for t, c := range val {
			list = append(list, &bizdata.FollowItem{
				ToUid:      cast.ParseInt(t, 0),
				CreateTime: cast.ParseInt(c, 0),
			})
		}
		lm[k] = list
	}
	return lm, missed, nil
}

func (c *Cache) setFollowsCacheL2(ctx context.Context, m map[int64][]*bizdata.FollowItem) error {
	// todo diff slot err?
	pipe := c.redis.Pipeline()
	for k, v := range m {
		field := make(map[string]interface{}, len(v))
		for _, u := range v {
			field[cast.FormatInt(u.ToUid)] = u.CreateTime
		}
		key := fmt.Sprintf(redisKeyHUserFollow, k)
		pipe.HMSet(ctx, key, field)
		pipe.Expire(ctx, key, relationCacheL2TTL)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (c *Cache) getFollowerList(ctx context.Context, uid, lastid, offset int64) ([]*bizdata.FollowItem, error) {
	key := fmt.Sprintf(redisKeyZUserFollower, uid)
	var (
		zs  []redis.Z
		err error
	)
	if lastid == 0 {
		zs, err = c.redis.ZRevRangeWithScores(ctx, key, 0, offset-1).Result()
	} else {
		zs, err = c.redis.ZRevRangeByMemberWithScores(ctx, key, lastid, offset)
	}
	if err != nil {
		return nil, err
	}
	list := make([]*bizdata.FollowItem, len(zs))
	for k, z := range zs {
		list[k] = &bizdata.FollowItem{
			ToUid:      cast.ParseInt(z.Member.(string), 0),
			CreateTime: int64(z.Score),
		}
	}
	return list, nil
}

func (c *Cache) isFollows(ctx context.Context, uid int64, uids []int64) (map[int64]int64, error) {
	key := fmt.Sprintf(redisKeyHUserFollow, uid)
	us := make([]string, len(uids))
	for k, v := range uids {
		us[k] = cast.FormatInt(v)
	}
	val, err := c.redis.HMGetXEX(ctx, key, relationCacheL2TTL, us...)
	if err != nil {
		return nil, err
	}
	rm := make(map[int64]int64, len(uids))
	for k, v := range val {
		x, ok := v.(string)
		if v != nil && ok {
			rm[uids[k]] = cast.ParseInt(x, 0)
		}
	}
	return rm, nil
}

func (c *Cache) isFollowers(ctx context.Context, uid int64, uids []int64) (map[int64]int64, []int64, error) {
	missed := make([]int64, 0, len(uids))
	rm := make(map[int64]int64, len(uids))
	for _, v := range uids {
		key := fmt.Sprintf(redisKeyHUserFollow, v)
		val, err := c.redis.HGetXEX(ctx, key, cast.FormatInt(uid), relationCacheL2TTL)
		if err != nil {
			missed = append(missed, v)
			continue
		} else {
			rm[v] = cast.ParseInt(val, 0)
		}
	}
	return rm, missed, nil
}

func (c *Cache) getRelationCount(ctx context.Context, uids []int64) (map[int64]*bizdata.RelationCount, []int64, error) {
	pipe := c.redis.Pipeline()
	cmdMap := make(map[int64]*redis.StringStringMapCmd, len(uids))
	for _, v := range uids {
		key := fmt.Sprintf(redisKeyHRelationCount, v)
		cmdMap[v] = pipe.HGetAll(ctx, key)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, uids, err
	}
	missed := make([]int64, 0, len(uids))
	fm := make(map[int64]*bizdata.RelationCount, len(uids))
	for k, v := range cmdMap {
		val, err := v.Result()
		if err != nil {
			missed = append(missed, k)
			continue
		}
		fc := &bizdata.RelationCount{}
		if c, ok := val[followField]; ok {
			fc.FollowCount = cast.Atoi(c, 0)
		}
		if c, ok := val[followerField]; ok {
			fc.FollowerCount = cast.Atoi(c, 0)
		}
	}
	return fm, missed, nil
}

func (c *Cache) setRelationCount(ctx context.Context, fm map[int64]*bizdata.RelationCount) error {
	pipe := c.redis.Pipeline()
	for k, v := range fm {
		key := fmt.Sprintf(redisKeyHRelationCount, k)
		fieldMap := map[string]interface{}{
			followField:   v.FollowCount,
			followerField: v.FollowerCount,
		}
		pipe.HMSet(ctx, key, fieldMap)
		pipe.Expire(ctx, key, relationCacheL2TTL)
	}
	_, err := pipe.Exec(ctx)
	return err
}
