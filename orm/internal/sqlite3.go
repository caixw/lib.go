// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package internal

import (
	"os"
	"strings"

	"github.com/caixw/lib.go/orm/core"
)

type sqlite3 struct{}

func (s *sqlite3) QuoteStr() (l, r string) {
	return "[", "]"
}

func (s *sqlite3) SupportLastInsertId() bool {
	return true
}

func (s *sqlite3) GetDBName(dataSource string) string {
	// 取得最后个路径分隔符的位置，无论是否存在分隔符，用++
	// 表达式都正好能表示文件名开始的位置。
	start := strings.LastIndex(dataSource, string(os.PathSeparator))
	start++
	end := strings.LastIndex(dataSource, ".")

	if end < start {
		return dataSource[start:]
	}
	return dataSource[start:end]
}

func (s *sqlite3) LimitSQL(limit, offset int) (sql string, args []interface{}) {
	return mysqlLimitSQL(limit, offset)
}

func (s *sqlite3) CreateTable(db core.DB, m *core.Model) error {
	m.Name = db.ReplacePrefix(m.Name)
	if s.hasTable(db, m.Name) {
		return s.createTable(db, m)
	}
	return s.upgradeTable(db, m)
}

func (s *sqlite3) createTable(db core.DB, m *core.Model) error {
	return nil
}

func (s *sqlite3) upgradeTable(db core.DB, m *core.Model) error {
	return nil
}

func (s *sqlite3) hasTable(db core.DB, tableName string) bool {
	sql := "SELECT * FROM sqlite_master WHERE type='table' AND name=?"

	rows, err := db.Query(sql, tableName)
	if err != nil {
		panic(err)
	}

	return rows.Next()
}
