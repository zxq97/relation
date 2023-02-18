package biz

import (
	"context"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewRelationshipUseCase)

type RelationshipRepo interface {
	Follow(context.Context, int64, int64) error
	Unfollow(context.Context, int64, int64) error
	FollowerCacheRebuild(context.Context, int64, int64) error
}

type RelationshipUseCase struct {
	repo RelationshipRepo
}

func NewRelationshipUseCase(repo RelationshipRepo) *RelationshipUseCase {
	return &RelationshipUseCase{repo: repo}
}

func (uc *RelationshipUseCase) Follow(ctx context.Context, uid, touid int64) error {
	return uc.repo.Follow(ctx, uid, touid)
}

func (uc *RelationshipUseCase) Unfollow(ctx context.Context, uid, touid int64) error {
	return uc.repo.Unfollow(ctx, uid, touid)
}

func (uc *RelationshipUseCase) FollowerCacheRebuild(ctx context.Context, uid, lastid int64) error {
	return uc.repo.FollowerCacheRebuild(ctx, uid, lastid)
}
