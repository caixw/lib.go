// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/caixw/lib.go/orm/core"
	"github.com/caixw/lib.go/orm/dialect"
)

// 实现两个internal.DB接口，分别对应sql包的DB和Tx结构，
// 供SQL和model包使用

type Engine struct {
	name   string // 数据库的名称
	prefix string // 表名前缀
	d      core.Dialect
	db     *sql.DB
	stmts  *core.Stmts
}

func newEngine(driverName, dataSourceName, prefix string) (*Engine, error) {
	d, found := dialect.Get(driverName)
	if !found {
		return nil, fmt.Errorf("未找到与driverName[%v]相同的Dialect", driverName)
	}

	dbInst, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	inst := &Engine{
		db:     dbInst,
		d:      d,
		prefix: prefix,
		name:   d.GetDBName(dataSourceName),
	}
	inst.stmts = core.NewStmts(inst)

	return inst, nil
}

// 对orm/core.DB.Name()的实现，返回当前操作的数据库名称。
func (e *Engine) Name() string {
	return e.name
}

// 对orm/core.DB.GetStmts()的实现，返回当前的sql.Stmt实例缓存容器。
func (e *Engine) GetStmts() *core.Stmts {
	return e.stmts
}

// 对orm/core.DB.PrepareSQL()的实现。替换语句的各种占位符。
func (e *Engine) PrepareSQL(sql string) string {
	// TODO 缓存replace
	l, r := e.Dialect().QuoteStr()
	replace := strings.NewReplacer("{", l, "}", r, "#", e.prefix)

	return replace.Replace(sql)
}

// 对orm/core.DB.Dialect()的实现。返回当前数据库对应的Dialect
func (e *Engine) Dialect() core.Dialect {
	return e.d
}

// 对orm/core.DB.Exec()的实现。执行一条非查询的SQL语句。
func (e *Engine) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return e.db.Exec(sql, args...)
}

// 对orm/core.DB.Query()的实现，执行一条查询语句。
func (e *Engine) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return e.db.Query(sql, args...)
}

// 对orm/core.DB.QueryRow()的实现。
// 执行一条查询语句，并返回第一条符合条件的记录。
func (e *Engine) QueryRow(sql string, args ...interface{}) *sql.Row {
	return e.db.QueryRow(sql, args...)
}

// 对orm/core.DB.Prepare()的实现。预处理SQL语句成sql.Stmt实例。
func (e *Engine) Prepare(sql string) (*sql.Stmt, error) {
	return e.db.Prepare(sql)
}

// 关闭当前的db，销毁所有的数据。不能再次使用。
func (e *Engine) close() {
	e.stmts.Close()
	e.db.Close()
}

// 开始一个新的事务
func (e *Engine) Begin() (*Tx, error) {
	tx, err := e.db.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{
		engine: e,
		tx:     tx,
	}, nil
}

// 查找缓存的sql.Stmt，在未找到的情况下，第二个参数返回false
func (e *Engine) Stmt(name string) (*sql.Stmt, bool) {
	return e.stmts.Get(name)
}

// 根据obj创建表
func (e *Engine) Create(obj interface{}) error {
	m, err := core.NewModel(obj)
	if err != nil {
		return err
	}
	return e.Dialect().CreateTable(e, m)
}

func (e *Engine) Update() *Update {
	return newUpdate(e)
}

func (e *Engine) Delete() *Delete {
	return newDelete(e)
}

func (e *Engine) Insert() *Insert {
	return newInsert(e)
}

func (e *Engine) Select() *Select {
	return newSelect(e)
}

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

func (t *Tx) Update() *Update {
	return newUpdate(t)
}

func (t *Tx) Delete() *Delete {
	return newDelete(t)
}

func (t *Tx) Insert() *Insert {
	return newInsert(t)
}

func (t *Tx) Select() *Select {
	return newSelect(t)
}

// 查找缓存的sql.Stmt，在未找到的情况下，第二个参数返回false
func (t *Tx) Stmt(name string) (*sql.Stmt, bool) {
	stmt, found := t.engine.Stmt(name)
	if !found {
		return nil, false
	}

	return t.tx.Stmt(stmt), true
}
