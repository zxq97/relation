package data

import (
	"context"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/zxq97/gotool/memcachex"
)

const (
	relationCacheL1TTL = 8 * 3600
	mcKeyUserFollow    = "rla_fow_%d" // uid
)

func getFollowsCacheL1(ctx context.Context, mcx *memcachex.MemcacheX, uids []int64) (map[int64][]*FollowItem, []int64, error) {
	keys := make([]string, len(uids))
	for k, v := range uids {
		keys[k] = fmt.Sprintf(mcKeyUserFollow, v)
	}
	val, err := mcx.GetMultiCtx(ctx, keys)
	if err != nil {
		return nil, nil, errors.Wrap(err, "mc get users follow")
	}
	m := make(map[int64][]*FollowItem, len(uids))
	missed := make([]int64, 0, len(uids))
	for i, k := range keys {
		if v, ok := val[k]; ok {
			list := &FollowList{}
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

func getFollowCacheL1(ctx context.Context, mcx *memcachex.MemcacheX, uid int64) ([]*FollowItem, error) {
	m, _, err := getFollowsCacheL1(ctx, mcx, []int64{uid})
	if err != nil {
		return nil, err
	}
	list, ok := m[uid]
	if !ok {
		return nil, memcache.ErrCacheMiss
	}
	return list, nil
}

func setFollowCacheL1(ctx context.Context, mcx *memcachex.MemcacheX, uid int64, list []*FollowItem) error {
	val := &FollowList{List: list}
	key := fmt.Sprintf(mcKeyUserFollow, uid)
	bs, err := proto.Marshal(val)
	if err != nil {
		return err
	}
	err = mcx.SetCtx(ctx, key, bs, relationCacheL1TTL)
	return err
}

func addFollowCacheL1(ctx context.Context, mcx *memcachex.MemcacheX, uid int64, itemList []*FollowItem) error {
	key := fmt.Sprintf(mcKeyUserFollow, uid)
	val, err := mcx.GetCtx(ctx, key)
	if err != nil {
		return err
	}
	list := FollowList{}
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

func delFollowCacheL1(ctx context.Context, mcx *memcachex.MemcacheX, uid, touid int64) error {
	key := fmt.Sprintf(mcKeyUserFollow, uid)
	val, err := mcx.GetCtx(ctx, key)
	if err != nil {
		return err
	}
	list := FollowList{}
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
