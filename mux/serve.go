// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mux

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// 所有支持的方法
var methods = []string{
	"GET",
	"POST",
	"PUT",
	"DELETE",
	"HEAD",
}

type muxEntry struct {
	patternC *regexp.Regexp
	handler  http.Handler
}

// 对http.ServeMux的简单扩展，可以实现对正则表达式和method的匹配。
type ServeMux struct {
	sync.Mutex

	entries map[string]map[string]*muxEntry
}

var _ http.Handler = &ServeMux{}

// 初始化ServeMux
func NewServeMux() *ServeMux {
	rs := make(map[string]map[string]*muxEntry)
	for _, m := range methods {
		rs[m] = make(map[string]*muxEntry, 0)
	}

	return &ServeMux{entries: rs}
}

// 添加路由。
//
// pattern参数可以是字符串，也可以是正则表达式。
// 所有的错误都直接panic，而不是像其它函数一样返回error。
func (srv *ServeMux) Handle(pattern string, handler http.Handler, methods ...string) {
	if len(methods) <= 0 {
		panic("请至少提供一个methods参数")
	}

	srv.Lock()
	defer srv.Unlock()

	for _, m := range methods {
		m = strings.ToUpper(m)
		_, found := srv.entries[m]
		if !found {
			panic("不存在的方法:" + m)
		}

		if _, found := srv.entries[m][pattern]; found {
			msg := fmt.Sprintf("该规则[%v]的路由已经存在", pattern)
			panic(msg)
		}

		srv.entries[m][pattern] = &muxEntry{
			patternC: regexp.MustCompile(pattern),
			handler:  handler,
		}
	}
}

// 添加路由
func (srv *ServeMux) HandleFunc(pattern string, handle http.HandlerFunc, methods ...string) {
	srv.Handle(pattern, http.Handler(handle), methods...)
}

func (srv *ServeMux) Get(pattern string, handler http.Handler) {
	srv.Handle(pattern, handler, "GET")
}

func (srv *ServeMux) GetFunc(pattern string, handle http.HandlerFunc) {
	srv.Get(pattern, http.Handler(handle))
}

func (srv *ServeMux) Post(pattern string, handler http.Handler) {
	srv.Handle(pattern, handler, "POST")
}

func (srv *ServeMux) PostFunc(pattern string, handle http.HandlerFunc) {
	srv.Post(pattern, http.Handler(handle))
}

func (srv *ServeMux) Delete(pattern string, handler http.Handler) {
	srv.Handle(pattern, handler, "DELETE")
}

func (srv *ServeMux) DeleteFunc(pattern string, handle http.Handler) {
	srv.Delete(pattern, http.Handler(handle))
}

func (srv *ServeMux) Put(pattern string, handler http.Handler) {
	srv.Handle(pattern, handler, "PUT")
}

func (srv *ServeMux) PutFunc(pattern string, handle http.HandlerFunc) {
	srv.Put(pattern, http.Handler(handle))
}

func (srv *ServeMux) Any(pattern string, handler http.Handler) {
	srv.Handle(pattern, handler, methods...)
}

func (srv *ServeMux) AnyFunc(pattern string, handle http.HandlerFunc) {
	srv.Any(pattern, http.Handler(handle))
}

// implement http.Handler
func (srv *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.serveHTTP(w, r)
}

// http.Handler::ServeHTTP()的具体实现，若有找到匹配了则执行完，并返回true，否则返回false
func (srv *ServeMux) serveHTTP(w http.ResponseWriter, r *http.Request) bool {
	entries, found := srv.entries[r.Method]
	if !found {
		return false
	}

	for _, entry := range entries {
		ok := entry.patternC.MatchString(r.URL.Path)
		if !ok {
			continue
		}

		// todo params
		entry.handler.ServeHTTP(w, r)

		freeContext(r) // 释放context的内容
		return true
	}
	return false
}
