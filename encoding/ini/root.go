// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ini

import (
	"errors"
	"fmt"
	"reflect"
)

type root struct {
	nodes map[string]*node
	trees map[string]*node
}

func newRoot(val interface{}) (*root, error) {
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.New("val参数必须为struct或是struct指针")
	}

	r := &root{nodes: make(map[string]*node), trees: make(map[string]*node)}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct {
			node, err := newNode(t.Field(i), field, reflect.Value{})
			if err != nil {
				return nil, err
			}
			r.trees[node.name] = node
		} else {
			node, err := newNode(t.Field(i), field, v)
			if err != nil {
				return nil, err
			}
			r.nodes[node.name] = node
		}
	} // end for

	return r, nil
}

// 将reader写入到r中。
func (r *root) unmarshal(reader *Reader) error {
	currNodes := r.nodes

	for {
		token, err := reader.Token()
		if err != nil {
			return err
		}

		switch token.Type {
		case EOF:
			return nil
		case Element:
			elem, found := currNodes[token.Key]
			if !found { // 该节点不存在于对象内，忽略
				continue
			}
			if elem == nil {
				return fmt.Errorf("节点的值为空:[%v]", token.Key)
			}

			if err := elem.callSet(token.Value); err != nil {
				return err
			}
		case Section:
			tree, found := r.trees[token.Value]
			if !found { // 该section不存在，忽略
				continue
			}
			currNodes = tree.nodes
		} // end switch
	} // end for
}

func (r *root) marshal(w *Writer) error {
	if err := writeNodes(r.nodes, w); err != nil {
		return err
	}

	for index, tree := range r.trees {
		w.AddSection(index)
		if err := writeNodes(tree.nodes, w); err != nil {
			return err
		}
	}
	return nil
}

func writeNodes(nodes map[string]*node, w *Writer) error {
	for k, v := range nodes {
		val, err := v.callGet()
		if err != nil {
			return err
		}
		w.AddElement(k, val)
	}
	return nil
}
