// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sqlbuilder

import (
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestDelete(t *testing.T) {
	a := assert.New(t)
	style := assert.StyleTrim | assert.StyleSpace

	e, err := newFakeDB()
	a.NotError(err).NotNil(e)

	d := NewDelete(e)
	a.NotNil(d)

	d.Table("#user").
		And("username like ?", "%admin%").
		OrIn("uid", 1, 2, 3, 4, 5).
		AndBetween(`{group}`, 1, 10)
	wont := "DELETE FROM prefix_user WHERE(username like ?) OR(uid IN(?,?,?,?,?)) AND([group] BETWEEN ? AND ?)"
	a.StringEqual(d.sqlString(true), wont, style)
}

func TestUpdate(t *testing.T) {
	a := assert.New(t)
	style := assert.StyleTrim | assert.StyleSpace

	e, err := newFakeDB()
	a.NotError(err).NotNil(e)

	u := NewUpdate(e)
	a.NotNil(u)

	u.Table("user").
		Columns("password", "username", `{group}`).
		And("id=?").
		Or(`{group}=?`)
	wont := "UPDATE user SET password=?,username=?,[group]=? WHERE(id=?) OR([group]=?)"
	a.StringEqual(u.sqlString(true), wont, style)
}

func TestInsert(t *testing.T) {
	a := assert.New(t)
	style := assert.StyleTrim | assert.StyleSpace

	e, err := newFakeDB()
	a.NotError(err).NotNil(e)

	i := NewInsert(e)
	a.NotNil(i)

	i.Table("#user").
		Columns("uid", "username", `{password}`).
		Columns("group", "age")
	wont := "INSERT INTO prefix_user(uid,username,[password],group,age) VALUES(?,?,?,?,?)"
	a.StringEqual(i.sqlString(true), wont, style).
		Equal(len(i.vals), 0)
}
