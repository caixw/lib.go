// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package dialects

import (
	"bytes"
	"testing"

	"github.com/caixw/lib.go/assert"
)

var _ base = &sqlite3{}

var s = &sqlite3{}

func TestSqlite3GetDBName(t *testing.T) {
	a := assert.New(t)

	a.Equal(s.GetDBName("./dbname.db"), "dbname")
	a.Equal(s.GetDBName("./dbname"), "dbname")
	a.Equal(s.GetDBName("abc/dbname.abc"), "dbname")
	a.Equal(s.GetDBName("dbname"), "dbname")
	a.Equal(s.GetDBName(""), "")
}

// mysql.quote() & mysql.QuoteStr()
func TestSqlite3QuoteQuoteStr(t *testing.T) {
	a := assert.New(t)
	l, r := s.QuoteStr()
	buf := bytes.NewBufferString("")
	s.quote(buf, "test")

	a.Equal(l+"test"+r, buf.String())
}
