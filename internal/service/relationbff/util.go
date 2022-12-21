package relationbff

import (
	"github.com/zxq97/relation/internal/model"
	"github.com/zxq97/relation/internal/service/relationsvc"
	"github.com/zxq97/relation/pkg/relation"
)

func translateCount(val *relationsvc.CountResponse) *relation.CountResponse {
	if val == nil {
		return &relation.CountResponse{}
	}
	m := make(map[int64]*relation.RelationCount, len(val.RelationCount))
	for k, v := range val.RelationCount {
		m[k] = &relation.RelationCount{
			FollowCount:   v.FollowCount,
			FollowerCount: v.FollowerCount,
		}
	}
	return &relation.CountResponse{RelationCount: m}
}

func translateList(val *model.FollowList) *relation.ListResponse {
	if val == nil {
		return &relation.ListResponse{}
	}
	list := make([]*relation.FollowItem, 0, len(val.List))
	for _, v := range val.List {
		list = append(list, &relation.FollowItem{
			Uid:        v.ToUid,
			CreateTime: v.CreateTime,
		})
	}
	return &relation.ListResponse{ItemList: list}
}

func translateRelation(val *relationsvc.RelationResponse) *relation.RelationResponse {
	if val == nil {
		return nil
	}
	m := make(map[int64]*relation.RelationItem, len(val.Rm))
	for k, v := range val.Rm {
		m[k].Relation = v.Relation
		m[k].FollowTime = v.FollowTime
		m[k].FollowedTime = v.FollowedTime
	}
	return &relation.RelationResponse{Rm: m}
}
