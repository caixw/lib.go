// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// tag包实现对特定格式的struct tag字符串的分析。
// 并不具有很强的通用性。
//  "id:1;unique;fun:add,1,2;"
// 将会被解析成：
//  [
//       "id"    :["1"],
//       "unique":nil,
//       "fun"   :["add","1","2"]
//  ]
package tag

import (
	"strings"
)

// 当前库的版本
const Version = "0.1.0.140910"

// 分析tag的内容，并以map的形式返回
func Parse(tag string) map[string][]string {
	ret := make(map[string][]string)

	if len(tag) == 0 {
		return nil
	}

	parts := strings.Split(tag, ";")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 1 {
			ret[kv[0]] = nil
		} else {
			ret[kv[0]] = strings.Split(kv[1], ",")
		}
	}

	return ret
}

// 从tag中查找名称为name的内容。
// 第二个参数用于判断该项是否存在。
func Get(tag, name string) ([]string, bool) {
	if len(tag) == 0 {
		return nil, false
	}

	parts := strings.Split(tag, ";")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		kv := strings.SplitN(part, ":", 2)
		if kv[0] != name {
			continue
		}
		if len(kv) == 1 {
			return nil, true
		}
		return strings.Split(kv[1], ","), true
	}

	return nil, false
}

// 查询指定名称的项是否存在，若只是查找是否存在该
// 项，使用Has()会比Get()要快上许多。
func Has(tag, name string) bool {
	if len(tag) == 0 {
		return false
	}

	parts := strings.Split(tag, ";")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		kv := strings.SplitN(part, ":", 2)
		if kv[0] == name {
			return true
		}
	}

	return false
}
