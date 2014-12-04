// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package dialects

import (
	"github.com/caixw/lib.go/orm/core"
	"strings"
)

type pq struct{}

// implement core.Dialect.QuoteStr()
func (p *pq) QuoteStr() (l, r string) {
	return "[", "]"
}

// implement core.Dialect.SupportLastInsertId()
func (p *pq) SupportLastInsertId() bool {
	return true
}

// implement core.Dialect.GetDBName()
func (p *pq) GetDBName(dataSource string) string {
	// dataSource样式：user=user dbname=db password=
	index := strings.Index(dataSource, "dbname=")
	dataSource = dataSource[index+1:]
	index = strings.Index(dataSource, " ") // BUG(caixw) 判断\t等其它字符
	return dataSource[:index]
}

// implement core.Dialect.LimitSQL()
func (p *pq) LimitSQL(limit, offset int) (sql string, args []interface{}) {
	return mysqlLimitSQL(limit, offset)
}

// implement core.Dialect.CreateTable()
func (s *pq) CreateTable(db core.DB, m *core.Model) error {
	// TODO

	return nil
}
