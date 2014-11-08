// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"database/sql"
)

// 所有sql语句的接口
type SQLStringer interface {
	// 返回当前实例的SQL语句。rebuild参数指定是否重新产生SQL语句。
	SQLString(rebuild bool) string
}

type Execer interface {
	Exec(args ...interface{}) (sql.Result, error)
}

// 重置语句接口
type Reseter interface {
	Reset()
}

// 缓存成stmt接口
type Stmter interface {
	// 将当前语句预编译成stmt并以name缓存
	Stmt(name string) (*sql.Stmt, error)
}

type Fetch interface {
	Query(args ...interface{}) (*sql.Rows, error)

	QueryRow(args ...interface{}) *sql.Row

	// 导出到v中
	// 若v是数组，则导出多条语句，若v是对象，则导出第一个对象
	Fetch(v interface{}, args ...interface{}) error

	Fetch2Map(args ...interface{}) (map[string]interface{}, error)

	Fetch2Maps(args ...interface{}) ([]map[string]interface{}, error)

	// 导出一列数据到v中；
	FetchColumns(col string, args ...interface{}) ([]interface{}, error)

	FetchColumn(col string, args ...interface{}) (interface{}, error)
}
