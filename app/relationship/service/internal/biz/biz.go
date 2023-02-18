package biz

import (
	"context"
	"sort"

	"github.com/zxq97/relation/app/relationship/pkg/bizdata"
)

const (
	relationFollowBit   = 1
	relationFollowedBit = 2
)

type RelationshipRepo interface {
	Follow(context.Context, int64, int64) error
	Unfollow(context.Context, int64, int64) error
	GetRelationCount(context.Context, []int64) (map[int64]*bizdata.RelationCount, error)
	GetUsersFollow(context.Context, []int64) (map[int64][]*bizdata.FollowItem, error)
	GetUserFollower(context.Context, int64, int64) ([]*bizdata.FollowItem, error)
	GetIsFollowMap(context.Context, int64, []int64) (map[int64]int64, error)
	GetIsFollowerMap(context.Context, int64, []int64) (map[int64]int64, error)
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

func (uc *RelationshipUseCase) GetFollowList(ctx context.Context, uid, lastid int64) ([]*bizdata.FollowItem, error) {
	m, err := uc.repo.GetUsersFollow(ctx, []int64{uid})
	if err != nil {
		return nil, err
	}
	list, ok := m[uid]
	if !ok {
		return nil, bizdata.ErrNotFound
	}
	idx := sort.Search(len(list), func(i int) bool {
		return list[i].ToUid == lastid
	})
	if idx != len(list) {
		right := idx + 20
		if right > len(list) {
			right = len(list)
		}
		list = list[idx:right]
	}
	return list, nil
}

func (uc *RelationshipUseCase) GetFollowerList(ctx context.Context, uid, lastid int64) ([]*bizdata.FollowItem, error) {
	return uc.repo.GetUserFollower(ctx, uid, lastid)
}

func (uc *RelationshipUseCase) GetRelation(ctx context.Context, uid int64, uids []int64) (map[int64]*bizdata.UserRelation, error) {
	followMap, err := uc.repo.GetIsFollowMap(ctx, uid, uids)
	if err != nil {
		return nil, err
	}
	followerMap, err := uc.repo.GetIsFollowerMap(ctx, uid, uids)
	if err != nil {
		return nil, err
	}
	relationMap := make(map[int64]*bizdata.UserRelation, len(uids))
	for k, v := range followMap {
		relationMap[k].Relation |= relationFollowBit
		relationMap[k].FollowTime = v
	}
	for k, v := range followerMap {
		relationMap[k].Relation |= relationFollowedBit
		relationMap[k].FollowedTime = v
	}
	return relationMap, nil
}

func (uc *RelationshipUseCase) GetRelationCount(ctx context.Context, uids []int64) (map[int64]*bizdata.RelationCount, error) {
	return uc.repo.GetRelationCount(ctx, uids)
}

func (uc *RelationshipUseCase) GetUsersFollow(ctx context.Context, uids []int64) (map[int64][]*bizdata.FollowItem, error) {
	return uc.repo.GetUsersFollow(ctx, uids)
}
