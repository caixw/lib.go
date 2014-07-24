// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// 是否输出详细的信息，如断言发生的文件，行数等。
const showCallerInfo = true

// 获取某个pc寄存器中的函数名，去掉函数名之前的路径信息。
func funcName(pc uintptr) string {
	if pc == 0 {
		return "<unknown func>"
	}
	name := runtime.FuncForPC(pc).Name()
	arr := strings.Split(name, "/")
	return arr[len(arr)-1]
}

// 定位到测试包中的信息，并输出信息。
// 若测试包中的函数是嵌套调用的，则有可能显示不正确。
func getCallerInfo() string {
	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			return ""
		}

		basename := path.Base(file)

		// 定位以_test.go结尾的文件，认定为起始调用的测试包。
		l := len(basename)
		if l < 8 || (basename[l-8:l] != "_test.go") {
			continue
		}

		return " @ " + funcName(pc) + "(" + basename + "):" + strconv.Itoa(line)
	}

	return ""
}

// 当expr条件不成立时，输出错误信息。
// expr 返回结果值为bool类型的表达式；msg, args输出的错误信息。
func Assert(t *testing.T, expr bool, msg string, args ...interface{}) {
	if expr {
		return
	}

	if showCallerInfo { // 输出一些调用堆栈的信息
		msg += getCallerInfo()
	}
	t.Errorf(msg, args...)
}

// 断言表达式expr为true，否则输出错误信息
func True(t *testing.T, expr bool, msg string, args ...interface{}) {
	Assert(t, expr, msg, args...)
}

// 断言表达式expr为false，否则输出错误信息
func False(t *testing.T, expr bool, msg string, args ...interface{}) {
	Assert(t, !expr, msg, args...)
}

func isNil(expr interface{}) bool {
	if nil == expr {
		return true
	}

	v := reflect.ValueOf(expr)
	k := v.Kind()

	if (k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.Interface ||
		k == reflect.Map ||
		k == reflect.Ptr ||
		k == reflect.Slice) && v.IsNil() {
		return true
	}

	return false
}

// 断言表达式expr为nil，否则输出错误信息
func Nil(t *testing.T, expr interface{}, msg string, args ...interface{}) {
	Assert(t, isNil(expr), msg, args...)
}

// 断言表达式expr为非nil值，否则输出错误信息
func NotNil(t *testing.T, expr interface{}, msg string, args ...interface{}) {
	Assert(t, !isNil(expr), msg, args...)
}

// 断言act与val两个值相等，否则输出错误信息
func Equal(t *testing.T, v1, v2 interface{}, msg string, args ...interface{}) {
	Assert(t, reflect.DeepEqual(v1, v2), msg, args...)
}

// 断言act与val两个值不相等，否则输出错误信息
func NotEqual(t *testing.T, v1, v2 interface{}, msg string, args ...interface{}) {
	Assert(t, !reflect.DeepEqual(v1, v2), msg, args...)
}

func isEmpty(expr interface{}) bool {
	if expr == nil {
		return true
	}

	switch v := expr.(type) {
	case bool:
		return false == v
	case int:
		return 0 == v
	case int8:
		return 0 == v
	case int16:
		return 0 == v
	case int32:
		return 0 == v
	case int64:
		return 0 == v
	case uint:
		return 0 == v
	case uint8:
		return 0 == v
	case uint16:
		return 0 == v
	case uint32:
		return 0 == v
	case uint64:
		return 0 == v
	case string:
		return "" == v
	}

	v := reflect.ValueOf(expr)
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Chan:
		return 0 == v.Len()
	case reflect.Ptr:
		return false
	}

	return false
}

// 断言expr的值为空(nil,"",0,false)，否则输出错误信息
func Empty(t *testing.T, expr interface{}, msg string, args ...interface{}) {
	Assert(t, isEmpty(expr), msg, args...)
}

// 断言expr的值为非空(除nil,"",0,false之外)，否则输出错误信息
func NotEmpty(t *testing.T, expr interface{}, msg string, args ...interface{}) {
	Assert(t, !isEmpty(expr), msg, args...)
}
