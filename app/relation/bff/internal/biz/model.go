package biz

type RelationCount struct {
	UID           int64
	FollowCount   int32
	FollowerCount int32
}

type FollowItem struct {
	ToUID      int64
	CreateTime int64
}

type UserRelation struct {
	Relation     int32
	FollowTime   int64
	FollowedTime int64
}
