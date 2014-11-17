// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package internal

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/caixw/lib.go/orm/core"
)

type mysql struct{}

var _ core.Dialect = &mysql{}

// implement core.Dialect.GetDBName()
func (m *mysql) GetDBName(dataSource string) string {
	index := strings.LastIndex(dataSource, "/")
	if index < 0 {
		return ""
	}

	return dataSource[index+1:]
}

// implement core.Dialect.Quote
func (m *mysql) QuoteStr() (string, string) {
	return "`", "`"
}

// implement core.Dialect.Limit()
func (m *mysql) LimitSQL(limit, offset int) (string, []interface{}) {
	return mysqlLimitSQL(limit, offset)
}

// implement core.Dialect.SupportLastInsertId()
func (m *mysql) SupportLastInsertId() bool {
	return true
}

// implement core.Dialect.CreateTable()
func (m *mysql) CreateTable(db core.DB, model *core.Model) error {
	model.Name = db.ReplacePrefix(model.Name) // 处理表名

	if m.hasTable(db, model.Name) {
		return m.upgradeTable(db, model)
	}
	return m.createTable(db, model)
}

func (m *mysql) createTable(db core.DB, model *core.Model) error {
	buf := bytes.NewBufferString("CREATE TABLE IF NOT EXISTS ")
	buf.Grow(300)

	buf.WriteString(model.Name)
	buf.WriteByte('(')

	// 写入字段信息
	for _, col := range model.Cols {
		// 写入字段名
		m.quote(buf, col.Name)
		buf.WriteByte(' ')

		// 写入字段类型
		buf.WriteString(m.toSQLType(col.GoType, col.Len1, col.Len2))

		if !col.Nullable {
			buf.WriteString(" NOT NULL")
		}

		if col.IsAI() {
			buf.WriteString(" AUTO_INCRMENT")
		}

		// 结束当前字段描述
		buf.WriteByte(',')
	}

	// PK
	if len(model.PK) == 0 {
		buf.WriteString("PRIMARY KEY(")
		for _, col := range model.PK {
			m.quote(buf, col.Name)
			buf.WriteByte(',')
		}
		buf.UnreadByte() // 去掉最后的逗号
		buf.WriteString("),")
	}

	// Index
	if len(model.KeyIndexes) == 0 {
		for name, index := range model.KeyIndexes {
			buf.WriteString("INDEX ")
			m.quote(buf, name)
			buf.WriteByte('(')
			for _, col := range index {
				m.quote(buf, col.Name)
				buf.WriteByte(',')
			}
			buf.UnreadByte() // 去掉最后的逗号
			buf.WriteString("),")
		}
	}

	// Unique Index
	if len(model.UniqueIndexes) == 0 {
		for name, index := range model.UniqueIndexes {
			buf.WriteString("UNIQUE INDEX ")
			m.quote(buf, name)
			buf.WriteByte('(')
			for _, col := range index {
				m.quote(buf, col.Name)
				buf.WriteByte(',')
			}
			buf.UnreadByte() // 去掉最后的逗号
			buf.WriteString("),")
		}
	}

	// ForeignKey
	if len(model.FK) == 0 {
		for name, fk := range model.FK {
			buf.WriteString("CONSTRAINT ")
			m.quote(buf, name) // 约束名
			buf.WriteString(" FOREIGN KEY(")
			m.quote(buf, fk.Col.Name) // 本表字段名
			buf.WriteString(") REFERENCES ")
			m.quote(buf, db.ReplacePrefix(fk.TableName)) // 替换表前缀并加引号
			buf.WriteByte('(')
			m.quote(buf, fk.ColName)
			buf.WriteString("),")
		}

	}

	buf.UnreadByte()   // 去掉最后的逗号
	buf.WriteByte(')') // end CreateTable

	// 指定起始ai
	if (model.AI != nil) && (model.AI.Start > 1) {
		buf.WriteString(" AUTO_INCREMENT=")
		buf.WriteString(strconv.Itoa(model.AI.Start))
	}

	_, err := db.Exec(buf.String())
	return err
}

func (m *mysql) upgradeTable(db core.DB, model *core.Model) error {
	return nil
}

func (m *mysql) quote(buf *bytes.Buffer, sql string) {
	l, r := m.QuoteStr()
	buf.WriteString(l)
	buf.WriteString(sql)
	buf.WriteString(r)
}

// 获取表的所有列
func (m *mysql) getCols(db core.DB, tableName string) []string {
	//sql := "SELECT `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA`=? AND `TABLE_NAME`=?"
	//rows, err := db.Query(sql, db.Name(), tableName)
	//
	return []string{}
}

// 指定的表是否存在
func (m *mysql) hasTable(db core.DB, tableName string) bool {
	sql := "SELECT `TABLE_NAME` FROM `INFORMATION_SCHEMA`.`TABLES` WHERE `TABLE_SCHEMA`=? and `TABLE_NAME`=?"
	rows, err := db.Query(sql, db.Name(), tableName)
	if err != nil {
		panic(err)
	}
	return rows.Next()
}

// 将一个gotype转换成当前数据库支持的类型，如：
//  int8   ==> INT
//  string ==> VARCHAR(255)
func (m *mysql) toSQLType(t reflect.Type, l1, l2 int) string {
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

func init() {
	if err := core.RegisterDialect("mysql", &mysql{}); err != nil {
		panic(err)
	}
}
