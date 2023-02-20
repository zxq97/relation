package method

import "gorm.io/gen"

type UserFollow interface {
	//sql(insert into user_follows (uid, to_uid) values (@uid, @touid))
	InsertFollow(uid, touid int64) error
	//sql(delete from user_follows where uid=@uid and to_uid=@uid limit 1)
	DeleteFollow(uid, touid int64) error
	//sql(select to_uid, create_at from user_follows where uid=@uid)
	FindUserFollow(uid int64) ([]*gen.T, error)
}

type UserFollower interface {
	//sql(insert into user_followers (uid, to_uid) values (@uid, @touid))
	InsertFollower(uid, touid int64) error
	//sql(delete from user_followers where uid=@uid and to_uid=@touid limit 1)
	DeleteFollower(uid, touid int64) error
	//sql(select to_uid, create_at from user_followers where uid=@uid order by create_at desc limit @limit)
	FindUserFollower(uid, limit int64) ([]*gen.T, error)
	//sql(select to_uid, create_at from user_followers where uid=@uid and id < (select id from user_followers where uid=@uid and to_uid=@lastid) order by create_at desc limit @limit)
	FindUserFollowerByLastID(uid, lastid, limit int64) ([]*gen.T, error)
}

type UserRelationCount interface {
	//sql(insert into user_relation_counts (uid, follow_count) values (@uid, 1) on duplicate key update follow_count=follow_count+1)
	IncrFollowCount(uid int64) error
	//sql(update user_relation_counts set follow_count=follow_count-1 where uid=@uid limit 1)
	DecrFollowCount(uid int64) error
	//sql(insert into user_relation_counts (uid, follower_count) values (@uid, @cnt) on duplicate key update follower_count=follower_count+@cnt)
	IncrByFollowerCount(uid, cnt int64) error
	//sql(update user_relation_counts set follower_count=follower_count-@cnt where uid=@uid limit 1)
	DecrByFollowerCount(uid, cnt int64) error
	//sql(select uid, follow_count, follower_count user_relation_counts where uid in (@uids))
	FindUsersRelationCount(uids []int64) ([]*gen.T, error)
}

type ExtraFollower interface {
	//sql(insert into extra_followers (uid) values (@uid))
	InsertFollower(uid int64) error
	//sql(insert into extra_followers (uid, stats) values (@uid, 1))
	DeleteFollower(uid int64) error
	//sql(select id, uid, stats from extra_followers limit @limit)
	FindUnSyncRecord(limit int64) ([]*gen.T, error)
	//sql(select id, uid, stats from extra_followers where uid=@uid)
	FindUnSyncRecordByUID(uid int64) ([]*gen.T, error)
	//sql(delete from extra_followers where id in (@ids))
	DeleteRecord(ids []int64) error
}
