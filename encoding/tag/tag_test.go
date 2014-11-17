// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package tag

import (
	"testing"

	"github.com/caixw/lib.go/assert"
)

var tag1 = "name,abc;name2,;;name3,n1,n2"
var tag2 = "name(abc);name2,;;name3(n1,n2)"

func TestReplace(t *testing.T) {
	tag := styleReplace.Replace(tag2)
	assert.Equal(t, tag, tag1)
}

func TestParse(t *testing.T) {
	a := assert.New(t)

	fn := func(t string) {
		m := Parse(tag1)
		a.Equal(3, len(m))

		a.Equal(m["name"][0], "abc")
		a.Equal(len(m["name"]), 1)

		a.Empty(m["name2"])

		a.Equal(len(m["name3"]), 2)
		a.Equal(m["name3"][0], "n1")
		a.Equal(m["name3"][1], "n2")
	}

	fn(tag1)
	fn(tag2)
}

func TestGet(t *testing.T) {
	a := assert.New(t)

	fn := func(t, name string, wont []string) {
		val, found := Get(t, name)
		a.True(found)
		a.Equal(val, wont, "[%v]与预期的值不符:结果值：[%v];预期值：[%v]", name, val, wont)
	}

	fn(tag1, "name", []string{"abc"})
	fn(tag1, "name2", []string{})
	fn(tag1, "name3", []string{"n1", "n2"})

	fn(tag2, "name", []string{"abc"})
	fn(tag2, "name2", []string{})
	fn(tag2, "name3", []string{"n1", "n2"})
}

func TestMustGet(t *testing.T) {
	a := assert.New(t)

	fn := func(tag1, name string, def []string, wont []string) {
		val := MustGet(tag1, name, def...)
		a.Equal(val, wont)
	}

	// name1不存在，测试默认值
	fn(tag1, "name1", []string{"def"}, []string{"def"})
	// abc不存在，测试默认值
	fn(tag1, "abc", []string{"defg", "abc"}, []string{"defg", "abc"})
	// name3存在，测试返回值
	fn(tag1, "name3", []string{"n3", "n4"}, []string{"n1", "n2"})

	// name1不存在，测试默认值
	fn(tag2, "name1", []string{"def"}, []string{"def"})
	// abc不存在，测试默认值
	fn(tag2, "abc", []string{"defg", "abc"}, []string{"defg", "abc"})
	// name3存在，测试返回值
	fn(tag2, "name3", []string{"n3", "n4"}, []string{"n1", "n2"})

}

func TestHas(t *testing.T) {
	a := assert.New(t)

	a.True(Has(tag1, "name"))
	a.True(Has(tag1, "name2"))
	a.True(Has(tag1, "name3"))
	a.False(Has(tag1, "name100"))

	a.True(Has(tag2, "name"))
	a.True(Has(tag2, "name2"))
	a.True(Has(tag2, "name3"))
	a.False(Has(tag2, "name100"))
}
