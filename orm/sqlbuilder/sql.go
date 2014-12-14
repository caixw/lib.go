// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sqlbuilder

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/caixw/lib.go/orm/core"
)

const (
	Delete = iota
	Insert
	Update
	Select
)

// sql := sqlbuild.New()
// sql.Table("#user").
//     Where("id>?",5).
//     And("username like ?", "%admin%").
type SQL struct {
	db        core.DB
	tableName string
	errors    []error       // 所有的错误缓存
	buf       *bytes.Buffer // 语句缓存

	// where
	cond     *bytes.Buffer
	condArgs []interface{}

	// data
	cols []string
	vals []interface{}

	// select
	join      *bytes.Buffer
	order     *bytes.Buffer
	limitSQL  string
	limitArgs []interface{}
}

// 新建一个SQL实例。
func New(db core.DB) *SQL {
	return &SQL{
		db:     db,
		errors: []error{},
		buf:    bytes.NewBuffer([]byte{}),

		// where
		cond:     bytes.NewBuffer([]byte{}),
		condArgs: []interface{}{},

		// data
		cols: []string{},
		vals: []interface{}{},

		// select
		join:      bytes.NewBuffer([]byte{}),
		order:     bytes.NewBuffer([]byte{}),
		limitArgs: []interface{}{},
	}
}

// 重置SQL语句的状态。除了SQL.db以外，
// 其它属性都将被重围为初始状态。
func (s *SQL) Reset() {
	s.tableName = ""
	s.errors = s.errors[:0]
	s.buf.Reset()

	// where
	s.resetWhere()

	// data
	s.cols = s.cols[:0]
	s.vals = s.vals[:0]

	// select
	s.join.Reset()
	s.order.Reset()
	s.limitArgs = s.limitArgs[:0]
}

// 是否存在错误
func (s *SQL) HasErrors() bool {
	return len(s.errors) > 0
}

// 返回错误内容
func (s *SQL) Errors() []error {
	return s.errors
}

// 设置表名
// 多次调用，只有最后一次启作用。
func (s *SQL) Table(name string) *SQL {
	s.tableName = name
	return s
}

// 指定列名。
// update/insert 语句可以用此方法指定需要更新的列。
// 若需要指定数据，请使用Data()或是Add()方法；
//
// select 语句可以用此方法指定需要获取的列。
func (s *SQL) Columns(cols ...string) *SQL {
	s.cols = append(s.cols, cols...)

	return s
}

// update/insert 语句可以用此方法批量指定需要更新的字段及相应的数据，
// 其它语句，忽略此方法产生的数据。
func (s *SQL) Data(data map[string]interface{}) *SQL {
	for k, v := range data {
		s.cols = append(s.cols, k)
		s.vals = append(s.vals, v)
	}
	return s
}

// update/insert 语句可以用此方法指定一条需要更新的字段及相应的数据。
// 其它语句，忽略此方法产生的数据。
func (s *SQL) Add(col string, val interface{}) *SQL {
	s.cols = append(s.cols, col)
	s.vals = append(s.vals, val)

	return s
}

// 将当前语句预编译并缓存到stmts中，方便之后再次使用。
func (s *SQL) Stmt(action int, name string) (*sql.Stmt, error) {
	var sql string
	switch action {
	case Delete:
		sql = s.db.PrepareSQL(s.deleteSQL())
	case Update:
		sql = s.db.PrepareSQL(s.updateSQL())
	case Insert:
		sql = s.db.PrepareSQL(s.insertSQL())
	case Select:
		sql = s.db.PrepareSQL(s.selectSQL())
	default:
		return nil, fmt.Errorf("无效的的action值[%v]", action)
	}

	return s.db.GetStmts().AddSQL(name, sql)
}

// 执行当前语句。
func (s *SQL) Exec(action int, args ...interface{}) (sql.Result, error) {
	switch action {
	case Delete:
		return s.Delete(args...)
	case Update:
		return s.Update(args...)
	case Insert:
		return s.Insert(args...)
	case Select:
		return nil, errors.New("select语句不能使用Exec()方法执行")
	default:
		return nil, fmt.Errorf("无效的的action值[%v]", action)
	}
}

func (s *SQL) deleteSQL() string {
	s.buf.Reset()
	s.buf.WriteString("DELETE FROM ")
	s.buf.WriteString(s.tableName)

	// where
	s.buf.WriteString(s.cond.String())

	return s.db.PrepareSQL(s.buf.String())
}

func (s *SQL) Delete(args ...interface{}) (sql.Result, error) {
	if s.HasErrors() {
		return nil, Errors(s.errors)
	}

	if len(args) == 0 {
		args = s.condArgs
	}

	return s.db.Exec(s.deleteSQL(), args...)
}

func (s *SQL) updateSQL() string {
	s.buf.Reset()
	s.buf.WriteString("UPDATE ")
	s.buf.WriteString(s.tableName)
	s.buf.WriteString(" SET ")
	for _, v := range s.cols {
		s.buf.WriteString(v)
		s.buf.WriteString("=?,")
	}
	s.buf.Truncate(s.buf.Len() - 1)

	// where
	s.buf.WriteString(s.cond.String())

	return s.db.PrepareSQL(s.buf.String())
}

func (s *SQL) Update(args ...interface{}) (sql.Result, error) {
	if s.HasErrors() {
		return nil, Errors(s.errors)
	}

	if len(args) == 0 {
		args = append(s.vals, s.condArgs)
	}

	return s.db.Exec(s.updateSQL(), args...)
}

func (s *SQL) insertSQL() string {
	s.buf.Reset() // 清空之前的内容

	s.buf.WriteString("INSERT INTO ")
	s.buf.WriteString(s.tableName)

	s.buf.WriteByte('(')
	s.buf.WriteString(strings.Join(s.cols, ","))
	s.buf.WriteString(") VALUES(")
	placeholder := strings.Repeat("?,", len(s.cols))
	// 去掉上面的最后一个逗号
	s.buf.WriteString(placeholder[0 : len(placeholder)-1])
	s.buf.WriteByte(')')

	return s.db.PrepareSQL(s.buf.String())
}

func (s *SQL) Insert(args ...interface{}) (sql.Result, error) {
	if s.HasErrors() {
		return nil, Errors(s.errors)
	}

	if len(args) == 0 {
		args = s.vals
	}

	return s.db.Exec(s.insertSQL(), args...)
}
