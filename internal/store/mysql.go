package store

import (
	"context"
	"fmt"

	"github.com/zxq97/relation/internal/model"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

const (
	tableUserFollow      = "user_follow"
	tableUserFollower    = "user_follower"
	tableUserFollowCount = "user_follow_count"
)

func Follow(ctx context.Context, uid, touid int64) error {
	return dbCli.Tx(ctx, func(sess sqlbuilder.Tx) error {
		sql := "INSERT INTO %s (`uid`, `to_uid`) VALUES (?, ?)"
		_, err := sess.Exec(fmt.Sprintf(sql, tableUserFollow), uid, touid)
		if err != nil {
			return err
		}
		_, err = sess.Exec(fmt.Sprintf(sql, tableUserFollower), touid, uid)
		if err != nil {
			return err
		}
		sql = "INSERT INTO %s (`uid`, `follow_count`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `follow_count` = `follow_count` + 1"
		_, err = sess.Exec(fmt.Sprintf(sql, tableUserFollowCount), uid, 1)
		if err != nil {
			return err
		}
		sql = "INSERT INTO %s (`uid`, `follower_count`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `follower_count` = `follower_count` + 1"
		_, err = sess.Exec(fmt.Sprintf(sql, tableUserFollowCount), touid, 1)
		return err
	})
}

func Unfollow(ctx context.Context, uid, touid int64) error {
	return dbCli.Tx(ctx, func(sess sqlbuilder.Tx) error {
		filter := db.Cond{"uid": uid, "to_uid": touid}
		_, err := sess.DeleteFrom(tableUserFollow).Where(filter).Limit(1).Exec()
		if err != nil {
			return err
		}
		filter = db.Cond{"uid": touid, "to_uid": uid}
		_, err = sess.DeleteFrom(tableUserFollower).Where(filter).Limit(1).Exec()
		if err != nil {
			return err
		}
		filter = db.Cond{"uid": uid}
		_, err = sess.Update(tableUserFollowCount).Set("`follow_count` = `follow_count` - 1 AND `follow_count` > 0").Where(filter).Exec()
		if err != nil {
			return err
		}
		filter = db.Cond{"uid": touid}
		_, err = sess.Update(tableUserFollowCount).Set("`follower_count` = `follower_count` - 1 AND `follower_count` > 0").Where(filter).Exec()
		return err
	})
}

func GetAllUserFollow(ctx context.Context, uid int64) ([]*model.FollowItem, error) {
	items := []*model.FollowItem{}
	filter := db.Cond{"uid": uid}
	err := dbCli.WithContext(ctx).Select("`to_uid`", "`create_time`").From(tableUserFollow).Where(filter).OrderBy("`create_time` DESC").All(&items)
	return items, err
}

func GetUserFollower(ctx context.Context, uid, lastid int64, limit int) ([]*model.FollowItem, error) {
	items := []*model.FollowItem{}
	filter := db.Cond{"uid": uid}
	if lastid != 0 {
		sql := fmt.Sprintf("SELECT `id` FROM %s WHERE `uid` = ? AND `to_uid` = ? LIMIT 1", tableUserFollower)
		row, err := dbCli.QueryRowContext(ctx, sql, uid, lastid)
		if err != nil {
			return items, err
		}
		var id int64
		err = row.Scan(&id)
		if err != nil {
			return items, err
		}
		filter["id < "] = id
	}
	err := dbCli.WithContext(ctx).Select("`to_uid`", "`create_time`").From(tableUserFollower).Where(filter).OrderBy("`create_time` DESC").Limit(limit).All(items)
	return items, err
}

func GetUsersFollowCount(ctx context.Context, uids []int64) (map[int64]*model.UserFollowCount, error) {
	counts := []*model.UserFollowCount{}
	filter := db.Cond{"uid IN": uids}
	err := dbCli.WithContext(ctx).SelectFrom(tableUserFollowCount).Where(filter).All(&counts)
	if err != nil {
		return nil, err
	}
	cntMap := make(map[int64]*model.UserFollowCount, len(uids))
	for _, v := range counts {
		cntMap[v.UID] = v
	}
	return cntMap, nil
}
