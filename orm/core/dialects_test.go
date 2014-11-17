// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/caixw/lib.go/assert"
)

// fakeDialect1
type fakeDialect1 struct {
}

var _ Dialect = &fakeDialect1{}

func (t *fakeDialect1) GetDBName(dataSource string) string {
	return ""
}

func (t *fakeDialect1) QuoteStr() (string, string) {
	return "[", "]"
}

func (t *fakeDialect1) CreateTable(db DB, m *Model) error {
	return nil
}

func (t *fakeDialect1) LimitSQL(limit, offset int) (string, []interface{}) {
	return "", nil
}

func (t *fakeDialect1) SupportLastInsertId() bool {
	return true
}

// fakeDialect2
type fakeDialect2 struct {
	num int
}

var _ Dialect = &fakeDialect2{}

func (t *fakeDialect2) GetDBName(dataSource string) string {
	return ""
}

func (t *fakeDialect2) QuoteStr() (string, string) {
	return "{", "}"
}

func (t *fakeDialect2) CreateTable(db DB, m *Model) error {
	return nil
}

func (t *fakeDialect2) LimitSQL(limit, offset int) (string, []interface{}) {
	return "", nil
}

func (t *fakeDialect2) SupportLastInsertId() bool {
	return true
}

func TestDialect(t *testing.T) {
	a := assert.New(t)

	clearDialects()
	a.Empty(dialects.items)

	err := RegisterDialect("fake1", &fakeDialect1{})
	a.NotError(err).
		True(IsRegistedDialect("fake1"))

	// 注册一个相同名称的
	err = RegisterDialect("fake1", &fakeDialect2{})
	a.Error(err)                    // 注册失败
	a.Equal(1, len(dialects.items)) // 数量还是1，注册没有成功

	// 再注册一个名称不相同的
	err = RegisterDialect("fake2", &fakeDialect2{})
	a.NotError(err)
	a.Equal(2, len(dialects.items))

	// 注册类型相同，但名称不同的实例
	err = RegisterDialect("fake3", &fakeDialect2{num: 2})
	a.Error(err)                    // 注册失败
	a.Equal(2, len(dialects.items)) // 数量还是2，注册没有成功

	// 清空
	clearDialects()
	a.Empty(dialects.items)
}
