package cache

import (
	"context"
	"fmt"
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/zxq97/gotool/concurrent"
	"github.com/zxq97/relation/internal/constant"
	"github.com/zxq97/relation/internal/model"
)

const (
	relationCacheL1TTL = 8 * 3600
	mcKeyUserFollow    = "rla_fow_%d" // uid
	mcKeyUserFollower  = "rla_foe_%d" // uid
)

func getRelationCacheL1(ctx context.Context, keyPrefix string, uid int64) ([]*model.FollowItem, error) {
	key := fmt.Sprintf(keyPrefix, uid)
	val, err := mcx.GetCtx(ctx, key)
	if err != nil {
		return nil, err
	}
	list := &model.FollowList{}
	err = proto.Unmarshal(val.Value, list)
	if err != nil {
		return nil, err
	}
	return list.List, nil
}

func setRelationCacheL1(ctx context.Context, keyPrefix string, uid int64, list []*model.FollowItem) error {
	val := &model.FollowList{List: list}
	key := fmt.Sprintf(keyPrefix, uid)
	bs, err := proto.Marshal(val)
	if err != nil {
		return err
	}
	err = mcx.SetCtx(ctx, key, bs, relationCacheL1TTL)
	return err
}

func addRelationCacheL1(ctx context.Context, keyPrefix string, uid int64, itemList []*model.FollowItem) error {
	key := fmt.Sprintf(keyPrefix, uid)
	val, err := mcx.GetCtx(ctx, key)
	if err != nil {
		return err
	}
	list := model.FollowList{}
	err = proto.Unmarshal(val.Value, &list)
	if err != nil {
		return err
	}
	for _, v := range itemList {
		list.List = append(list.List, v)
	}
	bs, err := proto.Marshal(&list)
	if err != nil {
		return err
	}
	err = mcx.SetCtx(ctx, key, bs, relationCacheL1TTL)
	return err
}

func delRelationCacheL1(ctx context.Context, keyPrefix string, uid, touid int64) error {
	key := fmt.Sprintf(keyPrefix, uid)
	val, err := mcx.GetCtx(ctx, key)
	if err != nil {
		return err
	}
	list := model.FollowList{}
	err = proto.Unmarshal(val.Value, &list)
	if err != nil {
		return err
	}
	for k, v := range list.List {
		if v.ToUid == touid {
			list.List = append(list.List[0:k], list.List[k:]...)
		}
	}
	bs, err := proto.Marshal(&list)
	if err != nil {
		return err
	}
	return mcx.SetCtx(ctx, key, bs, relationCacheL1TTL)
}

func getRelationList(ctx context.Context, uid, lastid, offset int64, follow bool) ([]*model.FollowItem, error) {
	mcKey := mcKeyUserFollow
	redisKey := redisKeyHUserFollow
	if !follow {
		mcKey = mcKeyUserFollower
	}
	list, err := getRelationCacheL1(ctx, mcKey, uid)
	if err != nil {
		if !follow {
			return nil, err
		}
		list, err = getRelationCacheL2(ctx, redisKey, uid)
		if err != nil {
			return nil, err
		}
		concurrent.Go(func() {
			_ = setRelationCacheL1(context.TODO(), mcKey, uid, list)
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
	_ = addRelationCacheL1(ctx, mcKeyUserFollow, uid, []*model.FollowItem{item})
	_ = addRelationCacheL2(ctx, redisKeyHUserFollow, uid, []*model.FollowItem{item})
	uid, item.ToUid = item.ToUid, uid
	_ = addRelationCacheL1(ctx, mcKeyUserFollower, uid, []*model.FollowItem{item})
	addRelationCount(ctx, uid, item.ToUid, 1)
}

func AddUserFollower(ctx context.Context, uid int64, list []*model.FollowItem) {
	_ = addRelationCacheL1(ctx, mcKeyUserFollower, uid, list)
}

func DelRelation(ctx context.Context, uid, touid int64) {
	_ = delRelationCacheL1(ctx, mcKeyUserFollow, uid, touid)
	_ = delRelationCacheL1(ctx, mcKeyUserFollower, touid, uid)
	_ = delRelationCacheL2(ctx, redisKeyHUserFollow, uid, touid)
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
