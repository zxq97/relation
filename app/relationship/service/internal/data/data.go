package data

import (
	"context"

	"github.com/google/wire"
	"github.com/zxq97/gokit/pkg/cast"
	"github.com/zxq97/gokit/pkg/concurrent"
	"github.com/zxq97/gokit/pkg/mq"
	"github.com/zxq97/gokit/pkg/mq/kafka"
	"github.com/zxq97/relation/app/relationship/pkg/bizdata"
	"github.com/zxq97/relation/app/relationship/pkg/dal/cache"
	"github.com/zxq97/relation/app/relationship/pkg/dal/query"
	"github.com/zxq97/relation/app/relationship/pkg/message"
	"github.com/zxq97/relation/app/relationship/service/internal/biz"
)

var ProviderSet = wire.NewSet(NewRelationshipRepo)
var _ biz.RelationshipRepo = (*relationshipRepo)(nil)

type relationshipRepo struct {
	p *kafka.Producer
	c *cache.Cache
	q *query.Query
}

func NewRelationshipRepo(p *kafka.Producer, c *cache.Cache, q *query.Query) biz.RelationshipRepo {
	return &relationshipRepo{p: p, c: c, q: q}
}

func (r *relationshipRepo) Follow(ctx context.Context, uid, touid int64) error {
	return r.p.SendMessage(ctx, message.TopicRelationFollow, cast.FormatInt(uid), mq.TagCreate, &message.AsyncFollow{Uid: uid, ToUid: touid})
}

func (r *relationshipRepo) Unfollow(ctx context.Context, uid, touid int64) error {
	return r.p.SendMessage(ctx, message.TopicRelationFollow, cast.FormatInt(uid), mq.TagDelete, &message.AsyncFollow{Uid: uid, ToUid: touid})
}

func (r *relationshipRepo) GetRelationCount(ctx context.Context, uids []int64) (map[int64]*bizdata.RelationCount, error) {
	m, missed, err := r.c.GetRelationCount(ctx, uids)
	if err != nil || len(missed) != 0 {
		dbm, err := r.getRelationCount(ctx, missed)
		if err != nil {
			return nil, err
		}
		for k, v := range dbm {
			m[k] = v
		}
		_ = r.c.SetRelationCount(ctx, dbm)
	}
	return m, nil
}

func (r *relationshipRepo) GetUsersFollow(ctx context.Context, uids []int64) (map[int64][]*bizdata.FollowItem, error) {
	m, missed, _ := r.c.GetUsersFollow(ctx, uids)
	if len(missed) != 0 {
		dbm, err := r.sfGetUsersFollow(ctx, missed)
		if err != nil {
			return nil, err
		}
		for k, v := range dbm {
			m[k] = v
		}
		concurrent.Go(func() {
			_ = r.c.SetUsersFollow(ctx, dbm)
		})
	}
	return m, nil
}

func (r *relationshipRepo) GetUserFollower(ctx context.Context, uid int64, lastid int64) ([]*bizdata.FollowItem, error) {
	list, err := r.c.GetFollowerList(ctx, uid, lastid)
	if err != nil {
		list, err = r.sfGetUserFollower(ctx, uid, lastid)
		if err != nil {
			return nil, err
		}
	}
	return list, nil
}

func (r *relationshipRepo) GetIsFollowMap(ctx context.Context, uid int64, uids []int64) (map[int64]int64, error) {
	followMap, err := r.c.IsFollow(ctx, uid, uids)
	if err != nil {
		list, err := r.sfGetUserFollow(ctx, uid)
		if err != nil {
			return nil, err
		}
		m := make(map[int64]int64, len(list))
		for _, v := range list {
			m[v.ToUid] = v.CreateTime
		}
		for _, v := range uids {
			if t, ok := m[v]; ok {
				followMap[v] = t
			}
		}
	}
	return followMap, nil
}

func (r *relationshipRepo) GetIsFollowerMap(ctx context.Context, uid int64, uids []int64) (map[int64]int64, error) {
	followerMap, missed, err := r.c.IsFollower(ctx, uid, uids)
	if err != nil || len(missed) != 0 {
		dbm, err := r.sfGetUsersFollow(ctx, missed)
		if err != nil {
			return nil, err
		}
		for k, v := range dbm {
			for _, f := range v {
				if f.ToUid == uid {
					followerMap[k] = f.CreateTime
					break
				}
			}
		}
	}
	return followerMap, nil
}
