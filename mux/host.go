// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mux

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"
)

type hostEntry struct {
	patternC *regexp.Regexp
	handler  http.Handler
}

// 域名路由器
type HostMux struct {
	sync.Mutex
	entries map[string]*hostEntry
}

// 添加域名路由
//
// 支持正则表达式
func (host *HostMux) Handle(pattern string, handler http.Handler) {
	host.Lock()
	defer host.Unlock()

	if _, found := host.entries[pattern]; found {
		msg := fmt.Sprintf("域名路由器[%s]已经存在", pattern)
		panic(msg)
	}

	host.entries[pattern] = &hostEntry{
		patternC: regexp.MustCompile(pattern),
		handler:  handler,
	}
}

// 添加域名路由
func (host *HostMux) HandleFunc(pattern string, handle http.HandlerFunc) {
	host.Handle(pattern, http.Handler(handle))
}

// implement http.Handler
func (host *HostMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host.serveHTTP(w, r)
}

func (host *HostMux) serveHTTP(w http.ResponseWriter, r *http.Request) bool {
	for _, entry := range host.entries {
		if !entry.patternC.MatchString(r.Host) {
			continue
		}

		ctx := GetContext(r)
		ctx.Add("domains", parseCaptures(entry.patternC, r.URL.Host))
		entry.handler.ServeHTTP(w, r)
		return true
	}
	return false
}
