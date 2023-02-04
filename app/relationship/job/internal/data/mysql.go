package data

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zxq97/relation/app/relationship/job/internal/biz"
	"gorm.io/gorm"
	"upper.io/db.v3"
)

func (r *relationshipRepo) follow(ctx context.Context, uid, touid int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		sql := "INSERT INTO user_follows (`uid`, `to_uid`) VALUES (?, ?)"
		if err := tx.Exec(sql, uid, touid).Error; err != nil {
			return errors.Wrap(err, "add follow")
		}
		sql = "INSERT INTO user_followers (`uid`, `to_uid`) VALUES (?, ?)"
		if err := tx.Exec(sql, touid, uid).Error; err != nil {
			return errors.Wrap(err, "add follower")
		}
		sql = "INSERT INTO user_relation_counts (`uid`, `follow_count`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `follow_count` = `follow_count` + 1"
		if err := tx.Exec(sql, uid, touid).Error; err != nil {
			return errors.Wrap(err, "add follow count")
		}
		sql = "INSERT INTO user_relation_counts (`uid`, `follower_count`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `follower_count` = `follower_count` + 1"
		return errors.Wrap(tx.Exec(sql, touid, uid).Error, "add follower count")
	})
}

func (r *relationshipRepo) unfollow(ctx context.Context, uid, touid int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		sql := "DELETE FROM user_follows WHERE `uid` = ? AND `to_uid` = ? LIMIT 1"
		if err := tx.Exec(sql, uid, touid).Error; err != nil {
			return errors.Wrap(err, "delete follow")
		}
		sql = "DELETE FROM user_followers WHERE `uid` = ? AND `to_uid` = ? LIMIT 1"
		if err := tx.Exec(sql, touid, uid).Error; err != nil {
			return errors.Wrap(err, "delete follower")
		}
		sql = "UPDATE user_relation_counts SET `follow_count` = `follow_count` - 1 WHERE `uid` = ? AND `follow_count` > 0 LIMIT 1"
		if err := tx.Exec(sql, uid).Error; err != nil {
			return errors.Wrap(err, "reduce follow count")
		}
		sql = "UPDATE user_relation_counts SET `follower_count` = `follower_count` - 1 WHERE `uid` = ? AND `follower_count` > 0 LIMIT 1"
		return errors.Wrap(tx.Exec(sql, touid).Error, "reduce follower count")
	})
}

func (r *relationshipRepo) getUserFollower(ctx context.Context, uid, lastid int64) ([]*biz.FollowItem, error) {
	uf := []*UserFollower{}
	filter := db.Cond{"uid": uid}
	if lastid != 0 {
		sql := "SELECT `id` FROM user_followers WHERE `uid` = ? AND `to_uid` = ? LIMIT 1"
		var id int64
		err := r.db.WithContext(ctx).Raw(sql, uid, lastid).Scan(&id).Error
		if err != nil {
			return nil, errors.Wrap(err, "get last id")
		}
		filter["id < "] = id
	}
	err := r.db.WithContext(ctx).Select("to_uid", "create_at").Where(filter).Order("create_at DESC").Limit(100).Find(&uf).Error
	if err != nil {
		return nil, errors.Wrap(err, "db get user follower")
	}
	list := make([]*biz.FollowItem, len(uf))
	for k, v := range uf {
		list[k] = &biz.FollowItem{
			ToUid:      v.ToUID,
			CreateTime: v.CreateAt.UnixMilli(),
		}
	}
	return list, nil
}
