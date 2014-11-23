// Copyright 2014 by caixw, All rights reserved.
// Use of s source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"bytes"
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/caixw/lib.go/orm/core"
)

const (
	and = iota
	or
)

// whereExpr语句部分
type whereExpr struct {
	cond     *bytes.Buffer
	condArgs []interface{}
}

// 重置所有状态为初始值。
func (w *whereExpr) Reset() {
	w.cond.Reset()
	w.condArgs = w.condArgs[:0]
}

// 所有where子句的构建，最终都调用此方法来写入实例中。
// op 与前一个语句的连接符号，可以是and或是or常量；
// cond 条件语句，值只能是占位符，不能直接写值；
// condArgs 占位符对应的值。
//  w := newwhereExpr(...)
//  w.build(and, "username=='abc'") // 错误：不能使用abc，只能使用？占位符。
//  w.build(and, "username=?", "abc") // 正确，将转换成: and username='abc'
func (w *whereExpr) build(op int, cond string, condArgs ...interface{}) *whereExpr {
	switch {
	case w.cond.Len() == 0:
		w.cond.WriteString(" WHERE(")
	case op == and:
		w.cond.WriteString(" AND(")
	case op == or:
		w.cond.WriteString(" OR(")
	default:
		panic("无效的op操作符")
	}

	w.cond.WriteString(cond)
	w.cond.WriteByte(')')

	w.condArgs = append(w.condArgs, condArgs...)

	return w
}

// 添加一条与前一条件语句关系为and的条件语句，若为第一条条件语句，则自动
// 忽略and关系。
func (w *whereExpr) And(cond string, condArgs ...interface{}) *whereExpr {
	return w.build(and, cond, condArgs...)
}

// 添加一条与前一条件语句关系为or的条件语句，若为第一条条件语句，则自动
// 忽略or关系。
func (w *whereExpr) Or(cond string, condArgs ...interface{}) *whereExpr {
	return w.build(or, cond, condArgs...)
}

// 添加In条件语句，与前一条的关系为And
//  w.AndIn("age", 55,56,57) ==> WHERE age in(55,56,57)
func (w *whereExpr) AndIn(col string, condArgs ...interface{}) *whereExpr {
	return w.in(and, col, condArgs...)
}

func (w *whereExpr) OrIn(col string, condArgs ...interface{}) *whereExpr {
	return w.in(or, col, condArgs...)
}

// whereExpr.AndIn()的别名
func (w *whereExpr) In(col string, condArgs ...interface{}) *whereExpr {
	return w.AndIn(col, condArgs...)
}

func (w *whereExpr) AndBetween(col string, start, end interface{}) *whereExpr {
	return w.between(and, col, start, end)
}

func (w *whereExpr) OrBetween(col string, start, end interface{}) *whereExpr {
	return w.between(or, col, start, end)
}

// whereExpr.AndBetween()的别名
func (w *whereExpr) Between(col string, start, end interface{}) *whereExpr {
	return w.AndBetween(col, start, end)
}

func (w *whereExpr) AndIsNull(col string) *whereExpr {
	return w.isNull(and, col)
}

func (w *whereExpr) OrIsNull(col string) *whereExpr {
	return w.isNull(or, col)
}

// whereExpr.AndIsNull()的别名
func (w *whereExpr) IsNull(col string) *whereExpr {
	return w.AndIsNull(col)
}

func (w *whereExpr) AndIsNotNull(col string) *whereExpr {
	return w.isNotNull(and, col)
}

func (w *whereExpr) OrIsNotNull(col string) *whereExpr {
	return w.isNotNull(or, col)
}

// whereExpr.AndIsNotNull()的别名
func (w *whereExpr) IsNotNull(col string) *whereExpr {
	return w.AndIsNotNull(col)
}

// where col in(v1,v2)语句的实现函数，供AndIn()和OrIn()函数调用。
func (w *whereExpr) in(op int, col string, condArgs ...interface{}) *whereExpr {
	if len(condArgs) <= 0 {
		panic("condArgs参数不能为空")
	}

	cond := bytes.NewBufferString(col)
	cond.WriteString(" IN(")
	cond.WriteString(strings.Repeat("?,", len(condArgs)))
	cond.Truncate(cond.Len() - 1) // 去掉最后的逗号
	cond.WriteByte(')')

	return w.build(op, cond.String(), condArgs...)
}

// where col between start and end 语句的实现函数，供AndBetween()和OrBetween()调用。
func (w *whereExpr) between(op int, col string, start, end interface{}) *whereExpr {
	return w.build(op, col+" BETWEEN ? AND ?", start, end)
}

// where col is null 语句的实现函数，供AndIsNull()和OrIsNull()调用。
func (w *whereExpr) isNull(op int, col string) *whereExpr {
	return w.build(op, col+" IS NULL")
}

// where col is not null 语句的实现函数，供AndIsNotNull()和OrIsNotNull()调用。
func (w *whereExpr) isNotNull(op int, col string) *whereExpr {
	return w.build(op, col+" IS NOT NULL")
}

// 获取条件语句的SQL格式，会将被引号包含的字符串替换成当前数据支持的符号
func (w *whereExpr) condString(db core.DB) string {
	return db.ReplaceQuote(w.cond.String())
}

// 用于产生sql的insert语句。
// 一般用法如下：
//  sql := newInsert(engine)
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

func newInsert(d core.DB) *Insert {
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
	i.table = i.db.ReplacePrefix(name)
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
		i.q.WriteString(i.db.ReplaceQuote(i.table))
		i.q.WriteByte('(')

		// 替换列名中的引号
		cols := i.db.ReplaceQuote(strings.Join(i.cols, ","))
		i.q.WriteString(cols)

		i.q.WriteString(") VALUES(")
		// 去掉上面的最后一个逗号
		placeholder := strings.Repeat("?,", len(i.cols))
		i.q.WriteString(placeholder[0 : len(placeholder)-1])
		i.q.WriteByte(')')
	}

	return i.q.String()
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
//  sql := newUpdate(db)
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

func newUpdate(d core.DB) *Update {
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
	u.table = u.db.ReplacePrefix(name)
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
		u.q.WriteString(u.db.ReplaceQuote((u.table)))
		u.q.WriteString(" SET ")
		for _, v := range u.cols {
			u.q.WriteString(u.db.ReplaceQuote(v))
			u.q.WriteString("=?,")
		}
		u.q.Truncate(u.q.Len() - 1)

		// where
		u.q.WriteString(u.condString(u.db))
	}

	return u.q.String()
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
//  sql := newDelete(e)
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

func newDelete(d core.DB) *Delete {
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
	d.table = d.db.ReplacePrefix(name)

	return d
}

func (d *Delete) sqlString(rebuild bool) string {
	if rebuild || d.q.Len() == 0 {
		d.q.Reset()

		d.q.WriteString("DELETE FROM ")
		d.q.WriteString(d.db.ReplaceQuote(d.table))

		// where
		d.q.WriteString(d.condString(d.db))
	}

	return d.q.String()
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
//  sql := newSelect(db)
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

func newSelect(d core.DB) *Select {
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
	s.table = s.db.ReplacePrefix(name)

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
	s.join.WriteString(s.db.ReplaceQuote(s.db.ReplacePrefix(table)))
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
		cols := s.db.ReplaceQuote(strings.Join(s.cols, ","))
		s.q.WriteString(cols)
		s.q.WriteString(" FROM ")
		s.q.WriteString(s.db.ReplaceQuote(s.table))
		s.q.WriteString(s.join.String())
		s.q.WriteString(s.condString(s.db)) // where
		s.q.WriteString(s.order.String())   // NOTE(caixw):mysql中若要limit，order字段是必须提供的
		s.q.WriteString(s.limitSQL)
	}

	return s.q.String()
}

// 将当前语句预编译并缓存到stmts中，方便之后再次使用。
func (s *Select) Stmt(name string) (*sql.Stmt, error) {
	return s.db.GetStmts().AddSQL(name, s.q.String())
}

// implement Fetch.Query()
func (s *Select) Query(args ...interface{}) (*sql.Rows, error) {
	if len(args) == 0 {
		// 与sqlString中添加的顺序相同，where在limit之前
		args = append(s.condArgs, s.limitArgs)
	}

	return s.db.Query(s.sqlString(false), args...)
}

// implement Fetch.QueryRow()
func (s *Select) QueryRow(args ...interface{}) *sql.Row {
	if len(args) == 0 {
		// 与sqlString中添加的顺序相同，where在limit之前
		args = append(s.condArgs, s.limitArgs)
	}

	return s.db.QueryRow(s.sqlString(false), args...)
}

// implement Fetch.Fetch2Map()
func (s *Select) Fetch2Map(args ...interface{}) (map[string]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := core.Fetch2Maps(true, rows)
	if err != nil {
		return nil, err
	}

	return data[0], nil
}

// implement Fetch.Fetch2Maps()
func (s *Select) Fetch2Maps(args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := core.Fetch2Maps(true, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// implement Fetch.FetchColumn()
func (s *Select) FetchColumn(col string, args ...interface{}) (interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := core.FetchColumns(true, col, rows)
	if err != nil {
		return nil, err
	}

	return data[0], nil
}

// implement Fetch.FetchColumns()
func (s *Select) FetchColumns(col string, args ...interface{}) ([]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := core.FetchColumns(true, col, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// implement Fetch.Fetch()
func (s *Select) Fetch(v interface{}, args ...interface{}) error {
	rows, err := s.Query(args...)
	if err != nil {
		return err
	}

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
			return errors.New("map的键值类型只能为interface{]")
		}

		mapped, err := s.Fetch2Map(args...)
		if err != nil {
			return err
		}
		vv.Set(reflect.ValueOf(mapped))
	case reflect.Slice:
		return core.Fetch2Objs(vv.Interface(), rows)
	case reflect.Struct:
		return core.Fetch2Objs(vv.Interface(), rows)
	default:
		return errors.New("只支持slice,map和struct指针")
	}
	return nil
}
