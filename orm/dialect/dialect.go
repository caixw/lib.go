// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package dialect

import (
	"fmt"
	"reflect"
	"sync"
)

// 通用但又没有统一标准的数据库功能接口。
//
// 有可能一个Dialect实例会被多个其它实例引用，不应该
// 在Dialect实例中保存状态值等内容。
type Dialect interface {
	// 对字段或是表名的引用符号
	Quote() (left, right string)

	// 将go类型转换成sql类型的表示方式
	ToSqlType(t reflect.Type, len1, len2 int) string

	// 生成limit n offset m语句
	// 返回的是对应数据库的limit语句以及语句中占位符对应的值
	Limit(limit, offset int) (sql string, vals []interface{})
}

type dialectMap struct {
	sync.Mutex
	items map[string]Dialect
}

// 所有注册的dialect
var dialects = &dialectMap{items: make(map[string]Dialect)}

// 清空所有已经注册的dialect
func clear() {
	dialects.Lock()
	defer dialects.Unlock()

	dialects.items = make(map[string]Dialect)
}

// 注册一个新的Dialect
// name值应该与sql.Open()中的driverName参数相同。
func Register(name string, d Dialect) error {
	// TODO(caixw) GO1.4 database/sql包可以查询已注册driverName
	// 列表，通过判断是否在该列表，再判断能否注册。
	dialects.Lock()
	defer dialects.Unlock()

	for k, v := range dialects.items {
		if k == name {
			return fmt.Errorf("该名称[%v]已经存在", name)
		}

		if reflect.TypeOf(d) == reflect.TypeOf(v) {
			return fmt.Errorf("该Dialect的实例已经存在，其注册名称为[%v]", k)
		}
	}

	dialects.items[name] = d
	return nil
}

// 指定名称的Dialect是否已经被注册
func IsRegisted(name string) bool {
	_, found := dialects.items[name]
	return found
}

// 所有已经注册的Dialect名称列表
func Registed() (ds []string) {
	dialects.Lock()
	defer dialects.Unlock()

	for k, _ := range dialects.items {
		ds = append(ds, k)
	}

	return
}

// 获取一个Dialect
func Get(name string) (d Dialect, found bool) {
	dialects.Lock()
	defer dialects.Unlock()

	d, found = dialects.items[name]
	return
}
