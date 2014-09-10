// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mux

import (
	"net/http"
)

// 这是一个控制器的封装示例
type controller struct {
	handler http.Handler
}

func (c *controller) before() {
	// todo
}

func (c *controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.before()
	// 在这里做一些其它事
	c.handler.ServeHTTP(w, r)
}
