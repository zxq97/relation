package relationbff

import (
	"context"

	account "github.com/zxq97/account/api"
	"github.com/zxq97/relation/api"
	"github.com/zxq97/relation/internal/env"
	"github.com/zxq97/relation/internal/service/relationsvc"
	"google.golang.org/grpc"
)

type RelationBFF struct {
	client        relationsvc.RelationSvcClient
	accountClient *account.AccountClientImpl
}

func InitRelationBFF(conf *RelationBffConfig, conn, blackConn *grpc.ClientConn) (*RelationBFF, error) {
	err := env.InitLog(conf.LogPath)
	if err != nil {
		return nil, err
	}
	client := relationsvc.NewRelationSvcClient(conn)
	return &RelationBFF{
		client:        client,
		accountClient: account.NewAccountClientImpl(blackConn),
	}, nil
}

func (bff *RelationBFF) Follow(ctx context.Context, req *api.FollowRequest) (*api.EmptyResponse, error) {
	if req.Source == api.Source_Undefined {
		return &api.EmptyResponse{}, api.ErrSourceUndefined
	}
	// todo check black
	m, err := bff.accountClient.CheckBlacked(ctx, req.Uid, []int64{req.ToUid})
	if err != nil {
		return &api.EmptyResponse{}, err
	}
	if _, ok := m[req.ToUid]; ok {
		return &api.EmptyResponse{}, api.ErrBlacked
	}
	res, err := bff.client.GetRelationCount(ctx, &relationsvc.BatchRequest{Uids: []int64{req.Uid}})
	if err != nil {
		return &api.EmptyResponse{}, err
	}
	rc, ok := res.RelationCount[req.Uid]
	if !ok || rc.FollowCount < env.FollowCountLimit {
		_, err = bff.client.Follow(ctx, &relationsvc.FollowRequest{Uid: req.Uid, ToUid: req.ToUid})
		return &api.EmptyResponse{}, err
	} else {
		return &api.EmptyResponse{}, api.ErrFollowLimit
	}
}

func (bff *RelationBFF) Unfollow(ctx context.Context, req *api.FollowRequest) (*api.EmptyResponse, error) {
	if req.Source == api.Source_Undefined {
		return &api.EmptyResponse{}, api.ErrSourceUndefined
	}
	_, err := bff.client.Unfollow(ctx, &relationsvc.FollowRequest{Uid: req.Uid, ToUid: req.ToUid})
	return &api.EmptyResponse{}, err
}

func (bff *RelationBFF) GetFollowList(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	if req.Source == api.Source_Undefined {
		return &api.ListResponse{}, api.ErrSourceUndefined
	}
	res, err := bff.client.GetFollowList(ctx, &relationsvc.ListRequest{Uid: req.Uid, LastId: req.LastId})
	return translateList(res), err
}

func (bff *RelationBFF) GetFollowerList(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	if req.Source == api.Source_Undefined {
		return &api.ListResponse{}, api.ErrSourceUndefined
	}
	res, err := bff.client.GetFollowerList(ctx, &relationsvc.ListRequest{Uid: req.Uid, LastId: req.LastId})
	return translateList(res), err
}

func (bff *RelationBFF) GetRelation(ctx context.Context, req *api.RelationRequest) (*api.RelationResponse, error) {
	if req.Source == api.Source_Undefined {
		return &api.RelationResponse{}, api.ErrSourceUndefined
	}
	res, err := bff.client.GetRelation(ctx, &relationsvc.RelationRequest{Uid: req.Uid, Uids: req.Uids})
	return translateRelation(res), err
}

func (bff *RelationBFF) GetRelationCount(ctx context.Context, req *api.CountRequest) (*api.CountResponse, error) {
	if req.Source == api.Source_Undefined {
		return &api.CountResponse{}, api.ErrSourceUndefined
	}
	res, err := bff.client.GetRelationCount(ctx, &relationsvc.BatchRequest{Uids: req.Uids})
	return translateCount(res), err
}

//GetCommonRelation 共同关注
func (bff *RelationBFF) GetCommonRelation(ctx context.Context, req *api.FollowRequest) (*api.BatchResponse, error) {
	if req.Source == api.Source_Undefined {
		return &api.BatchResponse{}, api.ErrSourceUndefined
	}
	fm, err := bff.client.GetUsersFollow(ctx, &relationsvc.BatchRequest{Uids: []int64{req.Uid, req.ToUid}})
	if err != nil {
		return &api.BatchResponse{}, err
	}
	uf, ok := fm.Fm[req.Uid]
	if !ok || uf == nil {
		return &api.BatchResponse{}, nil
	}
	tf, ok := fm.Fm[req.ToUid]
	if !ok || tf == nil {
		return &api.BatchResponse{}, nil
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
	return &api.BatchResponse{Uids: uids}, nil
}

//GetRelationChain 我关注的人关注了他
func (bff *RelationBFF) GetRelationChain(ctx context.Context, req *api.FollowRequest) (*api.BatchResponse, error) {
	if req.Source == api.Source_Undefined {
		return &api.BatchResponse{}, api.ErrSourceUndefined
	}
	fm, err := bff.client.GetUsersFollow(ctx, &relationsvc.BatchRequest{Uids: []int64{req.Uid, req.ToUid}})
	if err != nil {
		return &api.BatchResponse{}, err
	}
	uf, ok := fm.Fm[req.Uid]
	if !ok || uf == nil {
		return &api.BatchResponse{}, nil
	}
	uids := make([]int64, len(uf.List))
	for k, v := range uf.List {
		uids[k] = v.ToUid
	}
	rm, err := bff.client.GetRelation(ctx, &relationsvc.RelationRequest{Uid: req.ToUid, Uids: uids})
	if err != nil {
		return &api.BatchResponse{}, err
	}
	fs := make([]int64, 0, len(rm.Rm))
	for k, v := range rm.Rm {
		if v.FollowedTime > 0 && v.Relation&env.RelationFollowerBIT != 0 {
			fs = append(fs, k)
		}
	}
	return &api.BatchResponse{Uids: fs}, nil
}
