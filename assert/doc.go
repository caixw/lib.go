// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// assert是对testing包的一些简单包装。方便在测试包里少写一点代码。
//
// 提供了两种操作方式：直接调用包函数；或是使用assertion对象。两
// 种方式完全等价，可以根据自己需要，选择一种。
//  func TestAssert(t *testing.T) {
//      var v interface{} = 5
//
//      // 直接调用包函数
//      assert.True(t, v == 5, "v的值[%v]不等于5", v)
//      assert.Equal(t, 5, v, "v的值[%v]不等于5", v)
//      assert.Nil(t, v)
//
//      // 以assertion对象方式使用
//      a := assert.New(t)
//      a.True(v==5, "v的值[%v]不等于5", v)
//      a.Equal(5, v, "v的值[%v]不等于5", v)
//      a.Nil(v)
//  }
package assert

// 当前库的版本号
const Version = "1.3.3.140728"
