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

// fakeDb1 测试有的driver.Driver接口实例，并未真实实现。
type fakeDb1 struct {
}

var _ driver.Driver = &fakeDb1{}

func (f *fakeDb1) Open(name string) (driver.Conn, error) {
	return nil, nil
}

// 测试用的dialect.Dialect接口实例。
type fakeDialect1 struct {
	base
}

var _ core.Dialect = &fakeDialect1{}

func (t *fakeDialect1) QuoteStr() (string, string) {
	return "[", "]"
}

func (t *fakeDialect1) LimitSQL(limit, offset int) (string, []interface{}) {
	return "", nil
}

// 注册测试需要用到的Dialect
func init() {
	sql.Register("fakedb1", &fakeDb1{})

	if !core.IsRegistedDialect("fakedb1") {
		err := core.RegisterDialect("fakedb1", &fakeDialect1{})
		if err != nil {
			panic(err)
		}
	}
}
