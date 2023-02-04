package service

import (
	"context"

	"github.com/zxq97/relation/api/relation/bff/v1"
	"github.com/zxq97/relation/app/relation/bff/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RelationService struct {
	v1.UnimplementedRelationBFFServer
	relationUC *biz.RelationUseCase
	blackUC    *biz.BlackUseCase
}

func NewRelationService(relationUC *biz.RelationUseCase, blackUC *biz.BlackUseCase) *RelationService {
	return &RelationService{relationUC: relationUC, blackUC: blackUC}
}

func (s *RelationService) Follow(ctx context.Context, req *v1.FollowRequest) (*emptypb.Empty, error) {
	err := s.relationUC.Follow(ctx, req.Uid, req.ToUid)
	return &emptypb.Empty{}, err
}

func (s *RelationService) Unfollow(ctx context.Context, req *v1.FollowRequest) (*emptypb.Empty, error) {
	err := s.relationUC.Unfollow(ctx, req.Uid, req.ToUid)
	return &emptypb.Empty{}, err
}

func (s *RelationService) GetFollowList(ctx context.Context, req *v1.ListRequest) (*v1.ListResponse, error) {
	list, err := s.relationUC.GetFollowList(ctx, req.Uid, req.LastId)
	if err != nil {
		return &v1.ListResponse{}, err
	}
	return listDO2DTO(list), nil
}

func (s *RelationService) GetFollowerList(ctx context.Context, req *v1.ListRequest) (*v1.ListResponse, error) {
	list, err := s.relationUC.GetFollowerList(ctx, req.Uid, req.LastId)
	if err != nil {
		return &v1.ListResponse{}, err
	}
	return listDO2DTO(list), nil
}

func (s *RelationService) GetRelation(ctx context.Context, req *v1.RelationRequest) (*v1.RelationResponse, error) {
	rm, err := s.relationUC.GetRelation(ctx, req.Uid, req.Uids)
	if err != nil {
		return &v1.RelationResponse{}, err
	}
	return &v1.RelationResponse{Rm: rmDO2DTO(rm)}, nil
}

func (s *RelationService) GetRelationCount(ctx context.Context, req *v1.CountRequest) (*v1.CountResponse, error) {
	cm, err := s.relationUC.GetRelationCount(ctx, req.Uids)
	if err != nil {
		return &v1.CountResponse{}, err
	}
	return &v1.CountResponse{RelationCount: cmDO2DTO(cm)}, nil
}

func (s *RelationService) GetCommonRelation(ctx context.Context, req *v1.FollowRequest) (*v1.BatchResponse, error) {
	uids, err := s.relationUC.GetCommonRelation(ctx, req.Uid, req.ToUid)
	if err != nil {
		return &v1.BatchResponse{}, err
	}
	return &v1.BatchResponse{Uids: uids}, nil
}

func (s *RelationService) GetRelationChain(ctx context.Context, req *v1.FollowRequest) (*v1.BatchResponse, error) {
	uids, err := s.relationUC.GetRelationChain(ctx, req.Uid, req.ToUid)
	if err != nil {
		return &v1.BatchResponse{}, err
	}
	return &v1.BatchResponse{Uids: uids}, nil
}
