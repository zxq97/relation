package relation

type ListItem struct {
	UID        int64 `json:"uid"`
	CreateTime int64 `json:"create_time"`
}

type CountItem struct {
	FollowCount   int64 `json:"follow_count"`
	FollowerCount int64 `json:"follower_count"`
}

func listDTO2DO(val *ListResponse) []*ListItem {
	if val == nil || len(val.ItemList) == 0 {
		return nil
	}
	list := make([]*ListItem, 0, len(val.ItemList))
	for _, v := range val.ItemList {
		list = append(list, &ListItem{
			UID:        v.Uid,
			CreateTime: v.CreateTime,
		})
	}
	return list
}

func CountDTO2DO(val *CountResponse) map[int64]*CountItem {
	if val == nil || len(val.RelationCount) == 0 {
		return nil
	}
	m := make(map[int64]*CountItem, len(val.RelationCount))
	for k, v := range val.RelationCount {
		m[k] = &CountItem{
			FollowCount:   v.FollowCount,
			FollowerCount: v.FollowerCount,
		}
	}
	return m
}
