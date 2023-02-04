package data

import "time"

type UserFollower struct {
	ID       int64
	ToUID    int64
	CreateAt time.Time
}
