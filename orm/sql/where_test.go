// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestwhereExpr(t *testing.T) {
	a := assert.New(t)

	e, err := newDB()
	a.NotError(err).NotNil(e)

	w := &whereExpr{}
	a.NotNil(w)

	w.build(and, `"id"=? and username=?`, 5, "abc")
	a.Equal(w.condString(e), " WHERE(`id`=? and username=?)").
		Equal(w.condArgs, []interface{}{5, "abc"})

	// 重置
	w.Reset()
	a.Equal(w.cond.Len(), 0).
		Equal(len(w.condArgs), 0)

	// Between
	w.AndBetween("age", 5, 6)
	a.Equal(w.condString(e), " WHERE(age BETWEEN ? AND ?)").
		Equal(w.condArgs, []interface{}{5, 6})

	// In函数不指定数据，会触发panic
	a.Panic(func() { w.In("id") })

	w.Reset()
	w.AndIsNull("age")
	a.Equal(w.condString(e), " WHERE(age IS NULL)").
		Equal(len(w.condArgs), 0)

	w.Reset()
	w.OrIsNotNull("age")
	a.Equal(w.condString(e), " WHERE(age IS NOT NULL)").
		Equal(len(w.condArgs), 0)

	w.Reset()
	w.And("id=?", 5).AndIn("age", 7, 8, 9).OrIsNotNull("group")
	a.Equal(w.condString(e), " WHERE(id=?) AND(age IN(?,?,?)) OR(group IS NOT NULL)").
		Equal(w.condArgs, []interface{}{5, 7, 8, 9})
}
