package data

import "time"

type UserFollow struct {
	ID       int64
	ToUID    int64
	CreateAt time.Time
}

type UserFollower struct {
	ID       int64
	ToUID    int64
	CreateAt time.Time
}

type UserRelationCount struct {
	UID           int64
	FollowCount   int32
	FollowerCount int32
}
