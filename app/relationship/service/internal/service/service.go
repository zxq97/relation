package service

import (
	"context"

	"github.com/zxq97/relation/api/relationship/service/v1"
	"github.com/zxq97/relation/app/relationship/service/internal/biz"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RelationshipService struct {
	v1.UnimplementedRelationSvcServer
	uc *biz.RelationshipUseCase
}

func NewRelationshipService(uc *biz.RelationshipUseCase) *RelationshipService {
	return &RelationshipService{uc: uc}
}

func (s *RelationshipService) Follow(ctx context.Context, req *v1.FollowRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.uc.Follow(ctx, req.Uid, req.ToUid)
}

func (s *RelationshipService) Unfollow(ctx context.Context, req *v1.FollowRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.uc.Unfollow(ctx, req.Uid, req.ToUid)
}

func (s *RelationshipService) GetFollowList(ctx context.Context, req *v1.ListRequest) (*v1.ListResponse, error) {
	list, err := s.uc.GetFollowList(ctx, req.Uid, req.LastId)
	if err != nil {
		return &v1.ListResponse{}, err
	}
	return &v1.ListResponse{List: listDO2DTO(list)}, nil
}

func (s *RelationshipService) GetFollowerList(ctx context.Context, req *v1.ListRequest) (*v1.ListResponse, error) {
	list, err := s.uc.GetFollowerList(ctx, req.Uid, req.LastId)
	if err != nil {
		return &v1.ListResponse{}, err
	}
	return &v1.ListResponse{List: listDO2DTO(list)}, nil
}

func (s *RelationshipService) GetRelation(ctx context.Context, req *v1.RelationRequest) (*v1.RelationResponse, error) {
	rm, err := s.uc.GetRelation(ctx, req.Uid, req.Uids)
	if err != nil {
		return &v1.RelationResponse{}, err
	}
	return &v1.RelationResponse{Rm: rmDO2DTO(rm)}, nil
}

func (s *RelationshipService) GetRelationCount(ctx context.Context, req *v1.BatchRequest) (*v1.CountResponse, error) {
	cm, err := s.uc.GetRelationCount(ctx, req.Uids)
	if err != nil {
		return &v1.CountResponse{}, err
	}
	return &v1.CountResponse{RelationCount: cmDO2DTO(cm)}, nil
}

func (s *RelationshipService) GetUsersFollow(ctx context.Context, req *v1.BatchRequest) (*v1.UserFollowResponse, error) {
	fm, err := s.uc.GetUsersFollow(ctx, req.Uids)
	if err != nil {
		return &v1.UserFollowResponse{}, err
	}
	return &v1.UserFollowResponse{Fm: fmDO2DTO(fm)}, nil
}
