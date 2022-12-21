package relationbff

import (
	"context"

	"github.com/zxq97/relation/internal/constant"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/service/relationsvc"
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
		return &relation.EmptyResponse{}, relation.ErrSourceUndefined
	}
	// todo check black

	res, err := client.GetRelationCount(ctx, &relationsvc.BatchRequest{Uids: []int64{req.Uid}})
	if err != nil {
		return &relation.EmptyResponse{}, err
	}
	rc, ok := res.RelationCount[req.Uid]
	if !ok || rc.FollowCount < constant.FollowLimit {
		_, err = client.Follow(ctx, &relationsvc.FollowRequest{Uid: req.Uid, ToUid: req.ToUid})
		return &relation.EmptyResponse{}, err
	} else {
		return &relation.EmptyResponse{}, relation.ErrFollowLimit
	}
}

func (RelationBFF) Unfollow(ctx context.Context, req *relation.FollowRequest) (*relation.EmptyResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.EmptyResponse{}, relation.ErrSourceUndefined
	}
	_, err := client.Unfollow(ctx, &relationsvc.FollowRequest{Uid: req.Uid, ToUid: req.ToUid})
	return &relation.EmptyResponse{}, err
}

func (RelationBFF) GetFollowList(ctx context.Context, req *relation.ListRequest) (*relation.ListResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.ListResponse{}, relation.ErrSourceUndefined
	}
	res, err := client.GetFollowList(ctx, &relationsvc.ListRequest{Uid: req.Uid, LastId: req.LastId})
	return translateList(res), err
}

func (RelationBFF) GetFollowerList(ctx context.Context, req *relation.ListRequest) (*relation.ListResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.ListResponse{}, relation.ErrSourceUndefined
	}
	res, err := client.GetFollowerList(ctx, &relationsvc.ListRequest{Uid: req.Uid, LastId: req.LastId})
	return translateList(res), err
}

func (RelationBFF) GetRelation(ctx context.Context, req *relation.RelationRequest) (*relation.RelationResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.RelationResponse{}, relation.ErrSourceUndefined
	}
	res, err := client.GetRelation(ctx, &relationsvc.RelationRequest{Uid: req.Uid, Uids: req.Uids})
	return translateRelation(res), err
}

func (RelationBFF) GetRelationCount(ctx context.Context, req *relation.CountRequest) (*relation.CountResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.CountResponse{}, relation.ErrSourceUndefined
	}
	res, err := client.GetRelationCount(ctx, &relationsvc.BatchRequest{Uids: req.Uids})
	return translateCount(res), err
}

//GetCommonRelation 共同关注
func (RelationBFF) GetCommonRelation(ctx context.Context, req *relation.FollowRequest) (*relation.BatchResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.BatchResponse{}, relation.ErrSourceUndefined
	}
	fm, err := client.GetUsersFollow(ctx, &relationsvc.BatchRequest{Uids: []int64{req.Uid, req.ToUid}})
	if err != nil {
		return &relation.BatchResponse{}, err
	}
	uf, ok := fm.Fm[req.Uid]
	if !ok || uf == nil {
		return &relation.BatchResponse{}, nil
	}
	tf, ok := fm.Fm[req.ToUid]
	if !ok || tf == nil {
		return &relation.BatchResponse{}, nil
	}
	m := make(map[int64]struct{}, len(uf.List))
	for _, v := range uf.List {
		m[v.ToUid] = struct{}{}
	}
	uids := make([]int64, 0, len(tf.List))
	for _, v := range tf.List {
		if _, ok = m[v.ToUid]; ok {
			uids = append(uids, v.ToUid)
		}
	}
	return &relation.BatchResponse{Uids: uids}, nil
}

//GetRelationChain 我关注的人关注了他
func (RelationBFF) GetRelationChain(ctx context.Context, req *relation.FollowRequest) (*relation.BatchResponse, error) {
	if req.Source == relation.Source_Undefined {
		return &relation.BatchResponse{}, relation.ErrSourceUndefined
	}
	// todo need logic
	fm, err := client.GetUsersFollow(ctx, &relationsvc.BatchRequest{Uids: []int64{req.Uid, req.ToUid}})
	if err != nil {
		return &relation.BatchResponse{}, err
	}
	uf, ok := fm.Fm[req.Uid]
	if !ok || uf == nil {
		return &relation.BatchResponse{}, nil
	}
	return &relation.BatchResponse{}, nil
}
