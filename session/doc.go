// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// session的操作包。仅支持go1.3+
//
// session的存储方式多种多样，大家可以通过实现Store接口，之后
// 将该Store的实例传递给session.New()的第一个参数，就可以实现
// 自定义的session存储。
//
// 不能将一个Store指针同时传递给多个session.New()，它们必须是
// 独立的实例，否则将会发生串号等错误：
//  // 错误用法
//  mem := memory.New()
//  inst1 := session.New(mem, ...)
//  inst2 := session.New(mem, ...)
//  // inst1,inst2使用同一个Store，将会发生串号现象。
//
//  // 正确用法
//  mem1 := memory.New()
//  mem2 := memory.New()
//  inst1 := session.New(mem1, ...)
//  inst2 := session.New(mem2, ...)
//  // inst1,inst2将会是2个完全独立的存储系统，互不干扰。
package session

// 当前库的版本号
const Version = "0.1.0.140816"
