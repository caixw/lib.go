// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"testing"

	"github.com/caixw/lib.go/assert"
	_ "github.com/caixw/lib.go/orm/dialect/test"
)

// 测试engines的一些常用操作：New,Get,Close,CloseAll
func TestEngines(t *testing.T) {
	a := assert.New(t)

	e, err := New("mysql", "root:@/", "test", "test_")
	a.NotError(err).NotNil(e)

	// 不存在的实例
	e, found := Get("test1test")
	a.False(found).Nil(e)

	// 获取注册的名为test的Engine实例
	e, found = Get("test")
	a.True(found).NotNil(e)

	// 关闭之后，是否能再次正常获取
	Close("test")
	e, found = Get("test")
	a.False(found).Nil(e)

	// 重新添加2个Engine

	e, err = New("mysql", "root:@/", "test", "test_")
	a.NotError(err).NotNil(e)

	e, err = New("fakedb1", "root:@/", "fakedb1", "fakedb1_")
	a.NotError(err).NotNil(e)

	e, found = Get("test")
	a.True(found).NotNil(e)

	e, found = Get("fakedb1")
	a.True(found).NotNil(e)

	// 关闭所有
	CloseAll()
	e, found = Get("test")
	a.False(found).Nil(e)
	e, found = Get("fakedb1")
	a.False(found).Nil(e)
}
