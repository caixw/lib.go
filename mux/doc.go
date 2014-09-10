// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 这是对http.ServeMux及相关的一个简单扩展，使之可以支持
// 正则表达式的路由和对method的选择以及一个简单的多域名
// 路由器。
//
// 默认的路由操作：
//  // 普通的Get方法路由
//  mux.Get("/member/login", &memberLoginHandler{})
//  // 正则路由。
//  mux.Get("/member/(<?P<action>[a-zA-Z]+)", &memberAction{})
//  mux.ListenAndServe(":88")
//
// 多域名操作：
//  // s1只针对member.example.com域名才会路由
//  s1 := mux.NewServeMux()
//  s1.GetFunc("/member/login", loginHandle)
//  mux.Host("member.example.com", s1)
//
//  // s2针对除member.example.com域名以外的所有域名
//  s2 := mux.NewServeMux()
//  s2.GetFunc("/api/", apiHandle)
//  mux.Host("*.example.com", s2)
//  mux.ListenAndServe(":88")
package mux

// 当前包的版本
const Version = "0.1.0.140829"
