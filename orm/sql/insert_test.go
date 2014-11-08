// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"testing"

	"github.com/caixw/lib.go/assert"
	_ "github.com/caixw/lib.go/orm/dialect/test"
)

func TestInsert(t *testing.T) {
	a := assert.New(t)

	e, err := newDB()
	a.NotError(err).NotNil(e)

	i := NewInsert(e)
	a.NotNil(i)

	i.Table("table.user").
		Columns("uid", "username", `"password"`).
		Columns("group", "age")
	wont := "INSERT INTO prefix_user(uid,username,`password`,group,age) VALUES(?,?,?,?,?)"
	a.Equal(i.SQLString(true), wont).
		Equal(len(i.vals), 0)

	// Data方法传递map参数，无法确定元素顺序。
	/*
		    i.Reset()
			i.Table("table.user").
				Data(map[string]interface{}{"uid": 1, "username": "admin"}).
				Data(map[string]interface{}{"age": 100, "password": "123"})
			wont = "INSERT INTO prefix_user(uid,username,age,password) VALUES(?,?,?,?)"
			a.Equal(i.SQLString(true), wont).
				Equal(i.vals, []interface{}{1,"admin",100,"123"})
	*/

}
