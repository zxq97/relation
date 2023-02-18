package data

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zxq97/relation/app/relationship/pkg/bizdata"
	"github.com/zxq97/relation/app/relationship/pkg/dal/model"
	"github.com/zxq97/relation/app/relationship/pkg/message"
)

func (r *relationshipRepo) getRelationCount(ctx context.Context, uids []int64) (map[int64]*bizdata.RelationCount, error) {
	rows, err := r.q.WithContext(ctx).UserRelationCount.FindUsersRelationCount(uids)
	if err != nil {
		return nil, errors.Wrap(err, "db get relation count")
	}
	m := make(map[int64]*bizdata.RelationCount, len(uids))
	for _, v := range rows {
		m[v.UID] = &bizdata.RelationCount{
			Uid:           v.UID,
			FollowCount:   v.FollowCount,
			FollowerCount: v.FollowerCount,
		}
	}
	for _, v := range uids {
		if _, ok := m[v]; !ok {
			m[v] = &bizdata.RelationCount{
				Uid: v,
			}
		}
	}
	return m, nil
}

func (r *relationshipRepo) getUserFollow(ctx context.Context, uid int64) ([]*bizdata.FollowItem, error) {
	rows, err := r.q.WithContext(ctx).UserFollow.FindUserFollow(uid)
	if err != nil {
		return nil, errors.Wrap(err, "db get user follow")
	}
	list := make([]*bizdata.FollowItem, len(rows))
	for k, v := range rows {
		list[k] = &bizdata.FollowItem{
			ToUid:      v.ToUID,
			CreateTime: v.CreateAt.UnixMilli(),
		}
	}
	return list, nil
}

func (r *relationshipRepo) getUserFollower(ctx context.Context, uid, lastid int64) ([]*bizdata.FollowItem, error) {
	var (
		rows []*model.UserFollower
		err  error
	)
	if lastid == 0 {
		rows, err = r.q.WithContext(ctx).UserFollower.FindUserFollower(uid, message.ListBatchSize)
	} else {
		rows, err = r.q.WithContext(ctx).UserFollower.FindUserFollowerByLastID(uid, lastid, message.ListBatchSize)
	}
	if err != nil {
		return nil, errors.Wrap(err, "db get follower")
	}
	list := make([]*bizdata.FollowItem, len(rows))
	for k, v := range rows {
		list[k] = &bizdata.FollowItem{
			ToUid:      v.ToUID,
			CreateTime: v.CreateAt.UnixMilli(),
		}
	}
	return list, nil
}
