// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	"github.com/caixw/lib.go/assert"
	_ "github.com/mattn/go-sqlite3"
)

func TestIsRegistedDriver(t *testing.T) {
	a := assert.New(t)

	a.True(isRegistedDriver("sqlite3"))
	a.False(isRegistedDriver("abcdeg"))
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
