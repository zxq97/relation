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

func newUserFollow(db *gorm.DB, opts ...gen.DOOption) userFollow {
	_userFollow := userFollow{}

	_userFollow.userFollowDo.UseDB(db, opts...)
	_userFollow.userFollowDo.UseModel(&model.UserFollow{})

	tableName := _userFollow.userFollowDo.TableName()
	_userFollow.ALL = field.NewAsterisk(tableName)
	_userFollow.ID = field.NewInt64(tableName, "id")
	_userFollow.UID = field.NewInt64(tableName, "uid")
	_userFollow.ToUID = field.NewInt64(tableName, "to_uid")
	_userFollow.CreateAt = field.NewTime(tableName, "create_at")
	_userFollow.UpdateAt = field.NewTime(tableName, "update_at")

	_userFollow.fillFieldMap()

	return _userFollow
}

type userFollow struct {
	userFollowDo userFollowDo

	ALL      field.Asterisk
	ID       field.Int64
	UID      field.Int64
	ToUID    field.Int64
	CreateAt field.Time
	UpdateAt field.Time

	fieldMap map[string]field.Expr
}

func (u userFollow) Table(newTableName string) *userFollow {
	u.userFollowDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u userFollow) As(alias string) *userFollow {
	u.userFollowDo.DO = *(u.userFollowDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *userFollow) updateTableName(table string) *userFollow {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewInt64(table, "id")
	u.UID = field.NewInt64(table, "uid")
	u.ToUID = field.NewInt64(table, "to_uid")
	u.CreateAt = field.NewTime(table, "create_at")
	u.UpdateAt = field.NewTime(table, "update_at")

	u.fillFieldMap()

	return u
}

func (u *userFollow) WithContext(ctx context.Context) *userFollowDo {
	return u.userFollowDo.WithContext(ctx)
}

func (u userFollow) TableName() string { return u.userFollowDo.TableName() }

func (u userFollow) Alias() string { return u.userFollowDo.Alias() }

func (u *userFollow) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *userFollow) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 5)
	u.fieldMap["id"] = u.ID
	u.fieldMap["uid"] = u.UID
	u.fieldMap["to_uid"] = u.ToUID
	u.fieldMap["create_at"] = u.CreateAt
	u.fieldMap["update_at"] = u.UpdateAt
}

func (u userFollow) clone(db *gorm.DB) userFollow {
	u.userFollowDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u userFollow) replaceDB(db *gorm.DB) userFollow {
	u.userFollowDo.ReplaceDB(db)
	return u
}

type userFollowDo struct{ gen.DO }

// sql(insert into user_follows (uid, to_uid) values (@uid, @touid))
func (u userFollowDo) InsertFollow(uid int64, touid int64) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, uid)
	params = append(params, touid)
	generateSQL.WriteString("insert into user_follows (uid, to_uid) values (?, ?) ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// sql(delete from user_follows where uid=@uid and to_uid=@uid limit 1)
func (u userFollowDo) DeleteFollow(uid int64, touid int64) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, uid)
	params = append(params, uid)
	generateSQL.WriteString("delete from user_follows where uid=? and to_uid=? limit 1 ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// sql(select to_uid, create_at from user_follows where uid=@uid)
func (u userFollowDo) FindUserFollow(uid int64) (result []*model.UserFollow, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, uid)
	generateSQL.WriteString("select to_uid, create_at from user_follows where uid=? ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

func (u userFollowDo) Debug() *userFollowDo {
	return u.withDO(u.DO.Debug())
}

func (u userFollowDo) WithContext(ctx context.Context) *userFollowDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userFollowDo) ReadDB() *userFollowDo {
	return u.Clauses(dbresolver.Read)
}

func (u userFollowDo) WriteDB() *userFollowDo {
	return u.Clauses(dbresolver.Write)
}

func (u userFollowDo) Session(config *gorm.Session) *userFollowDo {
	return u.withDO(u.DO.Session(config))
}

func (u userFollowDo) Clauses(conds ...clause.Expression) *userFollowDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userFollowDo) Returning(value interface{}, columns ...string) *userFollowDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userFollowDo) Not(conds ...gen.Condition) *userFollowDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userFollowDo) Or(conds ...gen.Condition) *userFollowDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userFollowDo) Select(conds ...field.Expr) *userFollowDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userFollowDo) Where(conds ...gen.Condition) *userFollowDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userFollowDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *userFollowDo {
	return u.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (u userFollowDo) Order(conds ...field.Expr) *userFollowDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userFollowDo) Distinct(cols ...field.Expr) *userFollowDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userFollowDo) Omit(cols ...field.Expr) *userFollowDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userFollowDo) Join(table schema.Tabler, on ...field.Expr) *userFollowDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userFollowDo) LeftJoin(table schema.Tabler, on ...field.Expr) *userFollowDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userFollowDo) RightJoin(table schema.Tabler, on ...field.Expr) *userFollowDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userFollowDo) Group(cols ...field.Expr) *userFollowDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userFollowDo) Having(conds ...gen.Condition) *userFollowDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userFollowDo) Limit(limit int) *userFollowDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userFollowDo) Offset(offset int) *userFollowDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userFollowDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *userFollowDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userFollowDo) Unscoped() *userFollowDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userFollowDo) Create(values ...*model.UserFollow) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userFollowDo) CreateInBatches(values []*model.UserFollow, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userFollowDo) Save(values ...*model.UserFollow) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userFollowDo) First() (*model.UserFollow, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserFollow), nil
	}
}

func (u userFollowDo) Take() (*model.UserFollow, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserFollow), nil
	}
}

func (u userFollowDo) Last() (*model.UserFollow, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserFollow), nil
	}
}

func (u userFollowDo) Find() ([]*model.UserFollow, error) {
	result, err := u.DO.Find()
	return result.([]*model.UserFollow), err
}

func (u userFollowDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.UserFollow, err error) {
	buf := make([]*model.UserFollow, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userFollowDo) FindInBatches(result *[]*model.UserFollow, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userFollowDo) Attrs(attrs ...field.AssignExpr) *userFollowDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userFollowDo) Assign(attrs ...field.AssignExpr) *userFollowDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userFollowDo) Joins(fields ...field.RelationField) *userFollowDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userFollowDo) Preload(fields ...field.RelationField) *userFollowDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userFollowDo) FirstOrInit() (*model.UserFollow, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserFollow), nil
	}
}

func (u userFollowDo) FirstOrCreate() (*model.UserFollow, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.UserFollow), nil
	}
}

func (u userFollowDo) FindByPage(offset int, limit int) (result []*model.UserFollow, count int64, err error) {
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

func (u userFollowDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userFollowDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userFollowDo) Delete(models ...*model.UserFollow) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userFollowDo) withDO(do gen.Dao) *userFollowDo {
	u.DO = *do.(*gen.DO)
	return u
}
