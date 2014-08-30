// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mux

import (
	"net/http"
)

var defaultHostMux = &HostMux{
	entries: make(map[string]*hostEntry),
}

func Host(pattern string, handler http.Handler) {
	defaultHostMux.Handle(pattern, handler)
}

func HostFunc(pattern string, handle http.HandlerFunc) {
	defaultHostMux.HandleFunc(pattern, handle)
}

var defaultServeMux = NewServeMux()

func Handle(pattern string, handler http.Handler) {
	defaultServeMux.Handle(pattern, handler)
}

func HandleFunc(pattern string, handle http.HandlerFunc) {
	defaultServeMux.HandleFunc(pattern, handle)
}

func Get(pattern string, handler http.Handler) {
	defaultServeMux.Get(pattern, handler)
}

func GetFunc(pattern string, handle http.HandlerFunc) {
	defaultServeMux.GetFunc(pattern, handle)
}

func Post(pattern string, handler http.Handler) {
	defaultServeMux.Post(pattern, handler)
}

func PostFunc(pattern string, handle http.HandlerFunc) {
	defaultServeMux.PostFunc(pattern, handle)
}

func Delete(pattern string, handler http.Handler) {
	defaultServeMux.Delete(pattern, handler)
}

func DeleteFunc(pattern string, handle http.HandlerFunc) {
	defaultServeMux.DeleteFunc(pattern, handle)
}

func Put(pattern string, handler http.Handler) {
	defaultServeMux.Put(pattern, handler)
}

func PutFunc(pattern string, handle http.HandlerFunc) {
	defaultServeMux.PutFunc(pattern, handle)
}

func Any(pattern string, handler http.Handler) {
	defaultServeMux.Any(pattern, handler)
}

func AnyFunc(pattern string, handle http.HandlerFunc) {
	defaultServeMux.AnyFunc(pattern, handle)
}

// 默认供http.LintenAndServe()调用的函数。
func DefaultHandle(w http.ResponseWriter, r *http.Request) {
	ok := defaultHostMux.serveHTTP(w, r)
	if !ok {
		defaultServeMux.serveHTTP(w, r)
	}
}

func ListenAndServe(addr string) {
	http.ListenAndServe(addr, http.HandlerFunc(DefaultHandle))
}
