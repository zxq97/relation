package data

import (
	"context"
	"fmt"
	"sync"

	"github.com/zxq97/gotool/cast"
	"github.com/zxq97/gotool/concurrent"
	"github.com/zxq97/gotool/kafka"
	"golang.org/x/sync/singleflight"
	"upper.io/db.v3/lib/sqlbuilder"
)

const (
	sfKeyGetFollowList   = "follow_list_%d"      // uid
	sfKeyGetFollowerList = "follower_list_%d_%d" // uid lastid
)

var (
	sfg singleflight.Group
)

func sfGetUserFollow(ctx context.Context, sess sqlbuilder.Database, uid int64) ([]*FollowItem, error) {
	key := fmt.Sprintf(sfKeyGetFollowList, uid)
	val, err, _ := sfg.Do(key, func() (interface{}, error) {
		return getUserFollow(ctx, sess, uid)
	})
	if err != nil {
		return nil, err
	}
	list, ok := val.([]*FollowItem)
	if !ok {
		return nil, ErrNotFount
	}
	return list, nil
}

func sfGetUserFollower(ctx context.Context, sess sqlbuilder.Database, producer *kafka.Producer, uid, lastid int64) ([]*FollowItem, error) {
	key := fmt.Sprintf(sfKeyGetFollowerList, uid, lastid)
	val, err, _ := sfg.Do(key, func() (interface{}, error) {
		list, err := getUserFollower(ctx, sess, uid, lastid, 20)
		if err != nil {
			return nil, err
		}
		_ = sendKafkaMsg(ctx, producer, kafka.TopicRelationCacheRebuild, cast.FormatInt(uid), &RebuildKafka{Uid: uid, LastId: lastid}, kafka.EventTypeListMissed)
		return list, nil
	})
	if err != nil {
		return nil, err
	}
	list, ok := val.([]*FollowItem)
	if !ok {
		return nil, ErrNotFount
	}
	return list, nil
}

func sfGetUsersFollow(ctx context.Context, sess sqlbuilder.Database, uids []int64) (map[int64][]*FollowItem, error) {
	eg := concurrent.NewErrGroup(ctx)
	lock := sync.Mutex{}
	m := make(map[int64][]*FollowItem, len(uids))
	for _, v := range uids {
		u := v
		eg.Go(func() error {
			list, err := sfGetUserFollow(ctx, sess, u)
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
