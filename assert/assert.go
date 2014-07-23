// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// assert是对testing包的一些简单包装。方便在测试包里少写一点代码。
//  func TestAssert(t *testing.T) {
//      assert.True(t, v == 5, "v的值[%v]不等于5", v)
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
	if !expr {
		if Detail {
			pc, file, line, ok := runtime.Caller(callerDepth)
			if ok {
				base := path.Base(file)
				fn := funcName(pc)
				msg = msg + "@" + fn + "(" + base + "):" + strconv.Itoa(line)
			}
		}
		t.Errorf(msg, args...)
	}
}

// 在expr不成立(false)时，调用调用t.Errorf(msg,args...)输出内容。
func True(t *testing.T, expr bool, msg string, args ...interface{}) {
	Assert(2, t, expr, msg, args...)
}

// 与Assert()函数相反，在expr为true时，输出信息
func False(t *testing.T, expr bool, msg string, args ...interface{}) {
	Assert(2, t, !expr, msg, args...)
}

// 判断是否为nil，不为nil时，输出信息
func Nil(t *testing.T, expr interface{}, msg string, args ...interface{}) {
	Assert(2, t, expr == nil, msg, args...)
}

func NotNil(t *testing.T, expr interface{}, msg string, args ...interface{}) {
	Assert(2, t, expr != nil, msg, args...)
}

// 判断两个值是否相等，不相等时输出msg信息
func Equal(t *testing.T, act, val interface{}, msg string, args ...interface{}) {
	Assert(2, t, reflect.DeepEqual(act, val), msg, args...)
}

// 判断两个值是否不同。
func NotEqual(t *testing.T, act, val interface{}, msg string, args ...interface{}) {
	Assert(2, t, !reflect.DeepEqual(act, val), msg, args...)
}
