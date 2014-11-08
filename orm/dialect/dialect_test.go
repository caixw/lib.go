// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package dialect

import (
	"reflect"
	"testing"

	"github.com/caixw/lib.go/assert"
)

// fakeDialect1
type fakeDialect1 struct {
}

var _ Dialect = &fakeDialect1{}

func (t *fakeDialect1) Quote() (string, string) {
	return "[", "]"
}

func (t *fakeDialect1) ToSqlType(typ reflect.Type, l1, l2 int) string {
	return ""
}

func (t *fakeDialect1) Limit(limit, offset int) (string, []interface{}) {
	return "", nil
}

// fakeDialect2
type fakeDialect2 struct {
	num int
}

var _ Dialect = &fakeDialect2{}

func (t *fakeDialect2) Quote() (string, string) {
	return "{", "}"
}

func (t *fakeDialect2) ToSqlType(typ reflect.Type, l1, l2 int) string {
	return ""
}

func (t *fakeDialect2) Limit(limit, offset int) (string, []interface{}) {
	return "", nil
}

func TestDialect(t *testing.T) {
	a := assert.New(t)

	clear()
	a.Empty(dialects.items)

	err := Register("fake1", &fakeDialect1{})
	a.NotError(err).
		True(IsRegisted("fake1"))

	// 注册一个相同名称的
	err = Register("fake1", &fakeDialect2{})
	a.Error(err)                    // 注册失败
	a.Equal(1, len(dialects.items)) // 数量还是1，注册没有成功

	// 再注册一个名称不相同的
	err = Register("fake2", &fakeDialect2{})
	a.NotError(err)
	a.Equal(2, len(dialects.items))

	// 注册类型相同，但名称不同的实例
	err = Register("fake3", &fakeDialect2{num: 2})
	a.Error(err)                    // 注册失败
	a.Equal(2, len(dialects.items)) // 数量还是2，注册没有成功

	// 清空
	clear()
	a.Empty(dialects.items)
}
