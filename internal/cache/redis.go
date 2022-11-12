package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zxq97/gotool/cast"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/model"
)

const (
	relationCacheL2TTL     = time.Hour * 12
	redisKeyHUserFollow    = "rla_fow_%d"
	redisKeyHRelationCount = "lra_cnt_%d"

	followField   = "follow"
	followerField = "follower"
)

func getRelationCacheL2(ctx context.Context, keyPrefix string, uid int64) ([]*model.FollowItem, error) {
	key := fmt.Sprintf(keyPrefix, uid)
	val, err := rdx.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	} else if len(val) == 0 {
		return nil, redis.Nil
	}
	list := make([]*model.FollowItem, 0, len(val))
	for k, v := range val {
		list = append(list, &model.FollowItem{
			ToUid:      cast.ParseInt(k, 0),
			CreateTime: cast.ParseInt(v, 0),
		})
	}
	return list, nil
}

func setRelationCacheL2(ctx context.Context, keyPrefix string, uid int64, list []*model.FollowItem) error {
	field := make(map[string]interface{}, len(list))
	for _, v := range list {
		field[cast.FormatInt(v.ToUid)] = v.CreateTime
	}
	key := fmt.Sprintf(keyPrefix, uid)
	return rdx.HMSetEX(ctx, key, field, relationCacheL2TTL)
}

func addRelationCacheL2(ctx context.Context, keyPrefix string, uid int64, list []*model.FollowItem) error {
	fieldMap := make(map[string]interface{}, len(list))
	for _, v := range list {
		fieldMap[cast.FormatInt(v.ToUid)] = v.CreateTime
	}
	return rdx.HMSetXEX(ctx, fmt.Sprintf(keyPrefix, uid), fieldMap, relationCacheL2TTL)
}

func delRelationCacheL2(ctx context.Context, keyPrefix string, uid, touid int64) error {
	return rdx.HDel(ctx, fmt.Sprintf(keyPrefix, uid), cast.FormatInt(touid)).Err()
}

func addRelationCount(ctx context.Context, uid, touid, incr int64) {
	key := fmt.Sprintf(redisKeyHRelationCount, uid)
	_ = rdx.HIncrByXEX(ctx, key, followField, incr, relationCacheL2TTL)
	key = fmt.Sprintf(redisKeyHRelationCount, touid)
	_ = rdx.HIncrByXEX(ctx, key, followerField, incr, relationCacheL2TTL)
}

func SetFollowList(ctx context.Context, uid int64, list []*model.FollowItem) {
	_ = setRelationCacheL1(ctx, mcKeyUserFollow, uid, list)
	_ = setRelationCacheL2(ctx, redisKeyHUserFollow, uid, list)
}

func SetRelationCount(ctx context.Context, fm map[int64]*model.UserFollowCount) {
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
	if err != nil {
		env.ExcLogger.Println()
	}
}

func IsFollows(ctx context.Context, uid int64, uids []int64) (map[int64]int64, error) {
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

func IsFollowers(ctx context.Context, uid int64, uids []int64) (map[int64]int64, []int64, error) {
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

func GetRelationCount(ctx context.Context, uids []int64) (map[int64]*model.UserFollowCount, []int64, error) {
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
	fm := make(map[int64]*model.UserFollowCount, len(uids))
	for k, v := range cmdMap {
		val, err := v.Result()
		if err != nil {
			missed = append(missed, k)
			continue
		}
		fc := &model.UserFollowCount{}
		if c, ok := val[followField]; ok {
			fc.FollowCount = cast.Atoi(c, 0)
		}
		if c, ok := val[followerField]; ok {
			fc.FollowerCount = cast.Atoi(c, 0)
		}
	}
	return fm, missed, nil
}
