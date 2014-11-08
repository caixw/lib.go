// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestUpdate(t *testing.T) {
	a := assert.New(t)

	e, err := newDB()
	a.NotError(err).NotNil(e)

	u := NewUpdate(e)
	a.NotNil(u)

	u.Table("user").
		Columns("password", "username", `"group"`).
		And("id=?").
		Or(`"group"=?`)
	wont := "UPDATE user SET password=?,username=?,`group`=? WHERE(id=?) OR(`group`=?)"
	a.Equal(u.SQLString(true), wont)
}
