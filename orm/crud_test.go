// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"bytes"
	"testing"

	"github.com/caixw/lib.go/assert"
	"github.com/caixw/lib.go/orm/core"
)

func newDB() (core.DB, error) {
	return newEngine("fakedb1", "dataSourceName", "prefix_")
}

func TestWhereExpr(t *testing.T) {
	a := assert.New(t)
	style := assert.StyleTrim | assert.StyleSpace

	e, err := newDB()
	a.NotError(err).NotNil(e)

	w := &whereExpr{
		cond:     bytes.NewBufferString(""),
		condArgs: make([]interface{}, 0),
	}
	a.NotNil(w)

	w.build(and, `"id"=? and username=?`, 5, "abc")
	a.StringEqual(w.condString(e), " WHERE([id]=? and username=?)", style).
		Equal(w.condArgs, []interface{}{5, "abc"})

	// 重置
	w.Reset()
	a.Equal(w.cond.Len(), 0).
		Equal(len(w.condArgs), 0)

	// Between
	w.AndBetween("age", 5, 6)
	a.StringEqual(w.condString(e), " WHERE(age BETWEEN ? AND ?)", style).
		Equal(w.condArgs, []interface{}{5, 6})

	// In函数不指定数据，会触发panic
	a.Panic(func() { w.In("id") })

	w.Reset()
	w.AndIsNull("age")
	a.StringEqual(w.condString(e), " WHERE(age IS NULL)", style).
		Equal(len(w.condArgs), 0)

	w.Reset()
	w.OrIsNotNull("age")
	a.StringEqual(w.condString(e), " WHERE(age IS NOT NULL)", style).
		Equal(len(w.condArgs), 0)

	w.Reset()
	w.And("id=?", 5).AndIn("age", 7, 8, 9).OrIsNotNull("group")
	a.StringEqual(w.condString(e), " WHERE(id=?) AND(age IN(?,?,?)) OR(group IS NOT NULL)", style).
		Equal(w.condArgs, []interface{}{5, 7, 8, 9})
}

func TestDelete(t *testing.T) {
	a := assert.New(t)
	style := assert.StyleTrim | assert.StyleSpace

	e, err := newDB()
	a.NotError(err).NotNil(e)

	d := newDelete(e)
	a.NotNil(d)

	d.Table("table.user").
		And("username like ?", "%admin%").
		OrIn("uid", 1, 2, 3, 4, 5).
		AndBetween(`"group"`, 1, 10)
	wont := "DELETE FROM prefix_user WHERE(username like ?) OR(uid IN(?,?,?,?,?)) AND([group] BETWEEN ? AND ?)"
	a.StringEqual(d.sqlString(true), wont, style)
}

func TestUpdate(t *testing.T) {
	a := assert.New(t)
	style := assert.StyleTrim | assert.StyleSpace

	e, err := newDB()
	a.NotError(err).NotNil(e)

	u := newUpdate(e)
	a.NotNil(u)

	u.Table("user").
		Columns("password", "username", `"group"`).
		And("id=?").
		Or(`"group"=?`)
	wont := "UPDATE user SET password=?,username=?,[group]=? WHERE(id=?) OR([group]=?)"
	a.StringEqual(u.sqlString(true), wont, style)
}

func TestInsert(t *testing.T) {
	a := assert.New(t)
	style := assert.StyleTrim | assert.StyleSpace

	e, err := newDB()
	a.NotError(err).NotNil(e)

	i := newInsert(e)
	a.NotNil(i)

	i.Table("table.user").
		Columns("uid", "username", `"password"`).
		Columns("group", "age")
	wont := "INSERT INTO prefix_user(uid,username,[password],group,age) VALUES(?,?,?,?,?)"
	a.StringEqual(i.sqlString(true), wont, style).
		Equal(len(i.vals), 0)
}
