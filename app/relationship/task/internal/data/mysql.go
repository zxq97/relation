package data

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zxq97/relation/app/relationship/pkg/dal/query"
)

func (r *relationshipRepo) syncRecord(ctx context.Context, limit int64) (map[int64]int64, error) {
	m := make(map[int64]int64)
	err := r.q.Transaction(func(tx *query.Query) error {
		rows, err := tx.WithContext(ctx).ExtraFollower.FindUnSyncRecord(limit)
		if err != nil || len(rows) == 0 {
			return errors.Wrap(err, "get record")
		}
		ids := make([]int64, len(rows))
		for k, v := range rows {
			ids[k] = v.ID
			if v.Stats == 0 {
				m[v.UID]++
			} else {
				m[v.UID]--
			}
		}
		for k, v := range m {
			if v > 0 {
				if err = tx.WithContext(ctx).UserRelationCount.IncrByFollowerCount(k, v); err != nil {
					return errors.Wrap(err, "incr follower")
				}
			} else if v < 0 {
				if err = tx.WithContext(ctx).UserRelationCount.DecrByFollowerCount(k, -v); err != nil {
					return errors.Wrap(err, "decr follower")
				}
			}
		}
		return errors.Wrap(tx.WithContext(ctx).ExtraFollower.DeleteRecord(ids), "delete extra follower")
	})
	return m, err
}

func (r *relationshipRepo) syncRecordByUID(ctx context.Context, uid int64) (int64, error) {
	var cnt int64
	err := r.q.Transaction(func(tx *query.Query) error {
		rows, err := tx.WithContext(ctx).ExtraFollower.FindUnSyncRecordByUID(uid)
		if err != nil || len(rows) == 0 {
			return errors.Wrap(err, "get record")
		}
		ids := make([]int64, len(rows))
		for k, v := range rows {
			ids[k] = v.ID
			if v.Stats == 0 {
				cnt++
			} else {
				cnt--
			}
		}
		if cnt > 0 {
			if err = tx.WithContext(ctx).UserRelationCount.IncrByFollowerCount(uid, cnt); err != nil {
				return errors.Wrap(err, "incr follower")
			}
		} else if cnt < 0 {
			if err = tx.WithContext(ctx).UserRelationCount.DecrByFollowerCount(uid, -cnt); err != nil {
				return errors.Wrap(err, "decr follower")
			}
		}
		return errors.Wrap(tx.WithContext(ctx).ExtraFollower.DeleteRecord(ids), "delete extra follower")
	})
	return cnt, err
}
