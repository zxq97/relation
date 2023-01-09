package api

import (
	"context"

	"google.golang.org/grpc"
)

var (
	client RelationClient
)

func InitClient(conn *grpc.ClientConn) {
	client = NewRelationClient(conn)
}

func Follow(ctx context.Context, uid, touid int64, source Source) error {
	_, err := client.Follow(ctx, &FollowRequest{Uid: uid, ToUid: touid, Source: source})
	return err
}

func Unfollow(ctx context.Context, uid, touid int64, source Source) error {
	_, err := client.Unfollow(ctx, &FollowRequest{Uid: uid, ToUid: touid, Source: source})
	return err
}

func GetFollowList(ctx context.Context, uid, lastid int64, source Source) ([]*ListItem, error) {
	res, err := client.GetFollowList(ctx, &ListRequest{Uid: uid, LastId: lastid, Source: source})
	return listDTO2DO(res), err
}

func GetFollowerList(ctx context.Context, uid, lastid int64, source Source) ([]*ListItem, error) {
	res, err := client.GetFollowerList(ctx, &ListRequest{Uid: uid, LastId: lastid, Source: source})
	return listDTO2DO(res), err
}

func GetRelation(ctx context.Context, uid int64, uids []int64, source Source) (map[int64]*RItem, error) {
	res, err := client.GetRelation(ctx, &RelationRequest{Uid: uid, Uids: uids, Source: source})
	return itemDTO2DO(res), err
}

func GetRelationCount(ctx context.Context, uids []int64, source Source) (map[int64]*CountItem, error) {
	res, err := client.GetRelationCount(ctx, &CountRequest{Uids: uids, Source: source})
	return countDTO2DO(res), err
}
