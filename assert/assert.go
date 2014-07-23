// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// assert是对testing包的一些简单包装。方便在测试包里少写一点代码。
//  func TestAssert(t *testing.T) {
//      var v interface{} = 5
//      assert.True(t, v == 5, "v的值[%v]不等于5", v)
//      assert.Equal(t, 5, v, "v的值[%v]不等于5", v)
//      assert.Nil(t, v, "v的类型不为nil")
//  }
package assert

import (
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// 当前库的版本号
const Version = "0.1.0.140723"

// 是否输出详细的信息，如断言发生的文件，行数等。
var Detail = true

// 获取某个pc寄存器中的函数名，去掉函数名之前的路径信息。
func funcName(pc uintptr) string {
	if pc == 0 {
		return "<none>"
	}
	name := runtime.FuncForPC(pc).Name()
	arr := strings.Split(name, "/")
	return arr[len(arr)-1]
}

// 当expr条件不成立时，输出错误信息。
// callerDepth 调用者的深度，方便显示是在哪段测试代码中出错的，
// Assert()为0层，调用Assert()的为第1层，以此类推，当测试函数
// 直接调用Assert()时，此值为1；
// expr 返回结果值为bool类型的表达式；msg, args输出的错误信息。
func Assert(callerDepth int, t *testing.T, expr bool, msg string, args ...interface{}) {
	if expr {
		return
	}

	if Detail { // 输出一些调用堆栈的信息
		pc, file, line, ok := runtime.Caller(callerDepth)
		if ok {
			base := path.Base(file)
			fn := funcName(pc)
			msg = msg + "@" + fn + "(" + base + "):" + strconv.Itoa(line)
		}
	}
	t.Errorf(msg, args...)
}

// 断言表达式expr为true，否则输出错误信息
func True(t *testing.T, expr bool, msg string, args ...interface{}) {
	Assert(2, t, expr, msg, args...)
}

// 断言表达式expr为false，否则输出错误信息
func False(t *testing.T, expr bool, msg string, args ...interface{}) {
	Assert(2, t, !expr, msg, args...)
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
	Assert(2, t, isNil(expr), msg, args...)
}

// 断言表达式expr为非nil值，否则输出错误信息
func NotNil(t *testing.T, expr interface{}, msg string, args ...interface{}) {
	Assert(2, t, !isNil(expr), msg, args...)
}

// 断言act与val两个值相等，否则输出错误信息
func Equal(t *testing.T, act, val interface{}, msg string, args ...interface{}) {
	Assert(2, t, reflect.DeepEqual(act, val), msg, args...)
}

// 断言act与val两个值不相等，否则输出错误信息
func NotEqual(t *testing.T, act, val interface{}, msg string, args ...interface{}) {
	Assert(2, t, !reflect.DeepEqual(act, val), msg, args...)
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
	case reflect.Ptr: // TODO(caixw) 指针有为0的情况？
		return false
	}

	return false
}

// 断言expr的值为空(nil,"",0,false)，否则输出错误信息
func Empty(t *testing.T, expr interface{}, msg string, args ...interface{}) {
	Assert(2, t, isEmpty(expr), msg, args...)
}

// 断言expr的值为非空(除nil,"",0,false之外)，否则输出错误信息
func NotEmpty(t *testing.T, expr interface{}, msg string, args ...interface{}) {
	Assert(2, t, !isEmpty(expr), msg, args...)
}
