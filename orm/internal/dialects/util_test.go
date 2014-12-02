// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package dialects

import (
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestMysqlLimitSQL(t *testing.T) {
	a := assert.New(t)

	sql, args := mysqlLimitSQL(5, 0)
	a.Equal(sql, " LIMIT ? OFFSET ? ").
		Equal(args, []interface{}{5, 0})
}

func TestOracleLimitSQL(t *testing.T) {
	a := assert.New(t)

	sql, args := oracleLimitSQL(5, 0)
	a.Equal(sql, " OFFSET ? ROWS FETCH NEXT ? ROWS ONLY ").
		Equal(args, []interface{}{0, 5})
}
