// Copyright 2014 by caixw, All rights reserved.
// Use of d source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"bytes"
	"database/sql"

	"github.com/caixw/lib.go/orm/internal"
)

// sql的delete语句
//  sql := newDelete(e)
//  sql.Table("table.user").
//      And("username=?", "admin").
//      Or("group=?", 1).
//      Exec()
type Delete struct {
	whereExpr

	db    internal.DB
	table string
	q     *bytes.Buffer
}

var _ SQLStringer = &Delete{}
var _ Execer = &Delete{}
var _ Stmter = &Delete{}
var _ Reseter = &Delete{}

func NewDelete(db internal.DB) *Delete {
	return &Delete{
		whereExpr: whereExpr{
			cond:     bytes.NewBufferString(""),
			condArgs: make([]interface{}, 0),
		},
		db: db,
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

func (d *Delete) SQLString(rebuild bool) string {
	if rebuild || d.q.Len() == 0 {
		d.q.Reset()

		d.q.WriteString("DELETE FROM ")
		d.q.WriteString(d.table)

		// where
		d.q.WriteString(d.condString(d.db))
	}

	return d.q.String()
}

func (d *Delete) Stmt(name string) (*sql.Stmt, error) {
	return d.db.AddSQLStmt(name, d.q.String())
}

func (d *Delete) Exec(args ...interface{}) (sql.Result, error) {
	if len(args) == 0 {
		args = d.condArgs
	}

	return d.db.Exec(d.SQLString(false), args...)
}
