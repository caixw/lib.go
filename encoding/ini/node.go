// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ini

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/caixw/lib.go/conv"
)

// node.set的默认实现。
func defaultSetFunc(val string, v reflect.Value) error {
	return conv.To(val, v)
}

var defaultSet = reflect.ValueOf(defaultSetFunc)

// node.get的默认实现。
func defaultGetFunc(v reflect.Value) (string, error) {
	return conv.String(v.Interface())
}

var defaultGet = reflect.ValueOf(defaultGetFunc)

// ini 树形节点。
type node struct {
	// 节点的值。若存在子节点，则不保证此值的正确性。
	val reflect.Value

	// 将字符串转换成符合val要求的类型，并保存到Val中，函数原型为：
	// func(val string, v reflect.Value)error
	//
	// 若存在子节点，则此值无意义
	set reflect.Value

	// 将Val值转换成字符串并返回，函数原型为：
	// func(v reflect.Value) (string,error)
	//
	// 若存在子节点，则此值无意义
	get reflect.Value

	// 子节点。
	nodes map[string]*node

	// 节点名称
	name string
}

// 是否拥有子节点。
func (n *node) hasChild() bool {
	return len(n.nodes) > 0
}

// 调用节点的set函数
func (n *node) callSet(val string) error {
	if !n.set.IsValid() {
		return nil
	}

	ret := n.set.Call([]reflect.Value{reflect.ValueOf(val), reflect.ValueOf(n.val)})

	if i := ret[0].Interface(); i == nil {
		return nil
	} else {
		return i.(error)
	}
}

// 调用节点的get函数。
func (n *node) callGet() (string, error) {
	ret := n.get.Call([]reflect.Value{reflect.ValueOf(n.val)})

	val := ret[0].Interface().(string)

	if err := ret[1].Interface(); err == nil {
		return val, nil
	} else {
		return val, err.(error)
	}
}

// 声明一个新节点。
//
// field 该节点对应的reflect.StructField；
// v 该节点对应的reflect.Value；
// parent 该节点的父级reflect.Value，若不存在，则传递一个空值。
func newNode(field reflect.StructField, v, parent reflect.Value) (*node, error) {
	n := &node{
		val:  v,
		set:  defaultSet,
		get:  defaultGet,
		name: field.Name,
	}

	tag := field.Tag.Get("ini")
	if len(tag) > 0 {
		if err := parseTag(tag, parent, n); err != nil {
			return nil, err
		}
	}

	// 分析子节点
	if v.Kind() == reflect.Struct {
		n.nodes = make(map[string]*node)
		typ := v.Type()
		for i := 0; i < typ.NumField(); i++ {
			item, err := newNode(typ.Field(i), v.Field(i), v)
			if err != nil {
				return nil, err
			}
			n.nodes[item.name] = item
		}
	}

	return n, nil
}

// 分析tag字符串，将相关内容写入到n中
// tag的格式为："name:val;name2:val2"
func parseTag(tag string, p reflect.Value, n *node) error {
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		kv := strings.Split(part, ":")
		if len(kv) != 2 {
			return fmt.Errorf("tag格式不正确:[%v]", tag)
		}

		switch strings.ToLower(kv[0]) {
		case "name":
			n.name = kv[1]
		case "set":
			if p.IsValid() {
				continue
			}
			n.set = p.MethodByName(kv[1])
			if n.set.IsValid() {
				return errors.New("未知的Set函数")
			}
		case "get":
			if p.IsValid() {
				continue
			}
			n.get = p.MethodByName(kv[1])
			if n.get.IsValid() {
				return errors.New("未知的get函数")
			}
		default:
			return errors.New("未知的tag内容")
		}
	}
	return nil
}
