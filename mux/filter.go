// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mux

import (
	"net/http"
)

// 过滤器组合
type filters []http.Handler

// 过滤器组合也实现了http.Handler接口。
func (f filters) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, filter := range f {
		filter.ServeHTTP(w, r)
	}
}

// 为某一个单独的handler添加前置过滤器
//  srv.Handler("/api/", mux.Before(apiHandler, AuthHandler))
func Before(handler http.Handler, filter http.HandlerFunc) filters {
	if f, ok := handler.(filters); ok {
		ret := make(filters, 0, len(f)+1)
		return append(append(ret, f...), handler)
	}

	return filters{
		filter,
		handler,
	}
}

// 为某一个单独的handler添加后置过滤器
func After(handler http.Handler, filter http.HandlerFunc) filters {
	if f, ok := handler.(filters); ok {
		return append(f, filter)
	}

	return filters{
		handler,
		filter,
	}
}
