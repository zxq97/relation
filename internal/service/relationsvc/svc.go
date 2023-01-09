package relationsvc

import (
	"context"

	"github.com/zxq97/relation/internal/biz"
)

type RelationSvc struct {
	biz *biz.RelationSvcBIZ
}

func InitRelationSvc(conf *RelationSvcConfig) (*RelationSvc, error) {
	rsb, err := biz.NewRelationSvcBIZ(conf.Redis["redis"], conf.MC["mc"], conf.Mysql["relation"], conf.Kafka["kafka"].Addr)
	if err != nil {
		return nil, err
	}
	return &RelationSvc{
		biz: rsb,
	}, nil
}

//Follow 关注
func (svc *RelationSvc) Follow(ctx context.Context, req *FollowRequest) (*EmptyResponse, error) {
	return &EmptyResponse{}, svc.biz.Follow(ctx, req.Uid, req.ToUid)
}

//Unfollow 取关
func (svc *RelationSvc) Unfollow(ctx context.Context, req *FollowRequest) (*EmptyResponse, error) {
	return &EmptyResponse{}, svc.biz.Unfollow(ctx, req.Uid, req.ToUid)
}

//GetFollowList 关注列表
func (svc *RelationSvc) GetFollowList(ctx context.Context, req *ListRequest) (*FollowList, error) {
	list, err := svc.biz.GetFollowList(ctx, req.Uid, req.LastId)
	if err != nil {
		return &FollowList{}, err
	}
	return listDO2DTO(list), nil
}

//GetFollowerList 粉丝列表
func (svc *RelationSvc) GetFollowerList(ctx context.Context, req *ListRequest) (*FollowList, error) {
	list, err := svc.biz.GetFollowerList(ctx, req.Uid, req.LastId)
	if err != nil {
		return &FollowList{}, err
	}
	return listDO2DTO(list), nil
}

//GetRelation 好有关系
func (svc *RelationSvc) GetRelation(ctx context.Context, req *RelationRequest) (*RelationResponse, error) {
	m, err := svc.biz.GetRelation(ctx, req.Uid, req.Uids)
	if err != nil {
		return &RelationResponse{}, err
	}
	return rmDO2DTO(m), nil
}

//GetRelationCount 关注 粉丝 数量
func (svc *RelationSvc) GetRelationCount(ctx context.Context, req *BatchRequest) (*CountResponse, error) {
	m, err := svc.biz.GetRelationCount(ctx, req.Uids)
	if err != nil {
		return &CountResponse{}, err
	}
	return cmDO2DTO(m), nil
}

//GetUsersFollow 获取全量关注
func (svc *RelationSvc) GetUsersFollow(ctx context.Context, req *BatchRequest) (*UserFollowResponse, error) {
	m, err := svc.biz.GetUsersFollow(ctx, req.Uids)
	if err != nil {
		return &UserFollowResponse{}, err
	}
	return lmDO2DTO(m), nil
}
