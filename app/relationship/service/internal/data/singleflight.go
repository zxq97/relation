package data

import (
	"context"
	"fmt"
	"sync"

	"github.com/zxq97/gokit/pkg/cast"
	"github.com/zxq97/gokit/pkg/concurrent"
	"github.com/zxq97/gokit/pkg/mq"
	"github.com/zxq97/relation/app/relationship/pkg/bizdata"
	"github.com/zxq97/relation/app/relationship/pkg/message"
	"golang.org/x/sync/singleflight"
)

const (
	sfKeyGetFollowList   = "follow_list_%d"   // uid
	sfKeyGetFollowerList = "follower_list_%d" // uid
)

var (
	sfg singleflight.Group
)

func (r *relationshipRepo) sfGetUserFollow(ctx context.Context, uid int64) ([]*bizdata.FollowItem, error) {
	key := fmt.Sprintf(sfKeyGetFollowList, uid)
	val, err, _ := sfg.Do(key, func() (interface{}, error) {
		return r.getUserFollow(ctx, uid)
	})
	if err != nil {
		return nil, err
	}
	list, ok := val.([]*bizdata.FollowItem)
	if !ok {
		return nil, bizdata.ErrNotFound
	}
	return list, nil
}

func (r *relationshipRepo) sfGetUserFollower(ctx context.Context, uid, lastid int64) ([]*bizdata.FollowItem, error) {
	key := fmt.Sprintf(sfKeyGetFollowerList, uid)
	val, err, _ := sfg.Do(key, func() (interface{}, error) {
		list, err := r.getUserFollower(ctx, uid, lastid)
		if err != nil {
			return nil, err
		}
		_ = r.p.SendMessage(ctx, message.TopicRelationCacheRebuild, cast.FormatInt(uid), mq.TagListMissed, &message.CacheRebuild{Uid: uid, LastId: lastid})
		return list, nil
	})
	if err != nil {
		return nil, err
	}
	list, ok := val.([]*bizdata.FollowItem)
	if !ok {
		return nil, bizdata.ErrNotFound
	}
	return list, nil
}

func (r *relationshipRepo) sfGetUsersFollow(ctx context.Context, uids []int64) (map[int64][]*bizdata.FollowItem, error) {
	eg := concurrent.NewErrGroup(ctx)
	lock := sync.Mutex{}
	m := make(map[int64][]*bizdata.FollowItem, len(uids))
	for _, v := range uids {
		u := v
		eg.Go(func() error {
			list, err := r.sfGetUserFollow(ctx, u)
			if err != nil {
				return err
			}
			lock.Lock()
			defer lock.Unlock()
			m[u] = list
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	return m, nil
}
