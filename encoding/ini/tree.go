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
	"github.com/caixw/lib.go/encoding/tag"
)

// elem.set的默认实现。
func defaultSetFunc(val string, v reflect.Value) error {
	return conv.To(val, v)
}

var defaultSet = reflect.ValueOf(defaultSetFunc)

// elem.get的默认实现。
func defaultGetFunc(v reflect.Value) (string, error) {
	return conv.String(v.Interface())
}

var defaultGet = reflect.ValueOf(defaultGetFunc)

type elem struct {
	val reflect.Value

	// 将字符串转换并保存到elem.val中，函数原型为参照defaultSetFunc
	set reflect.Value

	// 将elem.val值转换成字符串并返回，函数原型参照defaultGetFunc
	get reflect.Value
}

// 调用elem.set函数
func (e *elem) callSet(val string) error {
	ret := e.set.Call([]reflect.Value{reflect.ValueOf(val), reflect.ValueOf(e.val)})

	if i := ret[0].Interface(); i == nil {
		return nil
	} else {
		return i.(error)
	}
}

// 调用elem.get函数。
func (e *elem) callGet() (string, error) {
	ret := e.get.Call([]reflect.Value{reflect.ValueOf(e.val)})

	val := ret[0].Interface().(string)

	if err := ret[1].Interface(); err == nil {
		return val, nil
	} else {
		return val, err.(error)
	}
}

// 将一个对象分析成树形目录结构。
//
// NOTE: ini只有两层结构。
type tree struct {
	elems map[string]*elem
	nodes map[string]map[string]*elem
}

// 向tree.nodes中添加元素。
// originV v的原始值，比如originV可能是指针，而v绝对不会是指针。
func (t *tree) addNode(field reflect.StructField, v, originV reflect.Value) error {
	name := field.Name
	if tmp, found := tag.Get(field.Tag.Get("ini"), "name"); found {
		name = tmp[0]
	}

	t.nodes[name] = make(map[string]*elem)
	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Ptr {
			f = f.Elem()
		}

		if f.Kind() == reflect.Ptr {
			return fmt.Errorf("[%v]的值为指针的指针", name)
		}
		if f.Kind() == reflect.Struct {
			return errors.New("不支持多层嵌套")
		}

		if err := addElem(t.nodes[name], typ.Field(i), v.Field(i), originV); err != nil {
			return err
		}
	}

	return nil
}

// 向elems中添加元素。
//
// 需要确保传递的v参数不是一个kind为reflect.Ptr和reflect.Struct的值。
// note:不要对parent做Value.Elem()处理，否则可能会无法取得自定义的get,set函数
func addElem(elems map[string]*elem, field reflect.StructField, v, parent reflect.Value) error {
	if !parent.IsValid() {
		return errors.New("无效的parent值")
	}

	e := &elem{val: v, set: defaultSet, get: defaultGet}
	name := field.Name

	// 提取tag内容
	tags := tag.Parse(field.Tag.Get("ini"))
	for key, vals := range tags {
		if vals == nil || len(vals[0]) == 0 { // 过滤第二个值为空的情况
			continue
		}
		switch strings.ToLower(key) {
		case "name":
			name = vals[0]
		case "set":
			if method := parent.MethodByName(vals[0]); !method.IsValid() {
				return fmt.Errorf("未知的Set函数:[%v]", vals[0])
			} else {
				e.set = method
			}
		case "get":
			if method := parent.MethodByName(vals[0]); !method.IsValid() {
				return fmt.Errorf("未知的get函数:[%v]", vals[0])
			} else {
				e.get = method
			}
		default:
			return fmt.Errorf("未知的tag字段:[%v]", key)
		}
	} // end for tags

	elems[name] = e

	return nil
}

// 从reader中初始化tree所代码的对象。
func (t *tree) unmarshal(reader *Reader) error {
	elems := t.elems

	for {
		token, err := reader.Token()
		if err != nil {
			return err
		}

		switch token.Type {
		case EOF:
			return nil
		case Element:
			elem, found := elems[token.Key]
			if !found { // 该节点不存在于对象内，忽略
				continue
			}

			if err := elem.callSet(token.Value); err != nil {
				return err
			}
		case Section:
			if es, found := t.nodes[token.Value]; found { // 该section不存在，忽略
				elems = es
			}
		} // end switch
	} // end for
}

// 将tree对象输出到writer中。
func (t *tree) marshal(w *Writer) error {
	if err := writeElems(t.elems, w); err != nil {
		return err
	}

	for index, n := range t.nodes {
		w.AddSection(index)
		if err := writeElems(n, w); err != nil {
			return err
		}
	}
	return nil
}

// 将elems写入writer对象中。
func writeElems(elems map[string]*elem, w *Writer) error {
	for k, v := range elems {
		val, err := v.callGet()
		if err != nil {
			return err
		}
		w.AddElement(k, val)
	}
	return nil
}

// 从obj对象构建tree结构。
func scan(obj interface{}) (*tree, error) {
	objv := reflect.ValueOf(obj)
	v := objv
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Ptr {
		return nil, errors.New("不接受指针的指针类型")
	}

	ret := &tree{
		elems: make(map[string]*elem),
		nodes: make(map[string]map[string]*elem),
	}

	t := v.Type()
	var err error
	for i := 0; i < t.NumField(); i++ {
		vf := v.Field(i)
		if vf.Kind() == reflect.Ptr {
			vf = vf.Elem()
		}

		if vf.Kind() == reflect.Ptr {
			return nil, fmt.Errorf("成员[%v]为指针的指针", t.Field(i).Name)
		}

		if vf.Kind() == reflect.Struct {
			err = ret.addNode(t.Field(i), vf, v.Field(i))
		} else {
			err = addElem(ret.elems, t.Field(i), vf, objv)
		}

		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}
