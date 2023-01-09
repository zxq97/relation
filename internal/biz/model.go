package biz

import "github.com/zxq97/relation/internal/data"

type FollowItem struct {
	ToUID      int64 `json:"to_uid"`
	CreateTime int64 `json:"create_time"`
}

type UserFollowCount struct {
	UID           int64 `json:"uid"`
	FollowCount   int32 `json:"follow_count"`
	FollowerCount int32 `json:"follower_count"`
}

type UserRelation struct {
	Relation     int32 `json:"relation"`
	FollowTime   int64 `json:"follow_time"`
	FollowedTime int64 `json:"followed_time"`
}

func itemPO2DO(item *data.FollowItem) *FollowItem {
	return &FollowItem{
		ToUID:      item.ToUid,
		CreateTime: item.CreateTime,
	}
}

func listPO2DO(list []*data.FollowItem) []*FollowItem {
	l := make([]*FollowItem, len(list))
	for k, v := range list {
		l[k] = itemPO2DO(v)
	}
	return l
}

func imPO2DO(im map[int64][]*data.FollowItem) map[int64][]*FollowItem {
	m := make(map[int64][]*FollowItem, len(im))
	for k, v := range im {
		m[k] = listPO2DO(v)
	}
	return m
}

func countPO2DO(c *data.UserFollowCount) *UserFollowCount {
	return &UserFollowCount{
		UID:           c.UID,
		FollowCount:   c.FollowCount,
		FollowerCount: c.FollowerCount,
	}
}

func cmPO2DO(cm map[int64]*data.UserFollowCount) map[int64]*UserFollowCount {
	m := make(map[int64]*UserFollowCount, len(cm))
	for k, v := range cm {
		m[k] = countPO2DO(v)
	}
	return m
}
