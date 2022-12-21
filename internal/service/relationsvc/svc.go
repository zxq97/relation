package relationsvc

import (
	"context"

	"github.com/zxq97/relation/internal/biz"
	"github.com/zxq97/relation/internal/data"
)

type RelationSvc struct {
	repo biz.RelationSVCRepo
}

func InitRelationSvc(conf *RelationSvcConfig) error {
	repo, err := data.NewRelationSVCRepo(conf.Redis["redis"], conf.MC["mc"], conf.Mysql["relation"], conf.Kafka["kafka"].Addr)
	if err != nil {
		return err
	}

}

//Follow 关注
func (svc *RelationSvc) Follow(ctx context.Context, req *FollowRequest) (*EmptyResponse, error) {
	err := svc.repo.Follow(ctx, req.Uid, req.ToUid)
	if err != nil {

	}
}

//Unfollow 取关
func (RelationSvc) Unfollow(ctx context.Context, req *FollowRequest) (*EmptyResponse, error) {

}

//GetFollowList 关注列表
func (RelationSvc) GetFollowList(ctx context.Context, req *ListRequest) (*model.FollowList, error) {

}

//GetFollowerList 粉丝列表
func (RelationSvc) GetFollowerList(ctx context.Context, req *ListRequest) (*model.FollowList, error) {

}

//GetRelation 好有关系
func (RelationSvc) GetRelation(ctx context.Context, req *RelationRequest) (*RelationResponse, error) {

}

//GetRelationCount 关注 粉丝 数量
func (RelationSvc) GetRelationCount(ctx context.Context, req *BatchRequest) (*CountResponse, error) {

}

//GetUsersFollow 获取全量关注
func (RelationSvc) GetUsersFollow(ctx context.Context, req *BatchRequest) (*UserFollowResponse, error) {

}
