// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package test

import (
	"database/sql"
	"database/sql/driver"
	"reflect"

	"github.com/caixw/lib.go/orm/dialect"
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
	num int
}

var _ dialect.Dialect = &fakeDialect2{}

func (t *fakeDialect2) Quote() (string, string) {
	return "{", "}"
}

func (t *fakeDialect2) ToSqlType(typ reflect.Type, l1, l2 int) string {
	return ""
}

func (t *fakeDialect2) Limit(limit, offset int) (string, []interface{}) {
	return "", nil
}

// 注册测试需要用到的Dialect
func init() {
	sql.Register("fakedb2", &fakeDb2{})

	if !dialect.IsRegisted("fakedb2") {
		err := dialect.Register("fakedb2", &fakeDialect2{})
		if err != nil {
			panic(err)
		}
	}
}
