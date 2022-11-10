package relationsvc

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/zxq97/gotool/cast"
	"github.com/zxq97/gotool/concurrent"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/relation/internal/cache"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/model"
	"github.com/zxq97/relation/internal/store"
	"golang.org/x/sync/singleflight"
)

const (
	sfKeyGetFollowList   = "follow_list_%d"   // uid
	sfKeyGetFollowerList = "follower_list_%d" // uid
)

var (
	producer *kafka.Producer
	sfg      singleflight.Group
)

type RelationSvc struct {
}

func InitRelationSvc(conf *RelationSvcConfig) error {
	err := env.InitLog(conf.LogPath)
	if err != nil {
		return err
	}
	cache.InitCache(conf.Redis["redis"], conf.MC["mc"])
	err = store.InitStore(conf.Mysql["relation"])
	if err != nil {
		return err
	}
	producer, err = kafka.InitKafkaProducer(conf.Kafka["kafka"].Addr, env.ApiLogger, env.ExcLogger)
	return err
}

func (RelationSvc) Follow(ctx context.Context, req *FollowRequest) (*EmptyResponse, error) {
	bs, err := packKafkaMsg(ctx, req, kafka.EventTypeCreate)
	if err != nil {
		return &EmptyResponse{}, err
	}
	shardKey := cast.FormatInt(req.Uid)
	return &EmptyResponse{}, producer.SendMessage(kafka.TopicRelationFollow, []byte(shardKey), bs)
}

func (RelationSvc) Unfollow(ctx context.Context, req *FollowRequest) (*EmptyResponse, error) {
	bs, err := packKafkaMsg(ctx, req, kafka.EventTypeDelete)
	if err != nil {
		return &EmptyResponse{}, err
	}
	shardKey := cast.FormatInt(req.Uid)
	return &EmptyResponse{}, producer.SendMessage(kafka.TopicRelationOperator, []byte(shardKey), bs)
}

func (RelationSvc) GetFollowList(ctx context.Context, req *ListRequest) (*model.FollowList, error) {
	val, err, _ := sfg.Do(fmt.Sprintf(sfKeyGetFollowList, req.Uid), func() (interface{}, error) {
		list, err := cache.GetFollowList(ctx, req.Uid, req.LastId, 20)
		if err != nil {
			list, err = store.GetAllUserFollow(ctx, req.Uid)
			if err != nil {
				return nil, nil
			}
			if len(list) != 0 {
				concurrent.Go(func() {
					cache.SetFollowList(context.TODO(), req.Uid, list)
				})
				idx := sort.Search(len(list), func(i int) bool {
					return list[i].ToUid == req.LastId
				})
				if idx != len(list) {
					right := idx + 20
					if right > len(list) {
						right = len(list)
					}
					list = list[idx:right]
				}
			}
		}
		return list, nil
	})
	res, ok := val.(*model.FollowList)
	if err != nil || !ok {
		return &model.FollowList{}, nil
	}
	return &model.FollowList{List: res.List}, nil
}

func (RelationSvc) GetFollowerList(ctx context.Context, req *ListRequest) (*model.FollowList, error) {
	val, err, _ := sfg.Do(fmt.Sprintf(sfKeyGetFollowerList, req.Uid), func() (interface{}, error) {
		list, err := cache.GetFollowerList(ctx, req.Uid, req.LastId, 20)
		if err != nil {
			list, err = store.GetUserFollower(ctx, req.Uid, req.LastId, 20)
			if err != nil {
				return nil, nil
			}
			bs, err := packKafkaMsg(ctx, req, kafka.EventTypeListMissed)
			if err != nil {
				return list, err
			}
			shardKey := cast.FormatInt(req.Uid)
			return list, producer.SendMessage(kafka.TopicRelationCacheRebuild, []byte(shardKey), bs)
		}
		return list, nil
	})
	res, ok := val.(*model.FollowList)
	if err != nil || !ok {
		return &model.FollowList{}, nil
	}
	return &model.FollowList{List: res.List}, nil
}

func (RelationSvc) GetRelation(ctx context.Context, req *RelationRequest) (*RelationResponse, error) {
	followMap, err := cache.IsFollows(ctx, req.Uid, req.Uids)
	if err != nil {
		list, err := getUserFollow(ctx, req.Uid)
		if err != nil {
			return &RelationResponse{}, nil
		}
		m := make(map[int64]struct{})
		for _, v := range list {
			m[v.ToUid] = struct{}{}
		}
		for _, v := range req.Uids {
			if _, ok := m[v]; ok {
				followMap[v] = struct{}{}
			}
		}
	}
	followerMap, missed, err := cache.IsFollowers(ctx, req.Uid, req.Uids)
	if err != nil || len(missed) != 0 {
		wg := concurrent.WaitGroup{}
		lock := sync.Mutex{}
		for _, v := range missed {
			wg.Go(func() {
				list, err := getUserFollow(ctx, v)
				if err != nil {
					return
				}
				for _, u := range list {
					if u.ToUid == v {
						lock.Lock()
						followerMap[v] = struct{}{}
						lock.Unlock()
					}
				}
			})
		}
	}
	relationMap := make(map[int64]int32, len(req.Uids))
	for k := range followMap {
		relationMap[k] |= 1
	}
	for k := range followerMap {
		relationMap[k] |= 2
	}
	return &RelationResponse{Rm: relationMap}, nil
}

func (RelationSvc) GetRelationCount(ctx context.Context, req *CountRequest) (*CountResponse, error) {
	fm, missed, err := cache.GetRelationCount(ctx, req.Uids)
	if err != nil || len(missed) > 0 {
		dbm, err := store.GetUsersFollowCount(ctx, req.Uids)
		if err != nil {
			// log
			return &CountResponse{}, err
		}
		for k, v := range dbm {
			fm[k] = v
		}
		// 缓存零关注数
		for _, v := range req.Uids {
			if _, ok := fm[v]; !ok {
				dbm[v] = &model.UserFollowCount{}
			}
		}
		cache.SetRelationCount(ctx, dbm)
	}
	return &CountResponse{RelationCount: model.FcDAO2DTO(fm)}, nil
}
