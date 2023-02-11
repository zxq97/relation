package data

import (
	"context"

	blackv1 "github.com/zxq97/relation/api/black/service/v1"
	relationshipv1 "github.com/zxq97/relation/api/relationship/service/v1"
	"github.com/zxq97/relation/app/relation/bff/internal/biz"
)

var _ biz.RelationRepo = (*relationRepo)(nil)

type relationRepo struct {
	relationshipClient relationshipv1.RelationSvcClient
	blackClient        blackv1.BlackSvcClient
}

func NewRelationRepo(relationshipClient relationshipv1.RelationSvcClient, blackClient blackv1.BlackSvcClient) *relationRepo {
	return &relationRepo{
		relationshipClient: relationshipClient,
		blackClient:        blackClient,
	}
}

func (r *relationRepo) Follow(ctx context.Context, uid, touid int64) error {
	res, err := r.blackClient.CheckBlacked(ctx, &blackv1.CheckRequest{Uid: uid, Uids: []int64{touid}})
	if err != nil {
		return err
	}
	if _, ok := res.BlackMap[touid]; ok {
		return biz.ErrBlacked
	}
	res, err = r.blackClient.CheckBlacked(ctx, &blackv1.CheckRequest{Uid: touid, Uids: []int64{uid}})
	if err != nil {
		return err
	}
	if _, ok := res.BlackMap[uid]; ok {
		return biz.ErrBlacked
	}
	cm, err := r.relationshipClient.GetRelationCount(ctx, &relationshipv1.BatchRequest{Uids: []int64{uid}})
	if err != nil {
		return err
	}
	if c, ok := cm.RelationCount[uid]; ok && c.FollowCount >= 3000 {
		return biz.ErrFollowLimit
	}
	_, err = r.relationshipClient.Follow(ctx, &relationshipv1.FollowRequest{Uid: uid, ToUid: touid})
	return err
}

func (r *relationRepo) Unfollow(ctx context.Context, uid, touid int64) error {
	_, err := r.relationshipClient.Unfollow(ctx, &relationshipv1.FollowRequest{Uid: uid, ToUid: touid})
	return err
}

func (r *relationRepo) GetRelationCount(ctx context.Context, uids []int64) (map[int64]*biz.RelationCount, error) {
	cm, err := r.relationshipClient.GetRelationCount(ctx, &relationshipv1.BatchRequest{Uids: uids})
	if err != nil {
		return nil, err
	}
	m := make(map[int64]*biz.RelationCount, len(cm.RelationCount))
	for k, v := range cm.RelationCount {
		m[k] = &biz.RelationCount{
			UID:           k,
			FollowCount:   v.FollowCount,
			FollowerCount: v.FollowerCount,
		}
	}
	return m, nil
}

func (r *relationRepo) GetUsersFollow(ctx context.Context, uids []int64) (map[int64][]*biz.FollowItem, error) {
	fm, err := r.relationshipClient.GetUsersFollow(ctx, &relationshipv1.BatchRequest{Uids: uids})
	if err != nil {
		return nil, err
	}
	m := make(map[int64][]*biz.FollowItem, len(fm.Fm))
	for k, v := range fm.Fm {
		m[k] = make([]*biz.FollowItem, len(v.List))
		for i := range v.List {
			m[k][i] = &biz.FollowItem{
				ToUID:      v.List[i].ToUid,
				CreateTime: v.List[i].CreateTime,
			}
		}
	}
	return m, nil
}

func (r *relationRepo) GetUserFollow(ctx context.Context, uid, lastid int64) ([]*biz.FollowItem, error) {
	res, err := r.relationshipClient.GetFollowList(ctx, &relationshipv1.ListRequest{Uid: uid, LastId: lastid})
	if err != nil {
		return nil, err
	}
	list := make([]*biz.FollowItem, len(res.List.List))
	for k, v := range res.List.List {
		list[k] = &biz.FollowItem{
			ToUID:      v.ToUid,
			CreateTime: v.CreateTime,
		}
	}
	return list, nil
}

func (r *relationRepo) GetUserFollower(ctx context.Context, uid int64, lastid int64) ([]*biz.FollowItem, error) {
	res, err := r.relationshipClient.GetFollowerList(ctx, &relationshipv1.ListRequest{Uid: uid, LastId: lastid})
	if err != nil {
		return nil, err
	}
	list := make([]*biz.FollowItem, len(res.List.List))
	for k, v := range res.List.List {
		list[k] = &biz.FollowItem{
			ToUID:      v.ToUid,
			CreateTime: v.CreateTime,
		}
	}
	return list, nil
}

func (r *relationRepo) GetRelation(ctx context.Context, uid int64, uids []int64) (map[int64]*biz.UserRelation, error) {
	rm, err := r.relationshipClient.GetRelation(ctx, &relationshipv1.RelationRequest{Uid: uid, Uids: uids})
	if err != nil {
		return nil, err
	}
	m := make(map[int64]*biz.UserRelation, len(rm.Rm))
	for k, v := range rm.Rm {
		m[k] = &biz.UserRelation{
			Relation:     v.Relation,
			FollowTime:   v.FollowTime,
			FollowedTime: v.FollowedTime,
		}
	}
	return m, nil
}
