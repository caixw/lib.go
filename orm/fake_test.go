// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 声明了两个core.Dialect接口实例及与之相对应的sql.DB接口实例。
// 方便测试函数调用。但是除了QuoteStr()函数以外，并未实现其它
// 接口。两个的driverName分别为：
// fakedb1,fakedb2

package orm

import (
	"database/sql"
	"database/sql/driver"

	"github.com/caixw/lib.go/orm/core"
)

type base struct {
}

func (t base) GetDBName(dataSource string) string {
	return ""
}

func (t *base) CreateTable(db core.DB, m *core.Model) error {
	return nil
}

func (m *base) SupportLastInsertId() bool {
	return true
}

func (m *base) LimitSQL(limit, offset int) (string, []interface{}) {
	return "", nil
}

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
}

var _ core.Dialect = &fakeDialect2{}

func (t *fakeDialect2) QuoteStr() (string, string) {
	return "{", "}"
}

// 注册测试需要用到的Dialect
func init() {
	// 注册fakedb1及相当的dialect
	sql.Register("fakedb1", &fakeDb1{})

	if !core.IsRegistedDialect("fakedb1") {
		err := core.RegisterDialect("fakedb1", &fakeDialect1{})
		if err != nil {
			panic(err)
		}
	}

	// 注册fakedb2及相关的dialect
	sql.Register("fakedb2", &fakeDb2{})

	if !core.IsRegistedDialect("fakedb2") {
		err := core.RegisterDialect("fakedb2", &fakeDialect2{})
		if err != nil {
			panic(err)
		}
	}
}
