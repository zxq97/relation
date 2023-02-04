package data

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/zxq97/relation/app/relationship/job/internal/biz"
)

const (
	mcKeyUserFollow = "rla_fow_%d" // uid
)

func (r *relationshipRepo) addFollowCacheL1(ctx context.Context, uid int64, itemList []*biz.FollowItem) error {
	key := fmt.Sprintf(mcKeyUserFollow, uid)
	val, err := r.mc.Get(key)
	if err != nil {
		return err
	}
	list := biz.FollowList{}
	err = proto.Unmarshal(val.Value, &list)
	if err != nil {
		return err
	}
	for _, v := range itemList {
		list.List = append(list.List, v)
	}
	bs, err := proto.Marshal(&list)
	if err != nil {
		return err
	}
	val.Value = bs
	return r.mc.CompareAndSwap(val)
}

func (r *relationshipRepo) delFollowCacheL1(ctx context.Context, uid, touid int64) error {
	key := fmt.Sprintf(mcKeyUserFollow, uid)
	val, err := r.mc.Get(key)
	if err != nil {
		return err
	}
	list := biz.FollowList{}
	err = proto.Unmarshal(val.Value, &list)
	if err != nil {
		return err
	}
	var flag bool
	for k, v := range list.List {
		if v.ToUid == touid {
			list.List = append(list.List[0:k], list.List[k:]...)
			flag = true
			break
		}
	}
	if !flag {
		return nil
	}
	bs, err := proto.Marshal(&list)
	if err != nil {
		return err
	}
	val.Value = bs
	return r.mc.CompareAndSwap(val)
}
