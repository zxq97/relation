package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zxq97/gotool/cast"
	"github.com/zxq97/relation/internal/model"
)

const (
	relationCacheL2TTL    = time.Hour * 12
	redisKeyHUserFollow   = "rla_fow_%d"
	redisKeyHUserFollower = "rla_foe_%d"
)

func getRelationCacheL2(ctx context.Context, keyPrefix string, uids []int64) (map[int64][]*model.FollowItem, []int64, error) {
	cmdMap := make(map[int64]*redis.StringStringMapCmd, len(uids))
	pipe := rdx.Pipeline()
	for _, v := range uids {
		cmdMap[v] = pipe.HGetAll(ctx, fmt.Sprintf(keyPrefix, v))
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, nil, err
	}
	itemMap := make(map[int64][]*model.FollowItem, len(uids))
	for k, v := range cmdMap {
		val, err := v.Result()
		if err != nil {
			// log
			continue
		}
		itemMap[k] = make([]*model.FollowItem, len(val))
		for u, t := range val {
			itemMap[k] = append(itemMap[k], &model.FollowItem{ToUid: cast.ParseInt(u, 0), CreateTime: cast.ParseInt(t, 0)})
		}
	}
	missed := make([]int64, 0, len(uids))
	for _, k := range uids {
		if _, ok := itemMap[k]; !ok {
			missed = append(missed, k)
		}
	}
	return itemMap, missed, nil
}

func setRelationCacheL2(ctx context.Context, keyPrefix string, listMap map[int64]*model.FollowList) {
	pipe := rdx.Pipeline()
	for k, v := range listMap {
		field := make(map[string]interface{}, len(listMap))
		for _, t := range v.List {
			field[cast.FormatInt(t.ToUid)] = t.CreateTime
		}
		key := fmt.Sprintf(keyPrefix, k)
		pipe.HMSet(ctx, key, field)
		pipe.Expire(ctx, key, relationCacheL2TTL)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		// log
	}
}

func addRelationCacheL2(ctx context.Context, keyPrefix string, uid int64, itemList []*model.FollowItem) {
	pipe := rdx.Pipeline()
	field := make(map[string]interface{}, len(itemList))
	key := fmt.Sprintf(keyPrefix, uid)
	pipe.HMSet(ctx, key, field)
	pipe.Expire(ctx, key, relationCacheL2TTL)
	_, err := pipe.Exec(ctx)
	if err != nil {
		// log
	}
}

func SetFollowList(ctx context.Context, uid int64, list []*model.FollowItem) {
	listMap := map[int64]*model.FollowList{uid: {List: list}}
	setRelationCacheL1(ctx, mcKeyUserFollow, listMap)
	setRelationCacheL2(ctx, redisKeyHUserFollow, listMap)
}

func IsFollows(ctx context.Context, uid int64, uids []int64) (map[int64]struct{}, error) {
	key := fmt.Sprintf(redisKeyHUserFollow, uid)
	us := make([]string, len(uids))
	for k, v := range uids {
		us[k] = cast.FormatInt(v)
	}
	val, err := rdx.HMGetEX(ctx, key, relationCacheL2TTL, us...)
	if err != nil {
		return nil, err
	}
	rm := make(map[int64]struct{}, len(uids))
	for k, v := range val {
		if v != nil {
			rm[uids[k]] = struct{}{}
		}
	}
	return rm, nil
}

func IsFollowers(ctx context.Context, uid int64, uids []int64) (map[int64]struct{}, []int64, error) {
	missed := make([]int64, 0, len(uids))
	rm := make(map[int64]struct{}, len(uids))
	for _, v := range uids {
		key := fmt.Sprintf(redisKeyHUserFollower, v)
		val, err := rdx.HMGetEX(ctx, key, relationCacheL2TTL, cast.FormatInt(uid))
		if err != nil {
			missed = append(missed, v)
			continue
		} else if val != nil {
			rm[v] = struct{}{}
		}
	}
	return rm, missed, nil
}
