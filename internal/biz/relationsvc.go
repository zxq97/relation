package biz

import (
	"context"

	"github.com/zxq97/relation/internal/data"
)

type RelationSVCRepo interface {
	Follow(context.Context, int64, int64) error
	Unfollow(context.Context, int64, int64) error
	GetFollowList(context.Context, int64, int64) ([]*data.FollowItem, error)
	GetFollowerList(context.Context, int64, int64) ([]*data.FollowItem, error)
	GetRelation(context.Context, int64, []int64) (map[int64]*data.UserRelation, error)
	GetRelationCount(context.Context, []int64) (map[int64]*data.UserFollowCount, error)
	GetUsersFollow(context.Context, []int64) (map[int64][]*data.UserFollow, error)
}
