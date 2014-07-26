// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"testing"
)

// Assertion对testing.T进行了简单的封装。可以以调用对象的方式调用包中的
// 各个函数。方便在一个测试函数中包含多个断言的情况下使用。
type Assertion struct {
	t *testing.T
}

func New(t *testing.T) *Assertion {
	return &Assertion{t: t}
}

func (a *Assertion) True(expr bool, msg ...interface{}) {
	True(a.t, expr, msg...)
}

func (a *Assertion) False(expr bool, msg ...interface{}) {
	False(a.t, expr, msg...)
}
func (a *Assertion) Nil(expr interface{}, msg ...interface{}) {
	Nil(a.t, expr, msg...)
}

func (a *Assertion) NotNil(expr interface{}, msg ...interface{}) {
	NotNil(a.t, expr, msg...)
}

func (a *Assertion) Equal(v1, v2 interface{}, msg ...interface{}) {
	Equal(a.t, v1, v2, msg...)
}

func (a *Assertion) NotEqual(v1, v2 interface{}, msg ...interface{}) {
	NotEqual(a.t, v1, v2, msg...)
}

func (a *Assertion) Empty(expr interface{}, msg ...interface{}) {
	Empty(a.t, expr, msg...)
}
func (a *Assertion) NotEmpty(expr interface{}, msg ...interface{}) {
	NotEmpty(a.t, expr, msg...)
}

func (a *Assertion) Error(expr interface{}, msg ...interface{}) {
	Error(a.t, expr, msg...)
}

func (a *Assertion) NotError(expr interface{}, msg ...interface{}) {
	NotError(a.t, expr, msg...)
}

func (a *Assertion) FileExists(path string, msg ...interface{}) {
	FileExists(a.t, path, msg...)
}

func (a *Assertion) FileNotExists(path string, msg ...interface{}) {
	FileNotExists(a.t, path, msg...)
}
