// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/caixw/lib.go/orm/dialect"
	"github.com/caixw/lib.go/orm/internal"
)

// 实现两个internal.DB接口，分别对应sql包的DB和Tx结构，
// 供SQL和model包使用

const tableName = "table."

type db struct {
	prefix  string
	dialect dialect.Dialect
	db      *sql.DB
	stmts   *Stmts
}

var _ internal.DB = &db{}

func newDB(driverName, dataSourceName, prefix string) (*db, error) {
	d, found := dialect.Get(driverName)
	if !found {
		return nil, fmt.Errorf("未找到与driverName[%v]相同的Dialect", driverName)
	}

	dbInst, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	inst := &db{
		db:      dbInst,
		dialect: d,
		prefix:  prefix,
	}
	inst.stmts = newStmts(inst)

	return inst, nil
}

// implement internal.DB.AddSQLStmt()
func (d *db) AddSQLStmt(name, sql string) (*sql.Stmt, error) {
	return d.stmts.AddSql(name, sql)
}

// implement internal.DB.SetSQLStmt()
func (d *db) SetSQLStmt(name, sql string) (*sql.Stmt, error) {
	return d.stmts.SetSql(name, sql)
}

// implement internal.DB.GetSQLStmt()
func (d *db) GetStmt(name string) (stmt *sql.Stmt, found bool) {
	return d.stmts.Get(name)
}

var replaceQuoteExpr = regexp.MustCompile(`("{1})([^\.\*," ]+)("{1})`)

// implement internal.DB.ReplaceQuote()
// 替换语句中的双引号为指定的符号。
// 若sql的值中包含双引号也会被替换，所以所有的值只能是占位符。
func (d *db) ReplaceQuote(sql string) string {
	left, right := d.Dialect().Quote()
	return replaceQuoteExpr.ReplaceAllString(sql, left+"$2"+right)
}

// implement internal.DB.ReplacePrefix()
// 替换表名的"table."虚前缀为e.prefix。
func (d *db) ReplacePrefix(table string) string {
	return strings.Replace(table, tableName, d.prefix, -1)
}

// implement internal.DB.Dialect()
func (d *db) Dialect() dialect.Dialect {
	return d.dialect
}

// implement internal.DB.Exec()
func (d *db) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(sql, args...)
}

// implement internal.DB.Query()
func (d *db) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(sql, args...)
}

// implement internal.DB.QueryRow()
func (d *db) QueryRow(sql string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(sql, args...)
}

// implement internal.DB.Prepare()
func (d *db) Prepare(sql string) (*sql.Stmt, error) {
	return d.db.Prepare(sql)
}

// 关闭当前的db
func (d *db) close() {
	d.stmts.free()
	d.db.Close()
}

type tx struct {
	db *db
	tx *sql.Tx
}

var _ internal.DB = &tx{}

func (t *tx) AddSQLStmt(name, sql string) (*sql.Stmt, error) {
	return t.db.AddSQLStmt(name, sql)
}

func (t *tx) SetSQLStmt(name, sql string) (*sql.Stmt, error) {
	return t.db.SetSQLStmt(name, sql)
}

func (t *tx) GetStmt(name string) (*sql.Stmt, bool) {
	return t.db.GetStmt(name)
}

// 替换语句中的双引号为指定的符号。
// 若sql的值中包含双引号也会被替换，所以所有的值只能是占位符。
func (t *tx) ReplaceQuote(sql string) string {
	return t.db.ReplaceQuote(sql)
}

// 替换表名的"table."虚前缀为e.prefix。
func (t *tx) ReplacePrefix(table string) string {
	return t.db.ReplacePrefix(table)
}

func (t *tx) Dialect() dialect.Dialect {
	return t.db.Dialect()
}

func (t *tx) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(sql, args...)
}

func (t *tx) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.Query(sql, args...)
}
func (t *tx) QueryRow(sql string, args ...interface{}) *sql.Row {
	return t.tx.QueryRow(sql, args...)
}

func (t *tx) Prepare(sql string) (*sql.Stmt, error) {
	return t.tx.Prepare(sql)
}

// 关闭当前的db
func (t *tx) close() {
	// 仅退出tx，但相关联的db还是继续运行
	//t.db.stmts.free()
	//t.db.db.Close()
}
