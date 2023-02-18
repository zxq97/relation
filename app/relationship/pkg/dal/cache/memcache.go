package cache

import (
	"context"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/zxq97/relation/app/relationship/pkg/bizdata"
)

const (
	relationCacheL1TTL = 8 * 3600
	mcKeyUserFollow    = "rla_fow_%d" // uid
)

func (c *Cache) getFollowsCacheL1(ctx context.Context, uids []int64) (map[int64][]*bizdata.FollowItem, []int64, error) {
	keys := make([]string, len(uids))
	for k, v := range uids {
		keys[k] = fmt.Sprintf(mcKeyUserFollow, v)
	}
	val, err := c.mc.GetMulti(keys)
	if err != nil {
		return nil, nil, errors.Wrap(err, "mc get users follow")
	}
	m := make(map[int64][]*bizdata.FollowItem, len(uids))
	missed := make([]int64, 0, len(uids))
	for i, k := range keys {
		if v, ok := val[k]; ok {
			list := &bizdata.FollowList{}
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

func (c *Cache) setFollowCacheL1(ctx context.Context, uid int64, list []*bizdata.FollowItem) error {
	val := &bizdata.FollowList{List: list}
	key := fmt.Sprintf(mcKeyUserFollow, uid)
	bs, err := proto.Marshal(val)
	if err != nil {
		return err
	}
	err = c.mc.Set(&memcache.Item{Key: key, Value: bs, Expiration: relationCacheL1TTL})
	return err
}

func (c *Cache) addFollowCacheL1(ctx context.Context, uid int64, itemList []*bizdata.FollowItem) error {
	key := fmt.Sprintf(mcKeyUserFollow, uid)
	val, err := c.mc.Get(key)
	if err != nil {
		return err
	}
	list := bizdata.FollowList{}
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
	return c.mc.CompareAndSwap(val)
}

func (c *Cache) delFollowCacheL1(ctx context.Context, uid, touid int64) error {
	key := fmt.Sprintf(mcKeyUserFollow, uid)
	val, err := c.mc.Get(key)
	if err != nil {
		return err
	}
	list := bizdata.FollowList{}
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
	return c.mc.CompareAndSwap(val)
}
