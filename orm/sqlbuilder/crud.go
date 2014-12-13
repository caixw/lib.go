// Copyright 2014 by caixw, All rights reserved.
// Use of s source code is governed by a MIT
// license that can be found in the LICENSE file.

package sqlbuilder

import (
	"bytes"
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/caixw/lib.go/orm/core"
	"github.com/caixw/lib.go/orm/fetch"
)

// 用于产生sql的insert语句。
// 一般用法如下：
//  sql := NewInsert(engine)
//
//  sql.Table(`"table.user"`).
//      Columns("username", "password", `"group"`).
//      Exec(nil, "admin", "123", 1)
//
//  sql.Table("user"). // 该表没有前缀
//      Data(map[string]interface{}{"username":"admin", "password":"123", `"group"`:1}).
//      Exec(nil) // 此处不指定参数，则表示直接使用上面Data函数指定的值。
type Insert struct {
	db    core.DB
	table string
	q     *bytes.Buffer
	cols  []string
	vals  []interface{}
}

func NewInsert(d core.DB) *Insert {
	ret := &Insert{
		db:   d,
		q:    bytes.NewBufferString(""),
		cols: make([]string, 0),
		vals: make([]interface{}, 0),
	}
	return ret
	//return ret.Reset()
}

// 重置表的所有状态。
func (i *Insert) Reset() {
	i.q.Reset()
	i.table = ""
	i.cols = i.cols[:0]
	i.vals = i.vals[:0]
}

// 指定操作的表名。
func (i *Insert) Table(name string) *Insert {
	i.table = name
	return i
}

// 指定多个列名。
// 不能将多个列名包含一个参数中，否则将在运行时出错。
func (i *Insert) Columns(cols ...string) *Insert {
	i.cols = append(i.cols, cols...)

	return i
}

// 添加一个键值对。
func (i *Insert) Add(col string, val interface{}) *Insert {
	i.cols = append(i.cols, col)
	i.vals = append(i.vals, val)

	return i
}

// 指定数据，相当于依次调用Set()函数
func (i *Insert) Data(data map[string]interface{}) *Insert {
	for c, v := range data {
		i.cols = append(i.cols, c)
		i.vals = append(i.vals, v)
	}

	return i
}

// 返回SQL语句。
func (i *Insert) sqlString(rebuild bool) string {
	if rebuild || i.q.Len() == 0 {
		i.q.Reset() // 清空之前的内容

		i.q.WriteString("INSERT INTO ")
		i.q.WriteString(i.table)

		i.q.WriteByte('(')
		i.q.WriteString(strings.Join(i.cols, ","))
		i.q.WriteString(") VALUES(")
		placeholder := strings.Repeat("?,", len(i.cols))
		// 去掉上面的最后一个逗号
		i.q.WriteString(placeholder[0 : len(placeholder)-1])
		i.q.WriteByte(')')
	}

	return i.db.PrepareSQL(i.q.String())
}

// 缓存当前语句到stmt
func (i *Insert) Stmt(name string) (*sql.Stmt, error) {
	return i.db.GetStmts().AddSQL(name, i.q.String())
}

// 执行当前的insert操作到数据库。若指定了args参数，则使用当前args参数
// 替换占位符，若不传递Args参数，则尝试使用Columns()等方法传递的值远的
// 占位符。
func (i *Insert) Exec(args ...interface{}) (sql.Result, error) {
	if len(args) == 0 { // 优先使用args参数，若没用，则调用i.vals中的值。
		args = i.vals
	}

	return i.db.Exec(i.sqlString(false), args...)
}

// 用于产生sql的update语句。
//  sql := NewUpdate(db)
//
//  sql.Table("table.user").
//      Columns("username",`"group"`).
//      Exec(tx, "admin2", 2) // 没用用AddValues()中指定的值，而是用了Exec()中的值
//
//  sql.Reset().
//      Table("table.user").
//      Data(map[string]interface{}{"username":"admin1",`"group"`:1}).
//      Exec(tx)
type Update struct {
	whereExpr

	db    core.DB
	table string
	q     *bytes.Buffer
	cols  []string
	vals  []interface{}
}

func NewUpdate(d core.DB) *Update {
	return &Update{
		whereExpr: whereExpr{
			cond:     bytes.NewBufferString(""),
			condArgs: make([]interface{}, 0),
		},
		db:   d,
		q:    bytes.NewBufferString(""),
		cols: make([]string, 0),
		vals: make([]interface{}, 0),
	}
}

func (u *Update) Reset() {
	u.whereExpr.Reset()
	u.q.Reset()
	u.table = ""
	u.cols = u.cols[0:0]
	u.vals = u.vals[0:0]
}

func (u *Update) Table(name string) *Update {
	u.table = name
	return u
}

func (u *Update) Columns(cols ...string) *Update {
	u.cols = append(u.cols, cols...)

	return u
}

func (u *Update) Data(data map[string]interface{}) *Update {
	for k, v := range data {
		u.Set(k, v)
	}
	return u
}

func (u *Update) Set(col string, val interface{}) *Update {
	u.cols = append(u.cols, col)
	u.vals = append(u.vals, val)

	return u
}

func (u *Update) sqlString(rebuild bool) string {
	if rebuild {
		u.q.Reset()

		u.q.WriteString("UPDATE ")
		u.q.WriteString(u.table)
		u.q.WriteString(" SET ")
		for _, v := range u.cols {
			u.q.WriteString(v)
			u.q.WriteString("=?,")
		}
		u.q.Truncate(u.q.Len() - 1)

		// where
		u.q.WriteString(u.cond.String())
	}

	return u.db.PrepareSQL(u.q.String())
}

func (u *Update) Stmt(name string) (*sql.Stmt, error) {
	return u.db.GetStmts().AddSQL(name, u.q.String())
}

func (u *Update) Exec(args ...interface{}) (sql.Result, error) {
	if len(args) == 0 {
		args = append(u.vals, u.whereExpr.condArgs)
	}

	return u.db.Exec(u.sqlString(false), args...)
}

// sql的delete语句
//  sql := NewDelete(e)
//  sql.Table("table.user").
//      And("username=?", "admin").
//      Or("group=?", 1).
//      Exec()
type Delete struct {
	whereExpr

	db    core.DB
	table string
	q     *bytes.Buffer
}

func NewDelete(d core.DB) *Delete {
	return &Delete{
		whereExpr: whereExpr{
			cond:     bytes.NewBufferString(""),
			condArgs: make([]interface{}, 0),
		},
		db: d,
		q:  bytes.NewBufferString(""),
	}
}

func (d *Delete) Reset() {
	d.table = ""
	d.whereExpr.Reset()
	d.q.Reset()
}

func (d *Delete) Table(name string) *Delete {
	d.table = name
	return d
}

func (d *Delete) sqlString(rebuild bool) string {
	if rebuild || d.q.Len() == 0 {
		d.q.Reset()

		d.q.WriteString("DELETE FROM ")
		d.q.WriteString(d.table)

		// where
		d.q.WriteString(d.cond.String())
	}

	return d.db.PrepareSQL(d.q.String())
}

func (d *Delete) Stmt(name string) (*sql.Stmt, error) {
	return d.db.GetStmts().AddSQL(name, d.q.String())
}

func (d *Delete) Exec(args ...interface{}) (sql.Result, error) {
	if len(args) == 0 {
		args = d.condArgs
	}

	return d.db.Exec(d.sqlString(false), args...)
}

// Select
//
//  sql := NewSelect(db)
//
//  sql.Select("username", "password", `"group"`).
//      From("table.users").
//      Where("id=?", 1)
//      And("username=?", 1)
//      Or(`"group"=?`, 1)
//      Fetch(tx)
type Select struct {
	whereExpr

	db        core.DB
	table     string
	q         *bytes.Buffer
	join      *bytes.Buffer
	order     *bytes.Buffer
	cols      []string
	limitSQL  string
	limitArgs []interface{}
}

func NewSelect(d core.DB) *Select {
	return &Select{
		whereExpr: whereExpr{},
		db:        d,
		q:         bytes.NewBufferString(""),
		// join: bytes.NewBufferString(""), 基本用不到的东西，就不要初始化了。
		order: bytes.NewBufferString(""),
		cols:  make([]string, 0),
	}
}

func (s *Select) Reset() {
	s.whereExpr.Reset()

	s.table = ""
	s.q.Reset()
	s.order.Reset()
	s.cols = s.cols[0:0]
	s.limitSQL = ""
	s.limitArgs = nil

	if s.join != nil {
		s.join.Reset()
	}
}

// 指定列名
func (s *Select) Columns(cols ...string) *Select {
	s.cols = append(s.cols, cols...)

	return s
}

// 指定表名
func (s *Select) Table(name string) *Select {
	s.table = name
	return s
}

// Select.Table的别名
func (s *Select) From(name string) *Select {
	return s.Table(name)
}

var joinType = []string{" LEFT JOIN ", " RIGHT JOIN ", " INNER JOIN ", " FULL JOIN "}

// join功能
// typ join的类型： 0 LEFT JOIN; 1 RIGHT JOIN; 2 INNER JOIN; 3 FULL JOIN
func (s *Select) joinOn(typ int, table string, on string) *Select {
	s.join.WriteString(joinType[typ])
	s.join.WriteString(table)
	s.join.WriteString(" ON ")
	s.join.WriteString(on)

	return s
}

// left join ... on ...
func (s *Select) LeftJoin(table, on string) *Select {
	return s.joinOn(0, table, on)
}

// right join ... on ...
func (s *Select) RightJoin(table, on string) *Select {
	return s.joinOn(1, table, on)
}

// inner join ... on ...
func (s *Select) InnerJoin(table, on string) *Select {
	return s.joinOn(2, table, on)
}

// full join ... on ...
func (s *Select) FullJoin(table, on string) *Select {
	return s.joinOn(3, table, on)
}

// Order by 功能
// sort 排序方式: 1=asc,2=desc，其它值无效
func (s *Select) orderBy(sort int, col string) *Select {
	if sort != 1 && sort != 2 {
		panic("sort参数错误，只能是1或是2")
	}

	if s.order.Len() == 0 {
		s.order.WriteString("ORDER BY ")
	} else {
		s.order.WriteString(", ")
	}

	s.order.WriteString(col)
	switch sort {
	case 1:
		s.order.WriteString("ASC ")
	case 2:
		s.order.WriteString("DESC ")
	default:
		panic("无效的order值，只能是asc或是desc")
	}

	return s
}

// order by asc
func (s *Select) Asc(cols ...string) *Select {
	for _, c := range cols {
		s.orderBy(1, c)
	}

	return s
}

// order by desc
func (s *Select) Desc(cols ...string) *Select {
	for _, c := range cols {
		s.orderBy(2, c)
	}

	return s
}

// limit ... offset ...
// offset值为0时，相当于limit N的效果。
func (s *Select) Limit(limit, offset int) *Select {
	s.limitSQL, s.limitArgs = s.db.Dialect().LimitSQL(limit, offset)
	return s
}

// 分页显示，调用Limit()实现，即Page()与Limit()方法会相互覆盖。
// start显示的页面，起始页为1，小于1时将触发panic；
// size每页显示的条数，小于1时将触发panic；
func (s *Select) Page(start, size int) *Select {
	if start < 1 {
		panic("start必须大于0")
	}
	if size < 1 {
		panic("size必须大于0")
	}

	start-- // 转到从0页开始
	return s.Limit(size, start*size)
}

// 将当前实例转换成SQL，所有的变量都以？代替。
// rebuild是否重新产生SQL，一般情况下，只有在调用sqlString()方法之后，
// 如果再修改了内容，需要再次调用sqlString()方法，并传递true值。
func (s *Select) sqlString(rebuild bool) string {
	if rebuild || s.q.Len() == 0 {
		s.q.Reset()

		s.q.WriteString("SELECT ")
		s.q.WriteString(strings.Join(s.cols, ","))
		s.q.WriteString(" FROM ")
		s.q.WriteString(s.table)
		s.q.WriteString(s.join.String())
		s.q.WriteString(s.cond.String())  // where
		s.q.WriteString(s.order.String()) // NOTE(caixw):mysql中若要limit，order字段是必须提供的
		s.q.WriteString(s.limitSQL)
	}

	return s.db.PrepareSQL(s.q.String())
}

// 将当前语句预编译并缓存到stmts中，方便之后再次使用。
func (s *Select) Stmt(name string) (*sql.Stmt, error) {
	return s.db.GetStmts().AddSQL(name, s.q.String())
}

// 导出数据到sql.Rows
func (s *Select) Query(args ...interface{}) (*sql.Rows, error) {
	if len(args) == 0 {
		// 与sqlString中添加的顺序相同，where在limit之前
		args = append(s.condArgs, s.limitArgs)
	}

	return s.db.Query(s.sqlString(false), args...)
}

// 导出第一条数据到sql.Row
func (s *Select) QueryRow(args ...interface{}) *sql.Row {
	if len(args) == 0 {
		// 与sqlString中添加的顺序相同，where在limit之前
		args = append(s.condArgs, s.limitArgs)
	}

	return s.db.QueryRow(s.sqlString(false), args...)
}

// 导出数据到map[string]interface{}
func (s *Select) Fetch2Map(args ...interface{}) (map[string]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := fetch.Map(true, rows)
	if err != nil {
		return nil, err
	}

	return data[0], nil
}

// 导出所有数据到[]map[string]interface{}
func (s *Select) Fetch2Maps(args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := fetch.Map(true, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 返回指定列的第一行内容
func (s *Select) FetchColumn(col string, args ...interface{}) (interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := fetch.Column(true, col, rows)
	if err != nil {
		return nil, err
	}

	return data[0], nil
}

// 返回指定列的所有数据
func (s *Select) FetchColumns(col string, args ...interface{}) ([]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := fetch.Column(false, col, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 将当前select语句查询的数据导出到v中
func (s *Select) Fetch(v interface{}, args ...interface{}) error {
	rows, err := s.Query(args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Map:
		if vv.Type().Key().Kind() != reflect.String {
			return errors.New("map的键名类型只能为string")
		}
		if vv.Type().Elem().Kind() != reflect.Interface {
			return errors.New("map的键值类型只能为interface{}")
		}

		mapped, err := s.Fetch2Map(args...)
		if err != nil {
			return err
		}
		vv.Set(reflect.ValueOf(mapped))
	case reflect.Slice:
		return fetch.Obj(vv.Interface(), rows)
	case reflect.Struct:
		return fetch.Obj(vv.Interface(), rows)
	default:
		return errors.New("只支持slice,map和struct指针")
	}
	return nil
}
