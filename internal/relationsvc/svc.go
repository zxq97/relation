package relationsvc

import (
	"context"
	"fmt"
	"sort"

	"github.com/zxq97/gotool/cast"
	"github.com/zxq97/gotool/concurrent"
	"github.com/zxq97/gotool/kafka"
	"github.com/zxq97/relation/internal/cache"
	"github.com/zxq97/relation/internal/constant"
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

//Follow 关注
func (RelationSvc) Follow(ctx context.Context, req *FollowRequest) (*EmptyResponse, error) {
	bs, err := packKafkaMsg(ctx, req, kafka.EventTypeCreate)
	if err != nil {
		return &EmptyResponse{}, err
	}
	shardKey := cast.FormatInt(req.Uid)
	return &EmptyResponse{}, producer.SendMessage(kafka.TopicRelationFollow, []byte(shardKey), bs)
}

//Unfollow 取关
func (RelationSvc) Unfollow(ctx context.Context, req *FollowRequest) (*EmptyResponse, error) {
	bs, err := packKafkaMsg(ctx, req, kafka.EventTypeDelete)
	if err != nil {
		return &EmptyResponse{}, err
	}
	shardKey := cast.FormatInt(req.Uid)
	return &EmptyResponse{}, producer.SendMessage(kafka.TopicRelationOperator, []byte(shardKey), bs)
}

//GetFollowList 关注列表
func (RelationSvc) GetFollowList(ctx context.Context, req *ListRequest) (*model.FollowList, error) {
	val, err, _ := sfg.Do(fmt.Sprintf(sfKeyGetFollowList, req.Uid), func() (interface{}, error) {
		list, err := cache.GetFollowList(ctx, req.Uid, req.LastId, constant.ListBatchSize)
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
					right := idx + constant.ListBatchSize
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

//GetFollowerList 粉丝列表
func (RelationSvc) GetFollowerList(ctx context.Context, req *ListRequest) (*model.FollowList, error) {
	val, err, _ := sfg.Do(fmt.Sprintf(sfKeyGetFollowerList, req.Uid), func() (interface{}, error) {
		list, err := cache.GetFollowerList(ctx, req.Uid, req.LastId, constant.ListBatchSize)
		if err != nil {
			list, err = store.GetUserFollower(ctx, req.Uid, req.LastId, constant.ListBatchSize)
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

//GetRelation 好有关系
func (RelationSvc) GetRelation(ctx context.Context, req *RelationRequest) (*RelationResponse, error) {
	followMap, err := cache.IsFollows(ctx, req.Uid, req.Uids)
	if err != nil {
		list, err := getUserFollow(ctx, req.Uid)
		if err != nil {
			return &RelationResponse{}, nil
		}
		m := make(map[int64]int64, len(list))
		for _, v := range list {
			m[v.ToUid] = v.CreateTime
		}
		for _, v := range req.Uids {
			if t, ok := m[v]; ok {
				followMap[v] = t
			}
		}
	}
	followerMap, missed, err := cache.IsFollowers(ctx, req.Uid, req.Uids)
	if err != nil || len(missed) != 0 {
		dbm, err := getUsersFollow(ctx, missed)
		if err != nil {
			return &RelationResponse{}, err
		}
		for k, v := range dbm {
			for _, f := range v {
				if f.ToUid == req.Uid {
					followerMap[k] = f.CreateTime
					break
				}
			}
		}
	}
	relationMap := make(map[int64]*RelationItem, len(req.Uids))
	for k, v := range followMap {
		relationMap[k].Relation |= 1
		relationMap[k].FollowTime = v
	}
	for k, v := range followerMap {
		relationMap[k].Relation |= 2
		relationMap[k].FollowedTime = v
	}
	return &RelationResponse{Rm: relationMap}, nil
}

//GetRelationCount 关注 粉丝 数量
func (RelationSvc) GetRelationCount(ctx context.Context, req *BatchRequest) (*CountResponse, error) {
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
	return &CountResponse{RelationCount: fcDAO2DTO(fm)}, nil
}

//GetUsersFollow 获取全量关注
func (RelationSvc) GetUsersFollow(ctx context.Context, req *BatchRequest) (*UserFollowResponse, error) {
	m, missed, err := cache.GetUsersFollow(ctx, req.Uids)
	if err != nil || len(missed) != 0 {
		dbm, err := getUsersFollow(ctx, req.Uids)
		if err != nil {
			return &UserFollowResponse{}, err
		}
		for k, v := range dbm {
			m[k] = v
		}
	}
	return &UserFollowResponse{Fm: usDAO2DTO(m)}, nil
}
