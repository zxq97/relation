package service

import (
	"github.com/zxq97/relation/api/relationship/service/v1"
	"github.com/zxq97/relation/app/relationship/service/internal/biz"
)

func listDO2DTO(list []*biz.FollowItem) *v1.FollowList {
	l := make([]*v1.FollowItem, len(list))
	for k, v := range list {
		l[k] = &v1.FollowItem{
			ToUid:      v.ToUid,
			CreateTime: v.CreateTime,
		}
	}
	return &v1.FollowList{List: l}
}

func rmDO2DTO(rm map[int64]*biz.UserRelation) map[int64]*v1.RelationItem {
	m := make(map[int64]*v1.RelationItem, len(rm))
	for k, v := range rm {
		m[k] = &v1.RelationItem{
			Relation:     v.Relation,
			FollowTime:   v.FollowTime,
			FollowedTime: v.FollowedTime,
		}
	}
	return m
}

func cmDO2DTO(cm map[int64]*biz.RelationCount) map[int64]*v1.RelationCount {
	m := make(map[int64]*v1.RelationCount, len(cm))
	for k, v := range cm {
		m[k] = &v1.RelationCount{
			FollowCount:   v.FollowCount,
			FollowerCount: v.FollowerCount,
		}
	}
	return m
}

func fmDO2DTO(fm map[int64][]*biz.FollowItem) map[int64]*v1.FollowList {
	m := make(map[int64]*v1.FollowList, len(fm))
	for k, v := range fm {
		m[k] = listDO2DTO(v)
	}
	return m
}
