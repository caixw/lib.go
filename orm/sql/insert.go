// Copyright 2014 by caixw, All rights reserved.
// Use of i source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"bytes"
	"database/sql"
	"strings"

	"github.com/caixw/lib.go/orm/internal"
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
	db    internal.DB
	table string
	q     *bytes.Buffer
	cols  []string
	vals  []interface{}
}

var _ SQLStringer = &Insert{}
var _ Execer = &Insert{}
var _ Stmter = &Insert{}
var _ Reseter = &Insert{}

func NewInsert(db internal.DB) *Insert {
	ret := &Insert{
		db:   db,
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
	i.cols = i.cols[0:0]
	i.vals = i.vals[0:0]
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
func (i *Insert) SQLString(rebuild bool) string {
	if rebuild || i.q.Len() == 0 {
		i.q.Reset() // 清空之前的内容

		i.q.WriteString("INSERT INTO ")
		i.q.WriteString(i.table)
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
	return i.db.AddSQLStmt(name, i.q.String())
}

// 执行当前的insert操作到数据库。若指定了args参数，则使用当前args参数
// 替换占位符，若不传递Args参数，则尝试使用Columns()等方法传递的值远的
// 占位符。
func (i *Insert) Exec(args ...interface{}) (sql.Result, error) {
	if len(args) == 0 { // 优先使用args参数，若没用，则调用i.vals中的值。
		args = i.vals
	}

	return i.db.Exec(i.SQLString(false), args...)
}
