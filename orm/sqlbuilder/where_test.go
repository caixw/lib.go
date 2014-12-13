// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sqlbuilder

import (
	"bytes"
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestWhereExpr(t *testing.T) {
	a := assert.New(t)
	style := assert.StyleTrim | assert.StyleSpace

	w := &whereExpr{
		cond:     bytes.NewBufferString(""),
		condArgs: make([]interface{}, 0),
	}
	a.NotNil(w)

	w.build(and, `{id}=? and username=?`, 5, "abc")
	a.StringEqual(w.cond.String(), " WHERE({id}=? and username=?)", style).
		Equal(w.condArgs, []interface{}{5, "abc"})

	// 重置
	w.Reset()
	a.Equal(w.cond.Len(), 0).
		Equal(len(w.condArgs), 0)

	// Between
	w.AndBetween("age", 5, 6)
	a.StringEqual(w.cond.String(), " WHERE(age BETWEEN ? AND ?)", style).
		Equal(w.condArgs, []interface{}{5, 6})

	// In函数不指定数据，会触发panic
	a.Panic(func() { w.In("id") })

	w.Reset()
	w.AndIsNull("age")
	a.StringEqual(w.cond.String(), " WHERE(age IS NULL)", style).
		Equal(len(w.condArgs), 0)

	w.Reset()
	w.OrIsNotNull("age")
	a.StringEqual(w.cond.String(), " WHERE(age IS NOT NULL)", style).
		Equal(len(w.condArgs), 0)

	w.Reset()
	w.And("id=?", 5).AndIn("age", 7, 8, 9).OrIsNotNull("group")
	a.StringEqual(w.cond.String(), " WHERE(id=?) AND(age IN(?,?,?)) OR(group IS NOT NULL)", style).
		Equal(w.condArgs, []interface{}{5, 7, 8, 9})
}
