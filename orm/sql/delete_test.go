// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestDelete(t *testing.T) {
	a := assert.New(t)

	e, err := newDB()
	a.NotError(err).NotNil(e)

	d := NewDelete(e)
	a.NotNil(d)

	d.Table("table.user").
		And("username like ?", "%admin%").
		OrIn("uid", 1, 2, 3, 4, 5).
		AndBetween(`"group"`, 1, 10)
	wont := "DELETE FROM prefix_user WHERE(username like ?) OR(uid IN(?,?,?,?,?)) AND(`group` BETWEEN ? AND ?)"
	a.Equal(d.SQLString(true), wont)
}
