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

// 指定的表是否存在
func (m *mysql) hasTable(db core.DB, tableName string) bool {
	sql := "SELECT `TABLE_NAME` FROM `INFORMATION_SCHEMA`.`TABLES` WHERE `TABLE_SCHEMA`=? and `TABLE_NAME`=?"
	rows, err := db.Query(sql, db.Name(), tableName)
	if err != nil {
		panic(err)
	}
	return rows.Next()
}

// 创建表
func (m *mysql) createTable(db core.DB, model *core.Model) error {
	buf := bytes.NewBufferString("CREATE TABLE IF NOT EXISTS ")
	buf.Grow(300)

	buf.WriteString(db.ReplacePrefix(model.Name))
	buf.WriteByte('(')

	// 写入字段信息
	for _, col := range model.Cols {
		createColSQL(buf, col, m)

		if col.IsAI() {
			buf.WriteString(" AUTO_INCRMENT")
		}
		buf.WriteByte(',')
	}

	// PK
	if len(model.PK) > 0 {
		createPKSQL(buf, model.PK, "pk", m)
		buf.WriteByte(',')
	}

	// Unique Index
	for name, index := range model.UniqueIndexes {
		createUniqueSQL(buf, index, name, m)
		buf.WriteByte(',')
	}

	// foreign  key
	for name, fk := range model.FK {
		fk.RefTableName = db.ReplacePrefix(fk.RefTableName)
		createFKSQL(buf, fk, name, m)
	}

	// Check
	for name, chk := range model.Check {
		createCheckSQL(buf, chk, name, m)
	}

	// key index不存在CONSTRAINT形式的语句
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

// 更新表
func (m *mysql) upgradeTable(db core.DB, model *core.Model) error {
	if err := upgradeCols(model, m, db); err != nil {
		return err
	}

	//

	return nil
}

func (m *mysql) quote(buf *bytes.Buffer, sql string) {
	buf.WriteByte('`')
	buf.WriteString(sql)
	buf.WriteByte('`')
}

// 获取表的列信息。
func (m *mysql) getCols(db core.DB, tableName string) ([]string, error) {
	sql := "SELECT `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"
	rows, err := db.Query(sql, db.Name(), tableName)
	if err != nil {
		return nil, err
	}

	return core.FetchColumnsString(false, "COLUMN_NAME", rows)
}

// 从数据库导出表的索引信息，保存到model中
func (m *mysql) getIndexes(db core.DB, model *core.Model, tableName string) {
	sql := "SELECT `INDEX_NAME`, `NON_UNIQUE`, `COLUMN_NAME` FROM `INFORMATION_SCHEMA`.`STATISTICS`" +
		" WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"
	rows, err := db.Query(sql, db.Name(), tableName)
	if err != nil {
		panic(err)
	}

	mapped, err := core.Fetch2MapsString(false, rows)
	if err != nil {
		panic(err)
	}

	for _, index := range mapped {
		name := index["INDEX_NAME"]
		col := model.Cols[index["COLUMN_NAME"]]

		if name == "PRIMARY" {
			model.PK = append(model.PK, col)
			continue
		}

		if index["NON_UNIQUE"] != "1" {
			model.KeyIndexes[name] = append(model.KeyIndexes[name], col)
		} else {
			model.UniqueIndexes[name] = append(model.UniqueIndexes[name], col)
		}
	}
}

// 将一个gotype转换成当前数据库支持的类型，如：
//  int8   ==> INT
//  string ==> VARCHAR(255)
func (m *mysql) sqlType(buf *bytes.Buffer, col *core.Column) {
	addIntLen := func() {
		if col.Len1 > 0 {
			buf.WriteByte('(')
			buf.WriteString(strconv.Itoa(col.Len1))
			buf.WriteByte(')')
		}
	}
	switch col.GoType.Kind() {
	case reflect.Bool:
		buf.WriteString("BOOLEAN")
	case reflect.Int8:
		buf.WriteString("TINYINT")
		addIntLen()
	case reflect.Int16:
		buf.WriteString("SMALLINT")
		addIntLen()
	case reflect.Int32:
		buf.WriteString("INT")
		addIntLen()
	case reflect.Int64, reflect.Int: // reflect.Int大小未知，都当作是BIGINT处理
		buf.WriteString("BIGINT")
		addIntLen()
	case reflect.Uint8:
		buf.WriteString("TINYINT")
		addIntLen()
		buf.WriteString(" UNSIGNED")
	case reflect.Uint16:
		buf.WriteString("SMALLINT")
		addIntLen()
		buf.WriteString(" UNSIGNED")
	case reflect.Uint32:
		buf.WriteString("INT")
		addIntLen()
		buf.WriteString(" UNSIGNED")
	case reflect.Uint64, reflect.Uint, reflect.Uintptr:
		buf.WriteString("BIGINT")
		addIntLen()
		buf.WriteString(" UNSIGNED")
	case reflect.Float32, reflect.Float64:
		buf.WriteString(fmt.Sprintf("DOUBLE(%d,%d)", col.Len1, col.Len2))
	case reflect.String:
		if col.Len1 < 65533 {
			buf.WriteString(fmt.Sprintf("VARCHAR(%d)", col.Len1))
		}
		buf.WriteString("LONGTEXT")
	case reflect.Slice, reflect.Array:
		// 若是数组，则特殊处理[]byte与[]rune两种情况。
		k := col.GoType.Elem().Kind()
		if (k != reflect.Int8) && (k != reflect.Int32) {
			panic("不支持数组类型")
		}

		if col.Len1 < 65533 {
			buf.WriteString(fmt.Sprintf("VARCHAR(%d)", col.Len1))
		}
		buf.WriteString("LONGTEXT")
	case reflect.Struct: // TODO(caixw) time,nullstring等处理
	default:
		panic(fmt.Sprintf("不支持的类型:[%v]", col.GoType.Name()))
	}
}

func init() {
	if err := core.RegisterDialect("mysql", &mysql{}); err != nil {
		panic(err)
	}
}
