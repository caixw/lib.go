// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"

	"github.com/caixw/lib.go/orm/core"
	"github.com/caixw/lib.go/orm/sqlbuilder"
)

// 事务对象
type Tx struct {
	engine *Engine
	tx     *sql.Tx
}

func (t *Tx) Name() string {
	return t.engine.Name()
}

func (t *Tx) GetStmts() *core.Stmts {
	return t.engine.GetStmts()
}

func (t *Tx) PrepareSQL(sql string) string {
	return t.engine.PrepareSQL(sql)
}

func (t *Tx) Dialect() core.Dialect {
	return t.engine.Dialect()
}

func (t *Tx) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(sql, args...)
}

func (t *Tx) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(sql, args...)
}
func (t *Tx) QueryRow(sql string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(sql, args...)
}

func (t *Tx) Prepare(sql string) (*sql.Stmt, error) {
	return t.tx.Prepare(sql)
}

// 关闭当前的db
func (t *Tx) close() {
	// 仅仅取消与engine的关联。
	t.engine = nil
}

// 提交事务
// 提交之后，整个Tx对象将不再有效。
func (t *Tx) Commit() (err error) {
	if err = t.tx.Commit(); err == nil {
		t.close()
	}
	return
}

// 回滚事务
func (t *Tx) Rollback() {
	t.tx.Rollback()
}

func (t *Tx) Update() *sqlbuilder.Update {
	return sqlbuilder.NewUpdate(t)
}

func (t *Tx) Delete() *sqlbuilder.Delete {
	return sqlbuilder.NewDelete(t)
}

func (t *Tx) Insert() *sqlbuilder.Insert {
	return sqlbuilder.NewInsert(t)
}

func (t *Tx) Select() *sqlbuilder.Select {
	return sqlbuilder.NewSelect(t)
}

// 查找缓存的sql.Stmt，在未找到的情况下，第二个参数返回false
func (t *Tx) Stmt(name string) (*sql.Stmt, bool) {
	stmt, found := t.engine.Stmt(name)
	if !found {
		return nil, false
	}

	return t.tx.Stmt(stmt), true
}
