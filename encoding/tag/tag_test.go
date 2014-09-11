// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package tag

import (
	"testing"

	"github.com/caixw/lib.go/assert"
)

var tag = "name:abc;name2;;name3:n1,n2"

func TestParse(t *testing.T) {
	a := assert.New(t)

	m := Parse(tag)
	a.Equal(3, len(m))

	a.Equal(m["name"][0], "abc")
	a.Equal(len(m["name"]), 1)

	a.Empty(m["name2"])

	a.Equal(len(m["name3"]), 2)
	a.Equal(m["name3"][0], "n1")
	a.Equal(m["name3"][1], "n2")
}

func TestGet(t *testing.T) {
	a := assert.New(t)

	fn := func(tag, name string, wont []string) {
		val, found := Get(tag, name)
		a.True(found)
		a.Equal(val, wont)
	}

	fn(tag, "name", []string{"abc"})
	fn(tag, "name2", nil)
	fn(tag, "name3", []string{"n1", "n2"})
}

func TestMustGet(t *testing.T) {
	a := assert.New(t)

	fn := func(tag, name string, def []string, wont []string) {
		val := MustGet(tag, name, def...)
		a.Equal(val, wont)
	}

	// name1不存在，测试默认值
	fn(tag, "name1", []string{"def"}, []string{"def"})
	// abc不存在，测试默认值
	fn(tag, "abc", []string{"defg", "abc"}, []string{"defg", "abc"})
	// name3存在，测试返回值
	fn(tag, "name3", []string{"n3", "n4"}, []string{"n1", "n2"})
}

func TestHas(t *testing.T) {
	a := assert.New(t)
	a.True(Has(tag, "name"))
	a.True(Has(tag, "name2"))
	a.True(Has(tag, "name3"))
	a.False(Has(tag, "name100"))
}
