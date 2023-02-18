package data

import (
	"context"
	"time"

	"github.com/google/wire"
	"github.com/zxq97/gokit/pkg/cast"
	"github.com/zxq97/gokit/pkg/mq/kafka"
	"github.com/zxq97/relation/app/relationship/job/internal/biz"
	"github.com/zxq97/relation/app/relationship/pkg/dal/cache"
	"github.com/zxq97/relation/app/relationship/pkg/dal/query"
	"github.com/zxq97/relation/app/relationship/pkg/message"
)

var ProviderSet = wire.NewSet(NewRelationshipRepo)
var _ biz.RelationshipRepo = (*relationshipRepo)(nil)

type relationshipRepo struct {
	p *kafka.Producer
	c *cache.Cache
	q *query.Query
}

func NewRelationshipRepo(p *kafka.Producer, c *cache.Cache, q *query.Query) *relationshipRepo {
	return &relationshipRepo{p: p, c: c, q: q}
}

func (r *relationshipRepo) Follow(ctx context.Context, uid, touid int64) error {
	if err := r.follow(ctx, uid, touid); err != nil {
		return err
	}
	if ok := r.c.BatchSyncCount(ctx, touid); ok {
		_ = r.p.SendMessage(ctx, message.TopicRelationSyncCount, cast.FormatInt(touid), message.TagSync, &message.SyncCount{Uid: touid, TimeWait: time.Now().Add(time.Second).UnixMilli()})
	}
	return r.c.Follow(ctx, uid, touid)
}

func (r *relationshipRepo) Unfollow(ctx context.Context, uid, touid int64) error {
	if err := r.unfollow(ctx, uid, touid); err != nil {
		return err
	}
	if ok := r.c.BatchSyncCount(ctx, touid); ok {
		_ = r.p.SendMessage(ctx, message.TopicRelationSyncCount, cast.FormatInt(touid), message.TagSync, &message.SyncCount{Uid: touid, TimeWait: time.Now().Add(time.Second).UnixMilli()})
	}
	return r.c.Unfollow(ctx, uid, touid)
}

func (r *relationshipRepo) FollowerCacheRebuild(ctx context.Context, uid, lastid int64) error {
	list, err := r.getUserFollower(ctx, uid, lastid)
	if err != nil {
		return err
	}
	return r.c.AppendFollower(ctx, uid, list)
}
