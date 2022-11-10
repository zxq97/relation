package model

import (
	"time"

	"github.com/zxq97/relation/internal/relationsvc"
)

type UserFollow struct {
	ToUID      int64     `db:"to_uid"`
	CreateTime time.Time `db:"create_time"`
}

type UserFollowCount struct {
	UID           int64 `db:"uid"`
	FollowCount   int32 `db:"uid"`
	FollowerCount int32 `db:"uid"`
}

func UfDAO2DTO(list []*UserFollow) []*FollowItem {
	itemList := make([]*FollowItem, len(list))
	for k, v := range list {
		itemList[k] = &FollowItem{
			ToUid:      v.ToUID,
			CreateTime: v.CreateTime.UnixMilli(),
		}
	}
	return itemList
}

func FcDAO2DTO(m map[int64]*UserFollowCount) map[int64]*relationsvc.RelationCount {
	rm := make(map[int64]*relationsvc.RelationCount, len(m))
	for k, v := range m {
		rm[k] = &relationsvc.RelationCount{
			FollowCount:   int64(v.FollowCount),
			FollowerCount: int64(v.FollowerCount),
		}
	}
	return rm
}
