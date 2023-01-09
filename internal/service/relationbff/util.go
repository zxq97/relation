package relationbff

import (
	"github.com/zxq97/relation/api"
	"github.com/zxq97/relation/internal/service/relationsvc"
)

func translateCount(val *relationsvc.CountResponse) *api.CountResponse {
	if val == nil {
		return &api.CountResponse{}
	}
	m := make(map[int64]*api.RelationCount, len(val.RelationCount))
	for k, v := range val.RelationCount {
		m[k] = &api.RelationCount{
			FollowCount:   v.FollowCount,
			FollowerCount: v.FollowerCount,
		}
	}
	return &api.CountResponse{RelationCount: m}
}

func translateList(val *model.FollowList) *api.ListResponse {
	if val == nil {
		return &api.ListResponse{}
	}
	list := make([]*api.FollowItem, 0, len(val.List))
	for _, v := range val.List {
		list = append(list, &api.FollowItem{
			Uid:        v.ToUid,
			CreateTime: v.CreateTime,
		})
	}
	return &api.ListResponse{ItemList: list}
}

func translateRelation(val *relationsvc.RelationResponse) *api.RelationResponse {
	if val == nil {
		return nil
	}
	m := make(map[int64]*api.RelationItem, len(val.Rm))
	for k, v := range val.Rm {
		m[k].Relation = v.Relation
		m[k].FollowTime = v.FollowTime
		m[k].FollowedTime = v.FollowedTime
	}
	return &api.RelationResponse{Rm: m}
}
