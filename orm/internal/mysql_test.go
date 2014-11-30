// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package internal

import (
	"bytes"
	"testing"

	"github.com/caixw/lib.go/assert"
)

var _ base = &mysql{}

var m = &mysql{}

func TestMysqlGetDBName(t *testing.T) {
	a := assert.New(t)

	a.Equal(m.GetDBName("root:password@/dbname"), "dbname")
	a.Equal(m.GetDBName("root:@/dbname"), "dbname")
	a.Equal(m.GetDBName("root:password@tcp(localhost:3066)/dbname"), "dbname")
	a.Equal(m.GetDBName("root:password@unix(/tmp/mysql.lock)/dbname?loc=Local"), "dbname")
	a.Equal(m.GetDBName("root:/"), "")
}

// mysql.quote() & mysql.QuoteStr()
func TestMysqlQuoteQuoteStr(t *testing.T) {
	a := assert.New(t)
	l, r := m.QuoteStr()
	buf := bytes.NewBufferString("")
	m.quote(buf, "test")

	a.Equal(l+"test"+r, buf.String())
}
