package data

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zxq97/relation/app/relationship/pkg/bizdata"
	"github.com/zxq97/relation/app/relationship/pkg/dal/query"
	"github.com/zxq97/relation/app/relationship/pkg/message"
)

func (r *relationshipRepo) follow(ctx context.Context, uid, touid int64) error {
	return r.q.Transaction(func(tx *query.Query) error {
		if err := tx.WithContext(ctx).UserFollow.InsertFollow(uid, touid); err != nil {
			return errors.Wrap(err, "add follow")
		}
		if err := tx.WithContext(ctx).UserFollower.InsertFollower(touid, uid); err != nil {
			return errors.Wrap(err, "add follower")
		}
		if err := tx.WithContext(ctx).UserRelationCount.IncrFollowCount(uid); err != nil {
			return errors.Wrap(err, "incr follow count")
		}
		return errors.Wrap(tx.WithContext(ctx).ExtraFollower.InsertFollower(touid), "add extra follower")
	})
}

func (r *relationshipRepo) unfollow(ctx context.Context, uid, touid int64) error {
	return r.q.Transaction(func(tx *query.Query) error {
		if err := tx.WithContext(ctx).UserFollow.DeleteFollow(uid, touid); err != nil {
			return errors.Wrap(err, "delete follow")
		}
		if err := tx.WithContext(ctx).UserFollower.DeleteFollower(touid, uid); err != nil {
			return errors.Wrap(err, "delete follower")
		}
		if err := tx.WithContext(ctx).UserRelationCount.DecrFollowCount(uid); err != nil {
			return errors.Wrap(err, "decr follow count")
		}
		return errors.Wrap(tx.WithContext(ctx).ExtraFollower.DeleteFollower(touid), "delete extra follower")
	})
}

func (r *relationshipRepo) getUserFollower(ctx context.Context, uid, lastid int64) ([]*bizdata.FollowItem, error) {
	rows, err := r.q.WithContext(ctx).UserFollower.FindUserFollowerByLastID(uid, lastid, message.RebuildBatchSize)
	if err != nil {
		return nil, errors.Wrap(err, "get user follower")
	}
	list := make([]*bizdata.FollowItem, len(rows))
	for k, v := range rows {
		list[k] = &bizdata.FollowItem{ToUid: v.ToUID, CreateTime: v.CreateAt.UnixMilli()}
	}
	return list, nil
}
