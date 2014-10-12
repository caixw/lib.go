// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// ini操作包
//
// 提供了最基本的Reader和Writer两个操作类。
//
// ini还提供了Unmarshal与Marshal两个函数，功能类似于xml
// 中的同名函数。也可以通过struct tag指定名称，但格式上
// 稍微有一点不一样，其格式如下：
//  type User struct {
//      ID int `ini:"name:id"`
//  }
// 通过struct tag还可能指定一个转换函数，指定数据之间是
// 是如何相互转换的。比如上图，如何从一个字符串的id转换
// 成int类型的User.ID：
//  type User struct {
//      ID int `ini:"name:id;get:GetId;set:SetId"`
//  }
//
//  // 声明一个SetId的同名函数，之后每次遇到ID变量，都会调用
//  // 此函数转换到User.ID中。
//  func (u *User)SetId(str string, v reflect.Value) error {
//      // do something...
//      return nil
//  }
//
// get与set的函数原型分别为：
//  // set
//  func (val string, v reflect.Value) error
//
//  // get
//  func (v reflect.Value) (string, error)
package ini

const Version = "0.1.6.141009"
