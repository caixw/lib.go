// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"reflect"
	"testing"

	"github.com/caixw/lib.go/assert"
)

type modelUser struct {
	Id       int    `orm:"name(id);ai(1,2);"`
	Email    string `orm:"unique(unique_index);nullable;pk(pk_name)"`
	Username string `orm:"index(index)"`
	Group    int    `orm:"name(group);"`

	Regdate int `orm:"-"`
}

func TestModel(t *testing.T) {
	a := assert.New(t)

	// todo 正确声明第二个参数！！
	m, err := NewModel(&modelUser{})
	a.NotError(err).NotNil(m)

	// cols
	idCol, found := m.Cols["id"] // 指定名称为小写
	a.True(found)

	emailCol, found := m.Cols["Email"] // 未指定别名，与字段名相同
	a.True(found).True(emailCol.Nullable)

	usernameCol, found := m.Cols["Username"]
	a.True(found)

	_, found = m.Cols["group"]
	a.True(found)

	regdate, found := m.Cols["Regdate"]
	a.False(found).Nil(regdate)

	// index
	index, found := m.KeyIndexes["index"]
	a.True(found).Equal(usernameCol, index[0])

	// ai
	a.Equal(m.AI.Col, idCol)

	// 主键应该和自增列相同
	a.NotNil(m.PK).Equal(m.PK[0], idCol)

	// unique_index
	unique, found := m.UniqueIndexes["unique_index"]
	a.True(found).Equal(unique[0], emailCol)
}

func TestColumnEqual(t *testing.T) {
	a := assert.New(t)
	// 断言c1,c2相等
	T := func(c1, c2 *Column) {
		a.True(c1.Equal(c2)).True(c2.Equal(c1))
	}
	// 断言c1,c2不相等
	F := func(c1, c2 *Column) {
		a.False(c1.Equal(c2)).False(c2.Equal(c1))
	}

	c1 := &Column{}
	c2 := &Column{}
	T(c1, c2)

	c1 = &Column{
		Name:       "id",
		Len1:       5,
		Len2:       0,
		Nullable:   false,
		GoType:     reflect.TypeOf(5),
		HasDefault: false,
		Default:    "",
	}

	c2 = &Column{
		Name:       "id",
		Len1:       5,
		Len2:       0,
		Nullable:   false,
		GoType:     reflect.TypeOf(5),
		HasDefault: false,
		Default:    "",
	}

	// name
	T(c1, c2)
	c1.Name = "test"
	F(c1, c2)
	c2.Name = "test"
	T(c1, c2)

	// gotype
	c1.GoType = reflect.TypeOf(1)
	c2.GoType = reflect.TypeOf(int8(2))
	F(c1, c2)
	c2.GoType = reflect.TypeOf(2)
	T(c1, c2)

	c1.Default = "abc"
	c2.Default = "def"
	T(c1, c2)
	c1.HasDefault = true
	F(c1, c2)
	c2.HasDefault = true
	F(c1, c2)
	c2.Default = "abc"
	T(c1, c2)
	F(c1, nil)
}

func TestColumnsEqual(t *testing.T) {
	a := assert.New(t)
	T := func(v1, v2 []*Column) {
		a.True(ColumnsEqual(v1, v2))
	}
	F := func(v1, v2 []*Column) {
		a.False(ColumnsEqual(v1, v2))
	}

	c1 := &Column{}
	c2 := &Column{}

	v1 := []*Column{}
	v2 := []*Column{}
	T(v1, v2)

	v1 = append(v1, c1)
	v2 = append(v2, c1)
	T(v1, v2)
	v2 = append(v2, c2)
	F(v1, v2)

	v2 = v2[:0]
	v2 = append(v2, c2)
	T(v1, v2)

	c1.Name = "Name"
	c2.Name = "Name"
	v1 = append(v1[:0], c1)
	v2 = append(v2[:0], c2)
	T(v1, v2)
}

func TestAutoIncrEqual(t *testing.T) {
	a := assert.New(t)
	T := func(a1, a2 *AutoIncr) {
		a.True(a1.Equal(a2)).
			True(a2.Equal(a1))
	}
	F := func(a1, a2 *AutoIncr) {
		a.False(a1.Equal(a2)).
			False(a2.Equal(a1))
	}

	a1 := &AutoIncr{}
	a2 := &AutoIncr{}

	T(a1, a1)
	T(a1, a2)
	a1.Start = 5
	F(a1, a2)
	a2.Start = 5
	T(a1, a2)

	a1.Col = &Column{Name: "abc"}
	F(a1, a2)
	a2.Col = &Column{}
	F(a1, a2)
	a2.Col.Name = "abc"
	T(a1, a2)
}
