package data

import (
	"context"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/zxq97/relation/app/relationship/pkg/dal/cache"
	"github.com/zxq97/relation/app/relationship/pkg/dal/query"
	"github.com/zxq97/relation/app/relationship/task/internal/biz"
)

var ProviderSet = wire.NewSet(NewRelationshipRepo)
var _ biz.RelationshipRepo = (*relationshipRepo)(nil)

type relationshipRepo struct {
	q *query.Query
	c *cache.Cache
}

func NewRelationshipRepo(c *cache.Cache, q *query.Query) *relationshipRepo {
	return &relationshipRepo{c: c, q: q}
}

func (r *relationshipRepo) SyncRecord(ctx context.Context, limit int64) error {
	m, err := r.syncRecord(ctx, limit)
	if err != nil {
		return err
	}
	for k, v := range m {
		if err = r.c.IncrFollowerCount(ctx, k, v); err != nil {
			err = errors.WithMessage(err, "incr follower")
		}
	}
	return err
}

func (r *relationshipRepo) SyncRecordByUID(ctx context.Context, uid int64) error {
	cnt, err := r.syncRecordByUID(ctx, uid)
	if err != nil {
		return err
	}
	return r.c.IncrFollowerCount(ctx, uid, cnt)
}
