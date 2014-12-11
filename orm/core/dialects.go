// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"database/sql"
	"fmt"
	"reflect"
	"sync"
)

type dialectMap struct {
	sync.Mutex
	items map[string]Dialect
}

// 所有注册的dialect
var dialects = &dialectMap{items: make(map[string]Dialect)}

// 清空所有已经注册的dialect
func clearDialects() {
	dialects.Lock()
	defer dialects.Unlock()

	dialects.items = make(map[string]Dialect)
}

// 注册一个新的Dialect
// name值应该与sql.Open()中的driverName参数相同。
func RegisterDialect(name string, d Dialect) error {
	dialects.Lock()
	defer dialects.Unlock()

	if !isRegistedDriver(name) {
		return fmt.Errorf("该名称[%v]的driver未注册", name)
	}

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

// 指定名称的driverName是否已经被注册
func isRegistedDriver(driverName string) bool {
	drivers := sql.Drivers()
	for _, driver := range drivers {
		if driver == driverName {
			return true
		}
	}

	return false
}

// 指定名称的Dialect是否已经被注册
func IsRegistedDialect(name string) bool {
	_, found := dialects.items[name]
	return found
}

// 所有已经注册的Dialect名称列表
func RegistedDialects() []string {
	dialects.Lock()
	defer dialects.Unlock()

	ds := make([]string, 0, len(dialects.items))
	for k, _ := range dialects.items {
		ds = append(ds, k)
	}

	return ds
}

// 获取一个Dialect
func GetDialect(name string) (d Dialect, found bool) {
	dialects.Lock()
	defer dialects.Unlock()

	d, found = dialects.items[name]
	return
}
