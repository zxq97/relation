package model

import "time"

type UserFollow struct {
	ToUID      int64     `db:"to_uid"`
	CreateTime time.Time `db:"create_time"`
}

type UserFollowCount struct {
	UID           int64 `db:"uid"`
	FollowCount   int32 `db:"uid"`
	FollowerCount int32 `db:"uid"`
}

func DAO2DTO(list []*UserFollow) []*FollowItem {
	itemList := make([]*FollowItem, len(list))
	for k, v := range list {
		itemList[k] = &FollowItem{
			ToUid:      v.ToUID,
			CreateTime: v.CreateTime.UnixMilli(),
		}
	}
	return itemList
}
