package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/zxq97/gotool/cast"
	"github.com/zxq97/gotool/redisx"
)

const (
	relationCacheL2TTL     = time.Hour * 12
	redisKeyHUserFollow    = "rla_fow_%d" // uid
	redisKeyHRelationCount = "lra_cnt_%d" // uid
	redisKeyZUserFollower  = "rla_foe_%d" // uid

	followField   = "follow"
	followerField = "follower"
)

func getFollowCacheL2(ctx context.Context, rdx *redisx.RedisX, uid int64) ([]*FollowItem, error) {
	key := fmt.Sprintf(redisKeyHUserFollow, uid)
	val, err := rdx.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "cache get follow")
	} else if len(val) == 0 {
		return nil, redis.Nil
	}
	list := make([]*FollowItem, 0, len(val))
	for k, v := range val {
		list = append(list, &FollowItem{
			ToUid:      cast.ParseInt(k, 0),
			CreateTime: cast.ParseInt(v, 0),
		})
	}
	return list, nil
}

func setFollowCacheL2(ctx context.Context, rdx *redisx.RedisX, uid int64, list []*FollowItem) error {
	field := make(map[string]interface{}, len(list))
	for _, v := range list {
		field[cast.FormatInt(v.ToUid)] = v.CreateTime
	}
	key := fmt.Sprintf(redisKeyHUserFollow, uid)
	return rdx.HMSetEX(ctx, key, field, relationCacheL2TTL)
}

func addFollowCacheL2(ctx context.Context, rdx *redisx.RedisX, uid int64, list []*FollowItem) error {
	fieldMap := make(map[string]interface{}, len(list))
	for _, v := range list {
		fieldMap[cast.FormatInt(v.ToUid)] = v.CreateTime
	}
	return rdx.HMSetXEX(ctx, fmt.Sprintf(redisKeyHUserFollow, uid), fieldMap, relationCacheL2TTL)
}

func delFollowCacheL2(ctx context.Context, rdx *redisx.RedisX, uid, touid int64) error {
	return rdx.HDel(ctx, fmt.Sprintf(redisKeyHUserFollow, uid), cast.FormatInt(touid)).Err()
}

func addFollower(ctx context.Context, rdx *redisx.RedisX, uid int64, item *FollowItem) error {
	key := fmt.Sprintf(redisKeyZUserFollower, item.ToUid)
	return rdx.ZAddXEX(ctx, key, []*redis.Z{{Member: uid, Score: float64(item.CreateTime)}}, time.Hour)
}

func delFollower(ctx context.Context, rdx *redisx.RedisX, uid, touid int64) error {
	key := fmt.Sprintf(redisKeyZUserFollower, touid)
	return rdx.ZRem(ctx, key, uid).Err()
}

func addRelationCount(ctx context.Context, rdx *redisx.RedisX, uid, touid, incr int64) {
	key := fmt.Sprintf(redisKeyHRelationCount, uid)
	_ = rdx.HIncrByXEX(ctx, key, followField, incr, relationCacheL2TTL)
	key = fmt.Sprintf(redisKeyHRelationCount, touid)
	_ = rdx.HIncrByXEX(ctx, key, followerField, incr, relationCacheL2TTL)
}

func addUserFollower(ctx context.Context, rdx *redisx.RedisX, uid int64, list []*FollowItem) error {
	key := fmt.Sprintf(redisKeyZUserFollower, uid)
	zs := make([]*redis.Z, len(list))
	for k, v := range list {
		zs[k] = &redis.Z{
			Member: v.ToUid,
			Score:  float64(v.CreateTime),
		}
	}
	return rdx.ZAddXEX(ctx, key, zs, time.Hour)
}

func getFollowerList(ctx context.Context, rdx *redisx.RedisX, uid, lastid, offset int64) ([]*FollowItem, error) {
	key := fmt.Sprintf(redisKeyZUserFollower, uid)
	zs, err := rdx.ZRevRangeByMemberWithScores(ctx, key, lastid, offset)
	if err != nil {
		return nil, err
	}
	list := make([]*FollowItem, len(zs))
	for k, z := range zs {
		list[k] = &FollowItem{
			ToUid:      z.Member.(int64),
			CreateTime: int64(z.Score),
		}
	}
	return list, nil
}

func setRelationCount(ctx context.Context, rdx *redisx.RedisX, fm map[int64]*UserFollowCount) error {
	pipe := rdx.Pipeline()
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

func isFollows(ctx context.Context, rdx *redisx.RedisX, uid int64, uids []int64) (map[int64]int64, error) {
	key := fmt.Sprintf(redisKeyHUserFollow, uid)
	us := make([]string, len(uids))
	for k, v := range uids {
		us[k] = cast.FormatInt(v)
	}
	val, err := rdx.HMGetXEX(ctx, key, relationCacheL2TTL, us...)
	if err != nil {
		return nil, err
	}
	rm := make(map[int64]int64, len(uids))
	for k, v := range val {
		c, ok := v.(string)
		if v != nil && ok {
			rm[uids[k]] = cast.ParseInt(c, 0)
		}
	}
	return rm, nil
}

func isFollowers(ctx context.Context, rdx *redisx.RedisX, uid int64, uids []int64) (map[int64]int64, []int64, error) {
	missed := make([]int64, 0, len(uids))
	rm := make(map[int64]int64, len(uids))
	for _, v := range uids {
		key := fmt.Sprintf(redisKeyHUserFollow, v)
		val, err := rdx.HGetXEX(ctx, key, cast.FormatInt(uid), relationCacheL2TTL)
		if err != nil {
			missed = append(missed, v)
			continue
		} else {
			rm[v] = cast.ParseInt(val, 0)
		}
	}
	return rm, missed, nil
}

func getRelationCount(ctx context.Context, rdx *redisx.RedisX, uids []int64) (map[int64]*UserFollowCount, []int64, error) {
	pipe := rdx.Pipeline()
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
	fm := make(map[int64]*UserFollowCount, len(uids))
	for k, v := range cmdMap {
		val, err := v.Result()
		if err != nil {
			missed = append(missed, k)
			continue
		}
		fc := &UserFollowCount{}
		if c, ok := val[followField]; ok {
			fc.FollowCount = cast.Atoi(c, 0)
		}
		if c, ok := val[followerField]; ok {
			fc.FollowerCount = cast.Atoi(c, 0)
		}
	}
	return fm, missed, nil
}
