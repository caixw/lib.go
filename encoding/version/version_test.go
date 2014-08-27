// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package version

import (
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestParse(t *testing.T) {
	a := assert.New(t)

	type parseType struct {
		version string
		val     []string
	}

	vals := []parseType{
		{"1.0.0", []string{"1", "0", "0"}},
		{"1...0.0", []string{"1", "0", "0"}},
		{"0.1.build1004", []string{"0", "1", "build", "1004"}},
		{"0.1+build1004.1", []string{"0", "1", "build", "1004", "1"}},
		{"0.1-1.0", []string{"0", "1", "1", "0"}},
		// {"1.0.1构建日期2014", []string{"1", "0", "1", "构建日期", "2014"}},
	}

	for _, v := range vals {
		verStr, err := Parse(v.version)
		a.NotError(err)
		a.Equal(verStr, v.val)
	}

}

func TestCompare(t *testing.T) {
	a := assert.New(t)

	const (
		gt = iota
		lt
		eq
	)

	type cmpType struct {
		v1, v2 string
		op     int
	}

	vals := []cmpType{
		{"0.1.0", "0.1.0", eq},
		{"1...0.0", "1.0.0", eq},
		{"1.0-alpha", "1.0-", lt},
		{"1.0+build1", "1.0build1.1", lt},
		{"1.0.build1.1", "1.0build", gt},
	}

	for k, v := range vals {
		switch v.op {
		case gt:
			a.True(Compare(v.v1, v.v2) > 0, "在%v个元素[%v]出错", k, v)
		case lt:
			a.True(Compare(v.v1, v.v2) < 0, "在%v个元素[%v]出错", k, v)
		case eq:
			a.Equal(Compare(v.v1, v.v2), 0, "在%v个元素[%v]出错", k, v)
		}
	}
}
