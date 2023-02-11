package data

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zxq97/gokit/pkg/cast"
	"github.com/zxq97/relation/app/relationship/service/internal/biz"
)

func (r *relationshipRepo) getRelationCount(ctx context.Context, uids []int64) (map[int64]*biz.RelationCount, error) {
	rcs := []*UserRelationCount{}
	err := r.db.WithContext(ctx).Model(&UserRelationCount{}).Where("`uid` IN ?", uids).Find(&rcs).Error
	if err != nil {
		return nil, errors.Wrap(err, "db get relation count")
	}
	m := make(map[int64]*biz.RelationCount, len(uids))
	for _, v := range rcs {
		m[v.UID] = &biz.RelationCount{
			Uid:           v.UID,
			FollowCount:   v.FollowCount,
			FollowerCount: v.FollowerCount,
		}
	}
	return m, nil
}

func (r *relationshipRepo) getUserFollow(ctx context.Context, uid int64) ([]*biz.FollowItem, error) {
	uf := []*UserFollow{}
	err := r.db.WithContext(ctx).Model(&UserFollow{}).Select("to_uid", "create_at").Where("uid = ?", uid).Order("create_at DESC").Find(&uf).Error
	if err != nil {
		return nil, errors.Wrap(err, "db get user follow")
	}
	list := make([]*biz.FollowItem, len(uf))
	for k, v := range uf {
		list[k] = &biz.FollowItem{
			ToUid:      v.ToUID,
			CreateTime: v.CreateAt.UnixMilli(),
		}
	}
	return list, nil
}

func (r *relationshipRepo) getUserFollower(ctx context.Context, uid, lastid int64) ([]*biz.FollowItem, error) {
	uf := []*UserFollower{}
	filter := "`uid` = ?"
	if lastid != 0 {
		sql := "SELECT `id` FROM user_followers WHERE `uid` = ? AND `to_uid` = ? LIMIT 1"
		var id int64
		err := r.db.WithContext(ctx).Raw(sql, uid, lastid).Scan(&id).Error
		if err != nil {
			return nil, errors.Wrap(err, "get last id")
		}
		filter += " AND `id` < " + cast.FormatInt(id)
	}
	err := r.db.WithContext(ctx).Select("to_uid", "create_at").Where(filter, uid).Order("create_at DESC").Limit(20).Find(&uf).Error
	if err != nil {
		return nil, errors.Wrap(err, "db get user follower")
	}
	list := make([]*biz.FollowItem, len(uf))
	for k, v := range uf {
		list[k] = &biz.FollowItem{
			ToUid:      v.ToUID,
			CreateTime: v.CreateAt.UnixMilli(),
		}
	}
	return list, nil
}
