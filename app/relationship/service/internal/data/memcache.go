package data

import (
	"context"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/zxq97/relation/app/relationship/service/internal/biz"
)

const (
	relationCacheL1TTL = 8 * 3600
	mcKeyUserFollow    = "rla_fow_%d" // uid
)

func (r *relationshipRepo) getFollowsCacheL1(ctx context.Context, uids []int64) (map[int64][]*biz.FollowItem, []int64, error) {
	keys := make([]string, len(uids))
	for k, v := range uids {
		keys[k] = fmt.Sprintf(mcKeyUserFollow, v)
	}
	val, err := r.mc.GetMulti(keys)
	if err != nil {
		return nil, nil, errors.Wrap(err, "mc get users follow")
	}
	m := make(map[int64][]*biz.FollowItem, len(uids))
	missed := make([]int64, 0, len(uids))
	for i, k := range keys {
		if v, ok := val[k]; ok {
			list := &biz.FollowList{}
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

func (r *relationshipRepo) getFollowCacheL1(ctx context.Context, uid int64) ([]*biz.FollowItem, error) {
	m, _, err := r.getFollowsCacheL1(ctx, []int64{uid})
	if err != nil {
		return nil, err
	}
	list, ok := m[uid]
	if !ok {
		return nil, memcache.ErrCacheMiss
	}
	return list, nil
}

func (r *relationshipRepo) setFollowCacheL1(ctx context.Context, uid int64, list []*biz.FollowItem) error {
	val := &biz.FollowList{List: list}
	key := fmt.Sprintf(mcKeyUserFollow, uid)
	bs, err := proto.Marshal(val)
	if err != nil {
		return err
	}
	err = r.mc.Set(&memcache.Item{Key: key, Value: bs, Expiration: relationCacheL1TTL})
	return err
}
