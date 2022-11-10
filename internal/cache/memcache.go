package cache

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/zxq97/gotool/cast"
	"github.com/zxq97/gotool/concurrent"
	"github.com/zxq97/relation/internal/constant"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/model"
)

const (
	relationCacheL1TTL = 8 * 3600
	mcKeyUserFollow    = "rla_fow_%d" // uid
	mcKeyUserFollower  = "rla_foe_%d" // uid
)

func getRelationCacheL1(ctx context.Context, keyPrefix string, uids []int64) (map[int64][]*model.FollowItem, []int64, error) {
	keys := make([]string, len(uids))
	for k, v := range uids {
		keys[k] = fmt.Sprintf(keyPrefix, v)
	}
	val, err := mcx.GetMultiCtx(ctx, keys)
	if err != nil {
		return nil, nil, err
	}
	itemMap := make(map[int64][]*model.FollowItem, len(uids))
	for k, v := range val {
		list := &model.FollowList{}
		err = proto.Unmarshal(v.Value, list)
		if err != nil {
			env.ExcLogger.Println()
			continue
		}
		s := strings.Split(k, "_")
		if len(s) == 0 {
			env.ExcLogger.Println()
			continue
		}
		itemMap[cast.ParseInt(s[len(s)-1], 0)] = list.List
	}
	missed := make([]int64, 0, len(uids))
	for _, k := range uids {
		if _, ok := itemMap[k]; !ok {
			missed = append(missed, k)
		}
	}
	return itemMap, missed, nil
}

func setRelationCacheL1(ctx context.Context, keyPrefix string, listMap map[int64]*model.FollowList) {
	for k, v := range listMap {
		key := fmt.Sprintf(keyPrefix, k)
		bs, err := proto.Marshal(v)
		if err != nil {
			env.ExcLogger.Println()
			continue
		}
		err = mcx.SetCtx(ctx, key, bs, relationCacheL1TTL)
		if err != nil {
			env.ExcLogger.Println()
		}
	}
}

func addRelationCacheL1(ctx context.Context, keyPrefix string, uid int64, itemList []*model.FollowItem) {
	key := fmt.Sprintf(keyPrefix, uid)
	val, err := mcx.GetCtx(ctx, key)
	if err != nil {
		env.ExcLogger.Println()
		return
	}
	list := model.FollowList{}
	err = proto.Unmarshal(val.Value, &list)
	if err != nil {
		env.ExcLogger.Println()
		return
	}
	for _, v := range itemList {
		list.List = append(list.List, v)
	}
	bs, err := proto.Marshal(&list)
	if err != nil {
		env.ExcLogger.Println()
		return
	}
	err = mcx.SetCtx(ctx, key, bs, relationCacheL1TTL)
	if err != nil {
		env.ExcLogger.Println()
	}
}

func delRelationCacheL1(ctx context.Context, keyPrefix string, uid, touid int64) {
	key := fmt.Sprintf(keyPrefix, uid)
	val, err := mcx.GetCtx(ctx, key)
	if err != nil {
		env.ExcLogger.Println()
		return
	}
	list := model.FollowList{}
	err = proto.Unmarshal(val.Value, &list)
	if err != nil {
		env.ExcLogger.Println()
		return
	}
	for k, v := range list.List {
		if v.ToUid == touid {
			list.List = append(list.List[0:k], list.List[k:]...)
		}
	}
	bs, err := proto.Marshal(&list)
	if err != nil {
		env.ExcLogger.Println()
		return
	}
	err = mcx.SetCtx(ctx, key, bs, relationCacheL1TTL)
	if err != nil {
		env.ExcLogger.Println()
	}
}

func getRelationList(ctx context.Context, uid, lastid, offset int64, follow bool) ([]*model.FollowItem, error) {
	mcKey := mcKeyUserFollow
	redisKey := redisKeyHUserFollow
	if !follow {
		mcKey = mcKeyUserFollower
		redisKey = redisKeyHUserFollower
	}
	listMap, _, err := getRelationCacheL1(ctx, mcKey, []int64{uid})
	list, ok := listMap[uid]
	if err != nil || !ok {
		listMap, _, err = getRelationCacheL2(ctx, redisKey, []int64{uid})
		list, ok = listMap[uid]
		if err != nil || !ok {
			return nil, err
		}
		concurrent.Go(func() {
			if follow {
				setRelationCacheL1(context.TODO(), mcKey, map[int64]*model.FollowList{uid: {List: list}})
			} else {
				addRelationCacheL1(context.TODO(), mcKey, uid, list)
			}
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].CreateTime > list[j].CreateTime
	})
	idx := sort.Search(len(list), func(i int) bool {
		return list[i].ToUid == lastid
	})
	if idx == len(list) {
		return nil, nil
	}
	right := idx + int(offset)
	if right > len(list) {
		right = len(list)
	}
	return list[idx+1 : right], nil
}

func AddRelation(ctx context.Context, uid int64, item *model.FollowItem) {
	addRelationCacheL1(ctx, mcKeyUserFollow, uid, []*model.FollowItem{item})
	addRelationCacheL2(ctx, redisKeyHUserFollow, uid, []*model.FollowItem{item})
	follower := item
	tuid := follower.ToUid
	follower.ToUid = uid
	addRelationCacheL1(ctx, mcKeyUserFollower, tuid, []*model.FollowItem{follower})
	addRelationCount(ctx, uid, item.ToUid, 1)
}

func AddUserFollower(ctx context.Context, uid int64, list []*model.FollowItem) {
	addRelationCacheL1(ctx, mcKeyUserFollower, uid, list)
}

func DelRelation(ctx context.Context, uid, touid int64) {
	delRelationCacheL1(ctx, mcKeyUserFollow, uid, touid)
	delRelationCacheL1(ctx, mcKeyUserFollower, touid, uid)
	delRelationCacheL2(ctx, redisKeyHUserFollow, uid, touid)
	addRelationCount(ctx, uid, touid, -1)
}

func GetFollowList(ctx context.Context, uid, lastid, offset int64) ([]*model.FollowItem, error) {
	return getRelationList(ctx, uid, lastid, offset, true)
}

func GetFollowerList(ctx context.Context, uid, lastid, offset int64) ([]*model.FollowItem, error) {
	list, err := getRelationList(ctx, uid, lastid, offset, false)
	if len(list) >= constant.ShowFollowerLimit {
		return nil, nil
	}
	return list, err
}
