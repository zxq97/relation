// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/zxq97/relation/app/relationship/pkg/dal/model"
)

func newUserFollower(db *gorm.DB, opts ...gen.DOOption) userFollower {
	_userFollower := userFollower{}

	_userFollower.userFollowerDo.UseDB(db, opts...)
	_userFollower.userFollowerDo.UseModel(&model.UserFollower{})

	tableName := _userFollower.userFollowerDo.TableName()
	_userFollower.ALL = field.NewAsterisk(tableName)
	_userFollower.ID = field.NewInt64(tableName, "id")
	_userFollower.UID = field.NewInt64(tableName, "uid")
	_userFollower.ToUID = field.NewInt64(tableName, "to_uid")
	_userFollower.CreateAt = field.NewTime(tableName, "create_at")
	_userFollower.UpdateAt = field.NewTime(tableName, "update_at")

	_userFollower.fillFieldMap()

	return _userFollower
}

type userFollower struct {
	userFollowerDo userFollowerDo

	ALL      field.Asterisk
	ID       field.Int64
	UID      field.Int64
	ToUID    field.Int64
	CreateAt field.Time
	UpdateAt field.Time

	fieldMap map[string]field.Expr
}

func (u userFollower) Table(newTableName string) *userFollower {
	u.userFollowerDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userFollower) As(alias string) *userFollower {
	u.userFollowerDo.DO = *(u.userFollowerDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userFollower) updateTableName(table string) *userFollower {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewInt64(table, "id")
	u.UID = field.NewInt64(table, "uid")
	u.ToUID = field.NewInt64(table, "to_uid")
	u.CreateAt = field.NewTime(table, "create_at")
	u.UpdateAt = field.NewTime(table, "update_at")

	u.fillFieldMap()

	return u
}

func (u *userFollower) WithContext(ctx context.Context) *userFollowerDo {
	return u.userFollowerDo.WithContext(ctx)
}

func (u userFollower) TableName() string { return u.userFollowerDo.TableName() }

func (u userFollower) Alias() string { return u.userFollowerDo.Alias() }

func (u *userFollower) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userFollower) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 5)
	u.fieldMap["id"] = u.ID
	u.fieldMap["uid"] = u.UID
	u.fieldMap["to_uid"] = u.ToUID
	u.fieldMap["create_at"] = u.CreateAt
	u.fieldMap["update_at"] = u.UpdateAt
}

func (u userFollower) clone(db *gorm.DB) userFollower {
	u.userFollowerDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userFollower) replaceDB(db *gorm.DB) userFollower {
	u.userFollowerDo.ReplaceDB(db)
	return u
}

type userFollowerDo struct{ gen.DO }

// sql(insert into user_followers (uid, to_uid) values (@uid, @touid))
func (u userFollowerDo) InsertFollower(uid int64, touid int64) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, uid)
	params = append(params, touid)
	generateSQL.WriteString("insert into user_followers (uid, to_uid) values (?, ?) ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// sql(delete from user_followers where uid=@uid and to_uid=@touid limit 1)
func (u userFollowerDo) DeleteFollower(uid int64, touid int64) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, uid)
	params = append(params, touid)
	generateSQL.WriteString("delete from user_followers where uid=? and to_uid=? limit 1 ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// sql(select to_uid, create_at from user_followers where uid=@uid order by create_at desc limit @limit)
func (u userFollowerDo) FindUserFollower(uid int64, limit int64) (result []*model.UserFollower, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, uid)
	params = append(params, limit)
	generateSQL.WriteString("select to_uid, create_at from user_followers where uid=? order by create_at desc limit ? ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// sql(select to_uid, create_at from user_followers where uid=@uid and id < (select id from user_followers where uid=@uid and to_uid=@lastid) order by create_at desc limit @limit)
func (u userFollowerDo) FindUserFollowerByLastID(uid int64, lastid int64, limit int64) (result []*model.UserFollower, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, uid)
	params = append(params, uid)
	params = append(params, lastid)
	params = append(params, limit)
	generateSQL.WriteString("select to_uid, create_at from user_followers where uid=? and id < (select id from user_followers where uid=? and to_uid=?) order by create_at desc limit ? ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (u userFollowerDo) Debug() *userFollowerDo {
	return u.withDO(u.DO.Debug())
}

func (u userFollowerDo) WithContext(ctx context.Context) *userFollowerDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userFollowerDo) ReadDB() *userFollowerDo {
	return u.Clauses(dbresolver.Read)
}

func (u userFollowerDo) WriteDB() *userFollowerDo {
	return u.Clauses(dbresolver.Write)
}

func (u userFollowerDo) Session(config *gorm.Session) *userFollowerDo {
	return u.withDO(u.DO.Session(config))
}

func (u userFollowerDo) Clauses(conds ...clause.Expression) *userFollowerDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userFollowerDo) Returning(value interface{}, columns ...string) *userFollowerDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userFollowerDo) Not(conds ...gen.Condition) *userFollowerDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userFollowerDo) Or(conds ...gen.Condition) *userFollowerDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userFollowerDo) Select(conds ...field.Expr) *userFollowerDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userFollowerDo) Where(conds ...gen.Condition) *userFollowerDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userFollowerDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *userFollowerDo {
	return u.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (u userFollowerDo) Order(conds ...field.Expr) *userFollowerDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userFollowerDo) Distinct(cols ...field.Expr) *userFollowerDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userFollowerDo) Omit(cols ...field.Expr) *userFollowerDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userFollowerDo) Join(table schema.Tabler, on ...field.Expr) *userFollowerDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userFollowerDo) LeftJoin(table schema.Tabler, on ...field.Expr) *userFollowerDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userFollowerDo) RightJoin(table schema.Tabler, on ...field.Expr) *userFollowerDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userFollowerDo) Group(cols ...field.Expr) *userFollowerDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userFollowerDo) Having(conds ...gen.Condition) *userFollowerDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userFollowerDo) Limit(limit int) *userFollowerDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userFollowerDo) Offset(offset int) *userFollowerDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userFollowerDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *userFollowerDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userFollowerDo) Unscoped() *userFollowerDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userFollowerDo) Create(values ...*model.UserFollower) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userFollowerDo) CreateInBatches(values []*model.UserFollower, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userFollowerDo) Save(values ...*model.UserFollower) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userFollowerDo) First() (*model.UserFollower, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserFollower), nil
	}
}

func (u userFollowerDo) Take() (*model.UserFollower, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserFollower), nil
	}
}

func (u userFollowerDo) Last() (*model.UserFollower, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserFollower), nil
	}
}

func (u userFollowerDo) Find() ([]*model.UserFollower, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserFollower), err
}

func (u userFollowerDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserFollower, err error) {
	buf := make([]*model.UserFollower, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userFollowerDo) FindInBatches(result *[]*model.UserFollower, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userFollowerDo) Attrs(attrs ...field.AssignExpr) *userFollowerDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userFollowerDo) Assign(attrs ...field.AssignExpr) *userFollowerDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userFollowerDo) Joins(fields ...field.RelationField) *userFollowerDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userFollowerDo) Preload(fields ...field.RelationField) *userFollowerDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userFollowerDo) FirstOrInit() (*model.UserFollower, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserFollower), nil
	}
}

func (u userFollowerDo) FirstOrCreate() (*model.UserFollower, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserFollower), nil
	}
}

func (u userFollowerDo) FindByPage(offset int, limit int) (result []*model.UserFollower, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userFollowerDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userFollowerDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userFollowerDo) Delete(models ...*model.UserFollower) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userFollowerDo) withDO(do gen.Dao) *userFollowerDo {
	u.DO = *do.(*gen.DO)
	return u
}