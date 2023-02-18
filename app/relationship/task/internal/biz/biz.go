package biz

import (
	"context"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewRelationshipUseCase)

type RelationshipRepo interface {
	SyncRecord(context.Context, int64) error
	SyncRecordByUID(context.Context, int64) error
}

type RelationshipUseCase struct {
	repo RelationshipRepo
}

func NewRelationshipUseCase(repo RelationshipRepo) *RelationshipUseCase {
	return &RelationshipUseCase{repo: repo}
}

func (uc *RelationshipUseCase) CronTaskSyncRecord(ctx context.Context, limit int64) error {
	return uc.repo.SyncRecord(ctx, limit)
}

func (uc *RelationshipUseCase) SyncRecordByUID(ctx context.Context, uid int64) error {
	return uc.repo.SyncRecordByUID(ctx, uid)
}
