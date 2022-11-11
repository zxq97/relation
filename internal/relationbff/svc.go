package relationbff

import (
	"context"

	"github.com/zxq97/relation/internal/constant"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/relationsvc"
	"github.com/zxq97/relation/pkg/relation"
	"google.golang.org/grpc"
)

var (
	client relationsvc.RelationSvcClient
)

type RelationBFF struct {
}

func InitRelationBFF(conf *RelationBffConfig, conn *grpc.ClientConn) error {
	err := env.InitLog(conf.LogPath)
	if err != nil {
		return err
	}
	client = relationsvc.NewRelationSvcClient(conn)
	return nil
}

func (RelationBFF) Follow(ctx context.Context, req *relation.FollowRequest) (*relation.EmptyResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.EmptyResponse{}, ErrSourceUndefined
	}
	// todo check black

	res, err := client.GetRelationCount(ctx, &relationsvc.CountRequest{Uids: []int64{req.Uid}})
	if err != nil {
		return &relation.EmptyResponse{}, err
	}
	rc, ok := res.RelationCount[req.Uid]
	if !ok || rc.FollowCount < constant.FollowLimit {
		_, err = client.Follow(ctx, &relationsvc.FollowRequest{Uid: req.Uid, ToUid: req.ToUid})
		return &relation.EmptyResponse{}, err
	} else {
		return &relation.EmptyResponse{}, ErrFollowLimit
	}
}

func (RelationBFF) Unfollow(ctx context.Context, req *relation.FollowRequest) (*relation.EmptyResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.EmptyResponse{}, ErrSourceUndefined
	}
	_, err := client.Unfollow(ctx, &relationsvc.FollowRequest{Uid: req.Uid, ToUid: req.ToUid})
	return &relation.EmptyResponse{}, err
}

func (RelationBFF) GetFollowList(ctx context.Context, req *relation.ListRequest) (*relation.ListResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.ListResponse{}, ErrSourceUndefined
	}
	res, err := client.GetFollowList(ctx, &relationsvc.ListRequest{Uid: req.Uid, LastId: req.LastId})
	return translateList(res), err
}

func (RelationBFF) GetFollowerList(ctx context.Context, req *relation.ListRequest) (*relation.ListResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.ListResponse{}, ErrSourceUndefined
	}
	res, err := client.GetFollowerList(ctx, &relationsvc.ListRequest{Uid: req.Uid, LastId: req.LastId})
	return translateList(res), err
}

func (RelationBFF) GetRelation(ctx context.Context, req *relation.RelationRequest) (*relation.RelationResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.RelationResponse{}, ErrSourceUndefined
	}
	res, err := client.GetRelation(ctx, &relationsvc.RelationRequest{Uid: req.Uid, Uids: req.Uids})
	return translateRelation(res), err
}

func (RelationBFF) GetRelationCount(ctx context.Context, req *relation.CountRequest) (*relation.CountResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.CountResponse{}, ErrSourceUndefined
	}
	res, err := client.GetRelationCount(ctx, &relationsvc.CountRequest{Uids: req.Uids})
	return translateCount(res), err
}
