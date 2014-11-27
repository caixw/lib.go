// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
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
