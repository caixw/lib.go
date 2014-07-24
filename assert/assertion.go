// Copyright 2014 by caixw, All rights reserved.
// Use of a source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"testing"
)

// 对testing.T的一个封装，方便在一个测试函数中包含多个断言的情况下使用。
type assertion struct {
	t *testing.T
}

func New(t *testing.T) *assertion {
	return &assertion{t: t}
}

func (a *assertion) Assert(expr bool, msg string, args ...interface{}) {
	Assert(a.t, expr, msg, args...)
}

func (a *assertion) True(expr bool, msg string, args ...interface{}) {
	True(a.t, expr, msg, args...)
}

func (a *assertion) False(expr bool, msg string, args ...interface{}) {
	False(a.t, expr, msg, args...)
}
func (a *assertion) Nil(expr interface{}, msg string, args ...interface{}) {
	Nil(a.t, expr, msg, args...)
}

func (a *assertion) NotNil(expr interface{}, msg string, args ...interface{}) {
	NotNil(a.t, expr, msg, args...)
}

func (a *assertion) Equal(v1, v2 interface{}, msg string, args ...interface{}) {
	Equal(a.t, v1, v2, msg, args...)
}

func (a *assertion) NotEqual(v1, v2 interface{}, msg string, args ...interface{}) {
	NotEqual(a.t, v1, v2, msg, args...)
}

func (a *assertion) Empty(expr interface{}, msg string, args ...interface{}) {
	Empty(a.t, expr, msg, args...)
}
func (a *assertion) NotEmpty(expr interface{}, msg string, args ...interface{}) {
	NotEmpty(a.t, expr, msg, args...)
}
