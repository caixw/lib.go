// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"database/sql"
)

// 通用但又没有统一标准的数据库功能接口。
//
// 有可能一个Dialect实例会被多个其它实例引用，
// 不应该在Dialect实例中保存状态值等内容。
type Dialect interface {
	// 对字段或是表名的引用字符
	QuoteStr() (left, right string)

	// 是否支持返回LastInsertId()特性
	SupportLastInsertId() bool

	// 从dataSourceName变量中获取数据库的名称
	GetDBName(dataSourceName string) string

	// 生成limit n offset m语句
	// 返回的是对应数据库的limit语句以及语句中占位符对应的值
	LimitSQL(limit, offset int) (sql string, args []interface{})

	// 根据一个Model创建或是更新表。
	// 表的创建虽然语法上大致上相同，但细节部分却又不一样，
	// 干脆整个过程完全交给Dialect去完成。
	CreateTable(db DB, m *Model) error
}

// 操作数据库的接口，用于统一普通数据库操作和事务操作。
type DB interface {
	// 当前操作数据库的名称
	Name() string

	// 获取一个缓存的sql.Stmt，若不存在found返回false
	GetStmts() *Stmts

	// 更换语句中的双引号为数据库中指定的字段引用符号
	ReplaceQuote(cols string) string

	// 替换表前缀为真实的前缀字符串
	ReplacePrefix(cols string) string

	// 返回Dialect接口
	Dialect() Dialect

	// 相当于sql.DB.Exec()
	Exec(sql string, args ...interface{}) (sql.Result, error)

	// 相当于sql.DB.Query()
	Query(sql string, args ...interface{}) (*sql.Rows, error)

	// 相当于sql.DB.QueryRow()
	QueryRow(sql string, args ...interface{}) *sql.Row

	// 相当于sql.DB.Prepare()
	Prepare(sql string) (*sql.Stmt, error)
}

type conType int

// 预定的约束类型，方便Model中使用。
const (
	none conType = iota
	index
	unique
	fk
	check
)

func (t conType) String() string {
	switch t {
	case none:
		return "<none>"
	case index:
		return "KEY INDEX"
	case unique:
		return "UNIQUE INDEX"
	case fk:
		return "FOREIGN KEY"
	case check:
		return "CHECK"
	default:
		return "<unknown>"
	}
}
