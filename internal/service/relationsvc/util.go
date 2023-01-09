package relationsvc

import "github.com/zxq97/relation/internal/biz"

func listDO2DTO(list []*biz.FollowItem) *FollowList {
	l := make([]*FollowItem, len(list))
	for k, v := range list {
		l[k] = &FollowItem{
			ToUid:      v.ToUID,
			CreateTime: v.CreateTime,
		}
	}
	return &FollowList{List: l}
}

func lmDO2DTO(lm map[int64][]*biz.FollowItem) *UserFollowResponse {
	m := make(map[int64]*FollowList, len(lm))
	for k, v := range lm {
		m[k] = listDO2DTO(v)
	}
	return &UserFollowResponse{Fm: m}
}

func rmDO2DTO(rm map[int64]*biz.UserRelation) *RelationResponse {
	m := make(map[int64]*RelationItem, len(rm))
	for k, v := range rm {
		m[k] = &RelationItem{
			Relation:     v.Relation,
			FollowTime:   v.FollowTime,
			FollowedTime: v.FollowedTime,
		}
	}
	return &RelationResponse{Rm: m}
}

func cmDO2DTO(cm map[int64]*biz.UserFollowCount) *CountResponse {
	m := make(map[int64]*RelationCount, len(cm))
	for k, v := range cm {
		m[k] = &RelationCount{
			FollowCount:   v.FollowCount,
			FollowerCount: v.FollowerCount,
		}
	}
	return &CountResponse{RelationCount: m}
}
