package service

import (
	"context"

	"github.com/zxq97/relation/api/black/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BlackService struct {
	//v1.UnimplementedBlackSvcServer
}

func NewBlackService() *BlackService {
	return &BlackService{}
}

func (s *BlackService) Black(ctx context.Context, req *v1.BlackRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *BlackService) CancelBlack(ctx context.Context, req *v1.BlackRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *BlackService) GetBlackList(ctx context.Context, req *v1.ListRequest) (*v1.ListResponse, error) {
	return &v1.ListResponse{}, nil
}

func (s *BlackService) CheckBlacked(ctx context.Context, req *v1.CheckRequest) (*v1.CheckResponse, error) {
	return &v1.CheckResponse{}, nil
}
