package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/zxq97/gokit/pkg/cast"
	"github.com/zxq97/relation/app/relationship/job/internal/biz"
)

const (
	relationCacheL2TTL     = time.Hour * 12
	redisKeyHUserFollow    = "rla_fow_%d" // uid
	redisKeyHRelationCount = "lra_cnt_%d" // uid
	redisKeyZUserFollower  = "rla_foe_%d" // uid

	followField   = "follow"
	followerField = "follower"
)

func (r *relationshipRepo) addRelationCount(ctx context.Context, uid, touid, incr int64) error {
	key := fmt.Sprintf(redisKeyHRelationCount, uid)
	err := r.redis.HIncrByXEX(ctx, key, followField, incr, relationCacheL2TTL)
	if err != nil {
		err = errors.WithMessage(err, "add follow count")
	}
	key = fmt.Sprintf(redisKeyHRelationCount, touid)
	err = r.redis.HIncrByXEX(ctx, key, followerField, incr, relationCacheL2TTL)
	return errors.WithMessage(err, "add follower count")
}

func (r *relationshipRepo) addFollower(ctx context.Context, uid int64, item *biz.FollowItem) error {
	key := fmt.Sprintf(redisKeyZUserFollower, item.ToUid)
	return r.redis.ZAddXEX(ctx, key, []*redis.Z{{Member: uid, Score: float64(item.CreateTime)}}, time.Hour)
}

func (r *relationshipRepo) appendFollower(ctx context.Context, uid int64, list []*biz.FollowItem) error {
	key := fmt.Sprintf(redisKeyZUserFollower, uid)
	zs := make([]*redis.Z, len(list))
	for k, v := range list {
		zs[k] = &redis.Z{
			Member: v.ToUid,
			Score:  float64(v.CreateTime),
		}
	}
	return r.redis.ZAddEX(ctx, key, zs, time.Hour)
}

func (r *relationshipRepo) addFollowCacheL2(ctx context.Context, uid int64, list []*biz.FollowItem) error {
	fieldMap := make(map[string]interface{}, len(list))
	for _, v := range list {
		fieldMap[cast.FormatInt(v.ToUid)] = v.CreateTime
	}
	return r.redis.HMSetXEX(ctx, fmt.Sprintf(redisKeyHUserFollow, uid), fieldMap, relationCacheL2TTL)
}

func (r *relationshipRepo) delFollowCacheL2(ctx context.Context, uid, touid int64) error {
	return r.redis.HDel(ctx, fmt.Sprintf(redisKeyHUserFollow, uid), cast.FormatInt(touid)).Err()
}

func (r *relationshipRepo) delFollower(ctx context.Context, uid, touid int64) error {
	key := fmt.Sprintf(redisKeyZUserFollower, touid)
	return r.redis.ZRem(ctx, key, uid).Err()
}
