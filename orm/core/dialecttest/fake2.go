// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package test

import (
	"database/sql"
	"database/sql/driver"

	"github.com/caixw/lib.go/orm/core"
	_ "github.com/go-sql-driver/mysql"
)

// 测试用driver.Driver接口实例。
type fakeDb2 struct {
}

var _ driver.Driver = &fakeDb2{}

func (f *fakeDb2) Open(name string) (driver.Conn, error) {
	return nil, nil
}

// 测试用dialect.Dialect接口实例
type fakeDialect2 struct {
	base
	num int
}

var _ core.Dialect = &fakeDialect2{}

func (t *fakeDialect2) QuoteStr() (string, string) {
	return "{", "}"
}

func (t *fakeDialect2) LimitSQL(limit, offset int) (string, []interface{}) {
	return "", nil
}

// 注册测试需要用到的Dialect
func init() {
	sql.Register("fakedb2", &fakeDb2{})

	if !core.IsRegistedDialect("fakedb2") {
		err := core.RegisterDialect("fakedb2", &fakeDialect2{})
		if err != nil {
			panic(err)
		}
	}
}
