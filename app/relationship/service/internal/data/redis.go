package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/zxq97/gokit/pkg/cast"
	"github.com/zxq97/relation/app/relationship/service/internal/biz"
)

const (
	relationCacheL2TTL     = time.Hour * 12
	redisKeyHUserFollow    = "rla_fow_%d" // uid
	redisKeyHRelationCount = "lra_cnt_%d" // uid
	redisKeyZUserFollower  = "rla_foe_%d" // uid

	followField   = "follow"
	followerField = "follower"
)

func (r *relationshipRepo) getFollowsCacheL2(ctx context.Context, uids []int64) (map[int64][]*biz.FollowItem, []int64, error) {
	cmdMap := make(map[int64]*redis.StringStringMapCmd, len(uids))
	pipe := r.redis.Pipeline()
	for _, v := range uids {
		key := fmt.Sprintf(redisKeyHUserFollow, v)
		cmdMap[v] = pipe.HGetAll(ctx, key)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cache get follow")
	}
	lm := make(map[int64][]*biz.FollowItem, len(uids))
	missed := make([]int64, 0, len(uids))
	for k, v := range cmdMap {
		val := v.Val()
		if len(val) == 0 {
			missed = append(missed, k)
			continue
		}
		list := make([]*biz.FollowItem, 0, len(val))
		for t, c := range val {
			list = append(list, &biz.FollowItem{
				ToUid:      cast.ParseInt(t, 0),
				CreateTime: cast.ParseInt(c, 0),
			})
		}
		lm[k] = list
	}
	return lm, missed, nil
}

func (r *relationshipRepo) getFollowCacheL2(ctx context.Context, uid int64) ([]*biz.FollowItem, error) {
	key := fmt.Sprintf(redisKeyHUserFollow, uid)
	val, err := r.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "cache get follow")
	} else if len(val) == 0 {
		return nil, redis.Nil
	}
	list := make([]*biz.FollowItem, 0, len(val))
	for k, v := range val {
		list = append(list, &biz.FollowItem{
			ToUid:      cast.ParseInt(k, 0),
			CreateTime: cast.ParseInt(v, 0),
		})
	}
	return list, nil
}

func (r *relationshipRepo) setFollowCacheL2(ctx context.Context, uid int64, list []*biz.FollowItem) error {
	field := make(map[string]interface{}, len(list))
	for _, v := range list {
		field[cast.FormatInt(v.ToUid)] = v.CreateTime
	}
	key := fmt.Sprintf(redisKeyHUserFollow, uid)
	return r.redis.HMSetEX(ctx, key, field, relationCacheL2TTL)
}

func (r *relationshipRepo) setFollowsCacheL2(ctx context.Context, m map[int64][]*biz.FollowItem) error {
	pipe := r.redis.Pipeline()
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

func (r *relationshipRepo) getFollowerList(ctx context.Context, uid, lastid, offset int64) ([]*biz.FollowItem, error) {
	key := fmt.Sprintf(redisKeyZUserFollower, uid)
	var (
		zs  []redis.Z
		err error
	)
	if lastid == 0 {
		zs, err = r.redis.ZRevRangeWithScores(ctx, key, 0, 19).Result()
	} else {
		zs, err = r.redis.ZRevRangeByMemberWithScores(ctx, key, lastid, offset)
	}
	if err != nil {
		return nil, err
	}
	list := make([]*biz.FollowItem, len(zs))
	for k, z := range zs {
		list[k] = &biz.FollowItem{
			ToUid:      z.Member.(int64),
			CreateTime: int64(z.Score),
		}
	}
	return list, nil
}

func (r *relationshipRepo) setRelationCount(ctx context.Context, fm map[int64]*biz.RelationCount) error {
	pipe := r.redis.Pipeline()
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

func (r *relationshipRepo) isFollows(ctx context.Context, uid int64, uids []int64) (map[int64]int64, error) {
	key := fmt.Sprintf(redisKeyHUserFollow, uid)
	us := make([]string, len(uids))
	for k, v := range uids {
		us[k] = cast.FormatInt(v)
	}
	val, err := r.redis.HMGetXEX(ctx, key, relationCacheL2TTL, us...)
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

func (r *relationshipRepo) isFollowers(ctx context.Context, uid int64, uids []int64) (map[int64]int64, []int64, error) {
	missed := make([]int64, 0, len(uids))
	rm := make(map[int64]int64, len(uids))
	for _, v := range uids {
		key := fmt.Sprintf(redisKeyHUserFollow, v)
		val, err := r.redis.HGetXEX(ctx, key, cast.FormatInt(uid), relationCacheL2TTL)
		if err != nil {
			missed = append(missed, v)
			continue
		} else {
			rm[v] = cast.ParseInt(val, 0)
		}
	}
	return rm, missed, nil
}

func (r *relationshipRepo) cacheGetRelationCount(ctx context.Context, uids []int64) (map[int64]*biz.RelationCount, []int64, error) {
	pipe := r.redis.Pipeline()
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
	fm := make(map[int64]*biz.RelationCount, len(uids))
	for k, v := range cmdMap {
		val, err := v.Result()
		if err != nil {
			missed = append(missed, k)
			continue
		}
		fc := &biz.RelationCount{}
		if c, ok := val[followField]; ok {
			fc.FollowCount = cast.Atoi(c, 0)
		}
		if c, ok := val[followerField]; ok {
			fc.FollowerCount = cast.Atoi(c, 0)
		}
	}
	return fm, missed, nil
}
