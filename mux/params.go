// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package mux

import (
	"regexp"
	"strings"
)

// 表示从request.URL.Path中传递过来的参数。
type param struct {
	Items []string          // path中以/分隔的种个字段
	Args  map[string]string // 从正规中获取被命名的值
}

// 分析客户端的请求路径到param对象中
//
// patternC 路由的正规表达式
// path http.Request.URL.Path中的值，不包含query的内容。
func parseParam(patternC *regexp.Regexp, path string) *param {
	p := &param{
		Items: strings.Split(path, "/"),
		Args:  make(map[string]string),
	}

	subexps := patternC.SubexpNames()
	args := patternC.FindStringSubmatch(path)
	for index, name := range subexps {
		if name == "" {
			continue
		}

		p.Args[name] = args[index]
	}

	return p
}
