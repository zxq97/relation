package api

import (
	"context"

	"google.golang.org/grpc"
)

type RelationClientImpl struct {
	client RelationClient
}

func NewRelationClientImpl(conn *grpc.ClientConn) *RelationClientImpl {
	return &RelationClientImpl{
		client: NewRelationClient(conn),
	}
}

func (rli *RelationClientImpl) Follow(ctx context.Context, uid, touid int64, source Source) error {
	_, err := rli.client.Follow(ctx, &FollowRequest{Uid: uid, ToUid: touid, Source: source})
	return err
}

func (rli *RelationClientImpl) Unfollow(ctx context.Context, uid, touid int64, source Source) error {
	_, err := rli.client.Unfollow(ctx, &FollowRequest{Uid: uid, ToUid: touid, Source: source})
	return err
}

func (rli *RelationClientImpl) GetFollowList(ctx context.Context, uid, lastid int64, source Source) ([]*ListItem, error) {
	res, err := rli.client.GetFollowList(ctx, &ListRequest{Uid: uid, LastId: lastid, Source: source})
	return listDTO2VO(res), err
}

func (rli *RelationClientImpl) GetFollowerList(ctx context.Context, uid, lastid int64, source Source) ([]*ListItem, error) {
	res, err := rli.client.GetFollowerList(ctx, &ListRequest{Uid: uid, LastId: lastid, Source: source})
	return listDTO2VO(res), err
}

func (rli *RelationClientImpl) GetRelation(ctx context.Context, uid int64, uids []int64, source Source) (map[int64]*RItem, error) {
	res, err := rli.client.GetRelation(ctx, &RelationRequest{Uid: uid, Uids: uids, Source: source})
	return itemDTO2VO(res), err
}

func (rli *RelationClientImpl) GetRelationCount(ctx context.Context, uids []int64, source Source) (map[int64]*CountItem, error) {
	res, err := rli.client.GetRelationCount(ctx, &CountRequest{Uids: uids, Source: source})
	return countDTO2VO(res), err
}
