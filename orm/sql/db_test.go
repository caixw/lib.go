// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/caixw/lib.go/orm/dialect"
	"github.com/caixw/lib.go/orm/internal"
	// 加载测试用例
	_ "github.com/caixw/lib.go/orm/dialect/test"
)

// 实现两个internal.DB接口，分别对应sql包的DB和Tx结构，
// 供SQL和model包使用

const tableName = "table."

// 用于测试的internal.DB接口实例
type db struct {
	prefix  string
	dialect dialect.Dialect
	db      *sql.DB
	stmts   map[string]*sql.Stmt
}

var _ internal.DB = &db{}

func newDB() (*db, error) {
	d, found := dialect.Get("mysql")
	if !found {
		return nil, errors.New("未找到与mysql相同的Dialect")
	}

	dbInst, err := sql.Open("mysql", "root:@/")
	if err != nil {
		return nil, err
	}

	return &db{
		db:      dbInst,
		dialect: d,
		prefix:  "prefix_",
		stmts:   map[string]*sql.Stmt{},
	}, nil
}

func (d *db) AddSQLStmt(name, sql string) (*sql.Stmt, error) {
	if _, found := d.stmts[name]; found {
		return nil, fmt.Errorf("名为[%v]的stmt已经存在", name)
	}

	stmt, err := d.db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	d.stmts[name] = stmt
	return stmt, nil
}

func (d *db) SetSQLStmt(name, sql string) (*sql.Stmt, error) {
	stmt, err := d.db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	d.stmts[name] = stmt
	return stmt, nil
}

func (d *db) GetStmt(name string) (stmt *sql.Stmt, found bool) {
	stmt, found = d.stmts[name]
	return
}

var replaceQuoteExpr = regexp.MustCompile(`("{1})([^\.\*," ]+)("{1})`)

// 替换语句中的双引号为指定的符号。
// 若sql的值中包含双引号也会被替换，所以所有的值只能是占位符。
func (d *db) ReplaceQuote(sql string) string {
	left, right := d.Dialect().Quote()
	return replaceQuoteExpr.ReplaceAllString(sql, left+"$2"+right)
}

// 替换表名的"table."虚前缀为e.prefix。
func (d *db) ReplacePrefix(table string) string {
	return strings.Replace(table, tableName, d.prefix, -1)
}

func (d *db) Dialect() dialect.Dialect {
	return d.dialect
}

func (d *db) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(sql, args...)
}

func (d *db) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(sql, args...)
}
func (d *db) QueryRow(sql string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(sql, args...)
}

func (d *db) Prepare(sql string) (*sql.Stmt, error) {
	return d.db.Prepare(sql)
}
