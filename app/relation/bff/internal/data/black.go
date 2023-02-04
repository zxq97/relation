package data

import (
	blackv1 "github.com/zxq97/relation/api/black/service/v1"
	relationshipv1 "github.com/zxq97/relation/api/relationship/service/v1"
	"github.com/zxq97/relation/app/relation/bff/internal/biz"
)

var _ biz.BlackRepo = (*blackRepo)(nil)

type blackRepo struct {
	relationshipClient relationshipv1.RelationSvcClient
	blackClient        blackv1.BlackSvcClient
}

func NewBlackRepo(relationshipClient relationshipv1.RelationSvcClient, blackClient blackv1.BlackSvcClient) *blackRepo {
	return &blackRepo{
		relationshipClient: relationshipClient,
		blackClient:        blackClient,
	}
}
