package model

import "time"

type UserFollow struct {
	UID        int64     `db:"uid"`
	ToUID      int64     `db:"to_uid"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
}

type UserFollowCount struct {
	UID           int64 `db:"uid"`
	FollowCount   int32 `db:"uid"`
	FollowerCount int32 `db:"uid"`
}
