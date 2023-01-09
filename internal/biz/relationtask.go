package biz

import (
	"context"

	"github.com/zxq97/gotool/config"
	"github.com/zxq97/relation/internal/data"
)

type RelationTaskRepo interface {
	Follow(context.Context, int64, int64) error
	Unfollow(context.Context, int64, int64) error
	FollowerCacheRebuild(context.Context, int64, int64) error
}

type RelationTaskBIZ struct {
	repo RelationTaskRepo
}

func NewRelationTaskBIZ(redisConf *config.RedisConf, mcConf *config.MCConf, mysqlConf *config.MysqlConf) (*RelationTaskBIZ, error) {
	repo, err := data.NewRelationTaskRepoImpl(redisConf, mcConf, mysqlConf)
	if err != nil {
		return nil, err
	}
	return &RelationTaskBIZ{
		repo: repo,
	}, nil
}

func (rtb *RelationTaskBIZ) Follow(ctx context.Context, uid, touid int64) error {
	return rtb.repo.Follow(ctx, uid, touid)
}

func (rtb *RelationTaskBIZ) Unfollow(ctx context.Context, uid, touid int64) error {
	return rtb.repo.Unfollow(ctx, uid, touid)
}

func (rtb *RelationTaskBIZ) FollowerCacheRebuild(ctx context.Context, uid, lastid int64) error {
	return rtb.repo.FollowerCacheRebuild(ctx, uid, lastid)
}
