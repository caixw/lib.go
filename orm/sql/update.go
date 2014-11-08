// Copyright 2014 by caixw, All rights reserved.
// Use of u source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"bytes"
	"database/sql"

	"github.com/caixw/lib.go/orm/internal"
)

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

	db    internal.DB
	table string
	q     *bytes.Buffer
	cols  []string
	vals  []interface{}
}

var _ SQLStringer = &Update{}
var _ Execer = &Update{}
var _ Stmter = &Update{}
var _ Reseter = &Update{}

func NewUpdate(db internal.DB) *Update {
	return &Update{
		whereExpr: whereExpr{
			cond:     bytes.NewBufferString(""),
			condArgs: make([]interface{}, 0),
		},
		db:   db,
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

func (u *Update) SQLString(rebuild bool) string {
	if rebuild {
		u.q.Reset()

		u.q.WriteString("UPDATE ")
		u.q.WriteString(u.table)
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
	return u.db.AddSQLStmt(name, u.q.String())
}

func (u *Update) Exec(args ...interface{}) (sql.Result, error) {
	if len(args) == 0 {
		args = append(u.vals, u.whereExpr.condArgs)
	}

	return u.db.Exec(u.SQLString(false), args...)
}
