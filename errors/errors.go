// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// errors包的简单扩展，增加了错误代码和嵌套错误的功能
//  err := errors.Newf(5, nil, "错误代码%v", 5)
//  err2 := errors.New(6, err, "嵌套错误")
//  err3 := err2.GetPrevious()
package errors

import (
	"fmt"
)

// 当前库的版本
const Version = "0.1.2.140819"

type Errors struct {
	previous *Errors //错误链接中的前一个错误
	code     int
	msg      string
}

func (e *Errors) Error() string {
	return e.msg
}

// 返回错误代码。
func (e *Errors) GetCode() int {
	return e.code
}

// 返回错误链接中的前一个错误，若没有，则返回nil。
func (e *Errors) GetPrevious() *Errors {
	return e.previous
}

func New(code int, previous error, s string) *Errors {
	if previous == nil {
		return &Errors{code: code, previous: nil, msg: s}
	}

	err, ok := previous.(*Errors)
	if !ok {
		err = &Errors{msg: previous.Error()}
	}

	return &Errors{code: code, previous: err, msg: s}

}

// 带格式化的New
func Newf(code int, previous error, format string, args ...interface{}) *Errors {
	return New(code, previous, fmt.Sprintf(format, args...))
}
