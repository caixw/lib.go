// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package internal

import (
	"database/sql"

	"github.com/caixw/lib.go/orm/dialect"
)

// 操作数据库的接口。
type DB interface {
	// 将sql预编译成sql.Stmt并缓存。若该名称已经存在，则将返回error
	AddSQLStmt(name, sql string) (*sql.Stmt, error)

	// 将sql预编译成sql.Stmt并缓存。若该名称已经存在，则将覆盖。
	SetSQLStmt(name, sql string) (*sql.Stmt, error)

	// 获取一个缓存的sql.Stmt，若不存在found返回false
	GetStmt(name string) (stmt *sql.Stmt, found bool)

	// 更换语句中的双引号为数据库中指定的字段引用符号
	ReplaceQuote(cols string) string

	// 替换表前缀为真实的前缀字符串
	ReplacePrefix(cols string) string

	// 返回dialect.Dialect接口
	Dialect() dialect.Dialect

	// 相当于sql.DB.Exec()
	Exec(sql string, args ...interface{}) (sql.Result, error)

	// 相当于sql.DB.Query()
	Query(sql string, args ...interface{}) (*sql.Rows, error)

	// 相当于sql.DB.QueryRow()
	QueryRow(sql string, args ...interface{}) *sql.Row

	// 相当于sql.DB.Prepare()
	Prepare(sql string) (*sql.Stmt, error)
}
