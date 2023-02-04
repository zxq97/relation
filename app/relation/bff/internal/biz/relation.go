package biz

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

var (
	ErrFollowLimit = errors.New("relation: follow limit")
)

type RelationRepo interface {
	Follow(context.Context, int64, int64) error
	Unfollow(context.Context, int64, int64) error
	GetRelationCount(context.Context, []int64) (map[int64]*RelationCount, error)
	GetUsersFollow(context.Context, []int64) (map[int64][]*FollowItem, error)
	GetUserFollow(context.Context, int64, int64) ([]*FollowItem, error)
	GetUserFollower(context.Context, int64, int64) ([]*FollowItem, error)
	GetRelation(context.Context, int64, []int64) (map[int64]*UserRelation, error)
}

type RelationUseCase struct {
	repo RelationRepo
}

func NewRelationUseCase(repo RelationRepo) *RelationUseCase {
	return &RelationUseCase{repo: repo}
}

func (uc *RelationUseCase) Follow(ctx context.Context, uid, touid int64) error {
	return uc.repo.Follow(ctx, uid, touid)
}

func (uc *RelationUseCase) Unfollow(ctx context.Context, uid, touid int64) error {
	return uc.repo.Unfollow(ctx, uid, touid)
}

func (uc *RelationUseCase) GetFollowList(ctx context.Context, uid, lastid int64) ([]*FollowItem, error) {
	return uc.repo.GetUserFollow(ctx, uid, lastid)
}

func (uc *RelationUseCase) GetFollowerList(ctx context.Context, uid, lastid int64) ([]*FollowItem, error) {
	return uc.repo.GetUserFollower(ctx, uid, lastid)
}

func (uc *RelationUseCase) GetRelation(ctx context.Context, uid int64, uids []int64) (map[int64]*UserRelation, error) {
	return uc.repo.GetRelation(ctx, uid, uids)
}

func (uc *RelationUseCase) GetRelationCount(ctx context.Context, uids []int64) (map[int64]*RelationCount, error) {
	return uc.repo.GetRelationCount(ctx, uids)
}

func (uc *RelationUseCase) GetCommonRelation(ctx context.Context, uid, touid int64) ([]int64, error) {
	m, err := uc.repo.GetUsersFollow(ctx, []int64{uid, touid})
	if err != nil {
		return nil, err
	}
	uf, ok := m[uid]
	if !ok {
		return nil, nil
	}
	tf, ok := m[touid]
	if !ok {
		return nil, nil
	}
	set := make(map[int64]struct{}, len(tf))
	for _, v := range tf {
		set[v.ToUID] = struct{}{}
	}
	cf := make([]int64, 0, len(uf))
	for _, v := range uf {
		if _, ok = set[v.ToUID]; ok {
			cf = append(cf, v.ToUID)
		}
	}
	return cf, nil
}

func (uc *RelationUseCase) GetRelationChain(ctx context.Context, uid, touid int64) ([]int64, error) {
	m, err := uc.repo.GetUsersFollow(ctx, []int64{uid})
	if err != nil {
		return nil, err
	}
	uf, ok := m[uid]
	if !ok {
		return nil, nil
	}
	list := make([]int64, len(uf))
	for k, v := range uf {
		list[k] = v.ToUID
	}
	cf := make([]int64, 0, 10)
	lock := sync.Mutex{}
	for i := 0; i < len(list); i += 10 {
		left := i
		right := i + 10
		if right > len(list) {
			right = len(list)
		}
		m, err = uc.repo.GetUsersFollow(ctx, list[left:right])
		if err != nil {
			return nil, err
		}
		wg := sync.WaitGroup{}
		for k, v := range m {
			wg.Add(1)
			u, l := k, v
			go func() {
				defer wg.Done()
				for j := range l {
					if l[j].ToUID == touid {
						lock.Lock()
						cf = append(cf, u)
						lock.Unlock()
						break
					}
				}
			}()
			wg.Wait()
		}
		if len(cf) >= 10 {
			break
		}
	}
	return cf, nil
}
