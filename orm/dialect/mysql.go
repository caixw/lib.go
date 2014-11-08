// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package dialect

import (
	"fmt"
	"reflect"
)

type mysql struct{}

var _ Dialect = &mysql{}

// Dialect.Quote
func (m *mysql) Quote() (string, string) {
	return "`", "`"
}

// Dialect.ToSqlType
func (m *mysql) ToSqlType(t reflect.Type, l1, l2 int) string {
	switch t.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int8:
		return "TINYINT"
	case reflect.Int16:
		return "SMALLINT"
	case reflect.Int32:
		return "INT"
	case reflect.Int64, reflect.Int:
		return "BIGINT"
	case reflect.Uint8:
		return "TINYINT UNSIGNED"
	case reflect.Uint16:
		return "SMALLINT UNSIGNED"
	case reflect.Uint32:
		return "INT UNSIGNED"
	case reflect.Uint64, reflect.Uint, reflect.Uintptr:
		return "BIGINT UNSIGNED"
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("DOUBLE(%d,%d)", l1, l2)
	case reflect.String:
		if l1 < 65533 {
			return fmt.Sprintf("VARCHAR(%d)", l1)
		} else {
			return "LONGTEXT"
		}
	case reflect.Struct: // TODO(caixw) time,nullstring等处理
	default:
	}
	return ""
}

func (m *mysql) Limit(limit, offset int) (string, []interface{}) {
	return mysqlLimit(limit, offset)
}

func init() {
	if err := Register("mysql", &mysql{}); err != nil {
		panic(err)
	}
}
