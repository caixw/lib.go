// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package dialects

import (
	"bytes"
	"testing"

	"github.com/caixw/lib.go/assert"
)

var _ base = &pq{}

var p = &pq{}

func TestPQGetDBName(t *testing.T) {
	a := assert.New(t)

	a.Equal(p.GetDBName("user=abc dbname = dbname password=abc"), "dbname")
	a.Equal(p.GetDBName("dbname=\tdbname user=abc"), "dbname")
	a.Equal(p.GetDBName("dbname=dbname\tuser=abc"), "dbname")
	a.Equal(p.GetDBName("\tdbname=dbname user=abc"), "dbname")
	a.Equal(p.GetDBName("\tdbname = dbname user=abc"), "dbname")
}

// pq.quote() & pq.QuoteStr()
func TestPQQuoteQuoteStr(t *testing.T) {
	a := assert.New(t)
	l, r := p.QuoteStr()
	buf := bytes.NewBufferString("")
	p.quote(buf, "test")

	a.Equal(l+"test"+r, buf.String())
}
