package cache

import (
	"context"
	"fmt"
	"sort"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/zxq97/gotool/concurrent"
	"github.com/zxq97/relation/internal/model"
)

const (
	relationCacheL1TTL = 8 * 3600
	mcKeyUserFollow    = "rla_fow_%d" // uid
)

func getFollowsCacheL1(ctx context.Context, uids []int64) (map[int64][]*model.FollowItem, []int64, error) {
	keys := make([]string, len(uids))
	for k, v := range uids {
		keys[k] = fmt.Sprintf(mcKeyUserFollow, v)
	}
	val, err := mcx.GetMultiCtx(ctx, keys)
	if err != nil {
		return nil, nil, errors.Wrap(err, "mc get users follow")
	}
	m := make(map[int64][]*model.FollowItem, len(uids))
	missed := make([]int64, 0, len(uids))
	for i, k := range keys {
		if v, ok := val[k]; ok {
			list := &model.FollowList{}
			err = proto.Unmarshal(v.Value, list)
			if err != nil {
				missed = append(missed, uids[i])
				continue
			}
			m[uids[i]] = list.List
		} else {
			missed = append(missed, uids[i])
		}
	}
	return m, missed, nil
}

func getFollowCacheL1(ctx context.Context, uid int64) ([]*model.FollowItem, error) {
	m, _, err := getFollowsCacheL1(ctx, []int64{uid})
	if err != nil {
		return nil, err
	}
	list, ok := m[uid]
	if !ok {
		return nil, memcache.ErrCacheMiss
	}
	return list, nil
}

func setFollowCacheL1(ctx context.Context, uid int64, list []*model.FollowItem) error {
	val := &model.FollowList{List: list}
	key := fmt.Sprintf(mcKeyUserFollow, uid)
	bs, err := proto.Marshal(val)
	if err != nil {
		return err
	}
	err = mcx.SetCtx(ctx, key, bs, relationCacheL1TTL)
	return err
}

func addFollowCacheL1(ctx context.Context, uid int64, itemList []*model.FollowItem) error {
	key := fmt.Sprintf(mcKeyUserFollow, uid)
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
	val.Value = bs
	return mcx.CompareAndSwap(val)
}

func delFollowCacheL1(ctx context.Context, uid, touid int64) error {
	key := fmt.Sprintf(mcKeyUserFollow, uid)
	val, err := mcx.GetCtx(ctx, key)
	if err != nil {
		return err
	}
	list := model.FollowList{}
	err = proto.Unmarshal(val.Value, &list)
	if err != nil {
		return err
	}
	var flag bool
	for k, v := range list.List {
		if v.ToUid == touid {
			list.List = append(list.List[0:k], list.List[k:]...)
			flag = true
			break
		}
	}
	if !flag {
		return nil
	}
	bs, err := proto.Marshal(&list)
	if err != nil {
		return err
	}
	val.Value = bs
	return mcx.CompareAndSwap(val)
}

func AddRelation(ctx context.Context, uid int64, item *model.FollowItem) {
	_ = addFollowCacheL1(ctx, uid, []*model.FollowItem{item})
	_ = addFollowCacheL2(ctx, uid, []*model.FollowItem{item})
	_ = addFollower(ctx, uid, item)
	addRelationCount(ctx, uid, item.ToUid, 1)
}

func DelRelation(ctx context.Context, uid, touid int64) {
	_ = delFollowCacheL1(ctx, uid, touid)
	_ = delFollowCacheL2(ctx, uid, touid)
	_ = delFollower(ctx, uid, touid)
	addRelationCount(ctx, uid, touid, -1)
}

func GetFollowList(ctx context.Context, uid, lastid, offset int64) ([]*model.FollowItem, error) {
	list, err := getFollowCacheL1(ctx, uid)
	if err != nil {
		list, err = getFollowCacheL2(ctx, uid)
		if err != nil {
			return nil, err
		}
		concurrent.Go(func() {
			_ = setFollowCacheL1(context.TODO(), uid, list)
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

func GetFollowerList(ctx context.Context, uid, lastid, offset int64) ([]*model.FollowItem, error) {
	key := fmt.Sprintf(redisKeyZUserFollower, uid)
	zs, err := rdx.ZRevRangeByMemberWithScores(ctx, key, lastid, offset)
	if err != nil {
		return nil, err
	}
	list := make([]*model.FollowItem, len(zs))
	for k, z := range zs {
		list[k] = &model.FollowItem{
			ToUid:      z.Member.(int64),
			CreateTime: int64(z.Score),
		}
	}
	return list, nil
}

func GetUsersFollow(ctx context.Context, uids []int64) (map[int64][]*model.FollowItem, []int64, error) {
	return getFollowsCacheL1(ctx, uids)
}
