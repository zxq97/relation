package data

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

const (
	tableUserFollow      = "user_follow"
	tableUserFollower    = "user_follower"
	tableUserFollowCount = "user_follow_count"
)

func follow(ctx context.Context, sess sqlbuilder.Database, uid, touid int64) error {
	return sess.Tx(ctx, func(sess sqlbuilder.Tx) error {
		sql := "INSERT INTO %s (`uid`, `to_uid`) VALUES (?, ?)"
		_, err := sess.Exec(fmt.Sprintf(sql, tableUserFollow), uid, touid)
		if err != nil {
			return errors.Wrap(err, "add follow")
		}
		_, err = sess.Exec(fmt.Sprintf(sql, tableUserFollower), touid, uid)
		if err != nil {
			return errors.Wrap(err, "add follower")
		}
		sql = "INSERT INTO %s (`uid`, `follow_count`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `follow_count` = `follow_count` + 1"
		_, err = sess.Exec(fmt.Sprintf(sql, tableUserFollowCount), uid, 1)
		if err != nil {
			return errors.Wrap(err, "incr follow count")
		}
		sql = "INSERT INTO %s (`uid`, `follower_count`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `follower_count` = `follower_count` + 1"
		_, err = sess.Exec(fmt.Sprintf(sql, tableUserFollowCount), touid, 1)
		return errors.Wrap(err, "incr follower count")
	})
}

func unfollow(ctx context.Context, sess sqlbuilder.Database, uid, touid int64) error {
	return sess.Tx(ctx, func(sess sqlbuilder.Tx) error {
		filter := db.Cond{"uid": uid, "to_uid": touid}
		_, err := sess.DeleteFrom(tableUserFollow).Where(filter).Limit(1).Exec()
		if err != nil {
			return errors.Wrap(err, "delete follow")
		}
		filter = db.Cond{"uid": touid, "to_uid": uid}
		_, err = sess.DeleteFrom(tableUserFollower).Where(filter).Limit(1).Exec()
		if err != nil {
			return errors.Wrap(err, "delete follower")
		}
		filter = db.Cond{"uid": uid}
		_, err = sess.Update(tableUserFollowCount).Set("`follow_count` = `follow_count` - 1 AND `follow_count` > 0").Where(filter).Exec()
		if err != nil {
			return errors.Wrap(err, "decr follow count")
		}
		filter = db.Cond{"uid": touid}
		_, err = sess.Update(tableUserFollowCount).Set("`follower_count` = `follower_count` - 1 AND `follower_count` > 0").Where(filter).Exec()
		return errors.Wrap(err, "decr follower count")
	})
}

func getUserFollow(ctx context.Context, sess sqlbuilder.Database, uid int64) ([]*FollowItem, error) {
	items := []*UserFollow{}
	filter := db.Cond{"uid": uid}
	err := sess.WithContext(ctx).Select("`to_uid`", "`create_time`").From(tableUserFollow).Where(filter).OrderBy("`create_time` DESC").All(&items)
	if err != nil {
		return nil, errors.Wrap(err, "get user follow")
	}
	return po2pb(items), nil
}

func getUserFollower(ctx context.Context, sess sqlbuilder.Database, uid, lastid int64, limit int) ([]*FollowItem, error) {
	items := []*UserFollow{}
	filter := db.Cond{"uid": uid}
	if lastid != 0 {
		sql := fmt.Sprintf("SELECT `id` FROM %s WHERE `uid` = ? AND `to_uid` = ? LIMIT 1", tableUserFollower)
		row, err := sess.QueryRowContext(ctx, sql, uid, lastid)
		if err != nil {
			return nil, errors.Wrap(err, "get last id")
		}
		var id int64
		err = row.Scan(&id)
		if err != nil {
			return nil, errors.Wrap(err, "row scan")
		}
		filter["id < "] = id
	}
	err := sess.WithContext(ctx).Select("`to_uid`", "`create_time`").From(tableUserFollower).Where(filter).OrderBy("`create_time` DESC").Limit(limit).All(items)
	if err != nil {
		return nil, errors.Wrap(err, "get user follower")
	}
	return po2pb(items), nil
}

func getFollowCount(ctx context.Context, sess sqlbuilder.Database, uids []int64) (map[int64]*UserFollowCount, error) {
	counts := []*UserFollowCount{}
	filter := db.Cond{"uid IN": uids}
	err := sess.WithContext(ctx).SelectFrom(tableUserFollowCount).Where(filter).All(&counts)
	if err != nil {
		return nil, err
	}
	cntMap := make(map[int64]*UserFollowCount, len(uids))
	for _, v := range counts {
		cntMap[v.UID] = v
	}
	return cntMap, nil
}
