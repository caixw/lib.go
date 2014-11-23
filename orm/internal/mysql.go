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
		// 写入字段名
		m.quote(buf, col.Name)
		buf.WriteByte(' ')

		// 写入字段类型
		buf.WriteString(m.toSQLType(col.GoType, col.Len1, col.Len2))

		if !col.Nullable {
			buf.WriteString(" NOT NULL")
		}

		if col.HasDefault {
			buf.WriteString(" DEFAULT '")
			buf.WriteString(col.Default)
			buf.WriteByte('\'')
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
	model.Name = db.ReplacePrefix(model.Name)

	dbModel := m.getModel(db, model.Name)

	if model.Equal(dbModel) {
		return nil
	}

	for modelName, modelCol := range model.Cols {
		dbModelCol, found := dbModel.Cols[modelName]
		if found {
			// upd col
		}
		// todo
	}
	//todo

	return nil
}

func (m *mysql) quote(buf *bytes.Buffer, sql string) {
	buf.WriteByte('`')
	buf.WriteString(sql)
	buf.WriteByte('`')
}

// 将表转换成core.Model
func (m *mysql) getModel(db core.DB, tableName string) *core.Model {
	sql := "SELECT `COLUMN_NAME`, `IS_NULLABLE`, `COLUMN_DEFAULT`, IF(`COLUMN_DEFAULT` IS NOT NULL,'YES','NO') as HAS_DEFAULT, `COLUMN_TYPE`, `COLUMN_KEY`, `EXTRA` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ?"
	rows, err := db.Query(sql, db.Name(), tableName)
	if err != nil {
		panic(err)
	}

	mapped, err := core.Fetch2MapsString(false, rows)
	if err != nil {
		panic(err)
	}

	model := &core.Model{
		Cols:          map[string]*core.Column{},
		KeyIndexes:    map[string][]*core.Column{},
		UniqueIndexes: map[string][]*core.Column{},
		Name:          tableName,
	}

	// 将列信息导出到model.Cols中
	for _, c := range mapped {
		col := &core.Column{
			Name:       c["COLUMN_NAME"],
			Nullable:   c["IS_NULLABLE"] == "YES",
			HasDefault: c["HAS_DEFAULT"] == "YES",
			Default:    c["COLUMN_DEFAULT"],
			GoType:     m.getType(c["COLUMN_TYPE"]),
		}

		model.Cols[col.Name] = col
	}

	m.getIndexes(db, model, tableName)

	return model
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

// 将SQL类型转换成reflect.Type。如：
//  VARCHAR(255) ==> string
func (m *mysql) getType(sql string) reflect.Type {
	sql = strings.ToUpper(sql)

	left := strings.IndexByte(sql, '(')
	if left < 0 {
		return mysqlTypeMap[sql]
	}

	right := strings.IndexByte(sql, ')')
	if right < 0 {
		panic(fmt.Sprintf("获取的表结构中关于字段类型描述格式不正确：[%v]", sql))
	}

	// 去掉中间部分关于长度的描述
	index := sql[0:left] + sql[right+1:]

	return mysqlTypeMap[index]
}

var mysqlTypeMap = map[string]reflect.Type{
	"BOOLEAN":  reflect.TypeOf(true),
	"TINYINT":  reflect.TypeOf(int8(1)),
	"SMALLINT": reflect.TypeOf(int16(1)),
	"INT":      reflect.TypeOf(int32(1)),
	"BIGINT":   reflect.TypeOf(int64(1)),

	"TINYINT UNSIGNED":  reflect.TypeOf(uint8(1)),
	"SMALLINT UNSIGNED": reflect.TypeOf(uint16(1)),
	"INT UNSIGNED":      reflect.TypeOf(uint32(1)),
	"BIGINT UNSIGNED":   reflect.TypeOf(uint64(1)),

	"DOUBLE": reflect.TypeOf(float64(1)),

	"VARCHAR":  reflect.TypeOf("1"),
	"LONGTEXT": reflect.TypeOf("1"),

	// datetime
}

// 将一个gotype转换成当前数据库支持的类型，如：
//  int8   ==> INT
//  string ==> VARCHAR(255)
func (m *mysql) toSQLType(t reflect.Type, l1, l2 int) string {
	switch t.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int8:
		return m.int2SQLType("TINYINT", l1)
	case reflect.Int16:
		return m.int2SQLType("SMALLINT", l1)
	case reflect.Int32:
		return m.int2SQLType("INT", l1)
	case reflect.Int64, reflect.Int: // reflect.Int大小未知，都当作是BIGINT处理
		return m.int2SQLType("BIGINT", l1)
	case reflect.Uint8:
		return m.uint2SQLType("TINYINT", l1)
	case reflect.Uint16:
		return m.uint2SQLType("SMALLINT", l1)
	case reflect.Uint32:
		return m.uint2SQLType("INT", l1)
	case reflect.Uint64, reflect.Uint, reflect.Uintptr:
		return m.uint2SQLType("BIGINT", l1)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("DOUBLE(%d,%d)", l1, l2)
	case reflect.String:
		if l1 < 65533 {
			return fmt.Sprintf("VARCHAR(%d)", l1)
		}
		return "LONGTEXT"
	case reflect.Slice, reflect.Array:
		// 若是数组，则特殊处理[]byte与[]rune两种情况。
		k := t.Elem().Kind()
		if (k != reflect.Int8) && (k != reflect.Int32) {
			panic("不支持数组类型")
		}

		if l1 < 65533 {
			return fmt.Sprintf("VARCHAR(%d)", l1)
		}
		return "LONGTEXT"
	case reflect.Struct: // TODO(caixw) time,nullstring等处理
	default:
		panic(fmt.Sprintf("不支持的类型:[%v]", t.Name()))
	}
	return ""
}

// 将一个有符号整数名称与长度拼接SQL语句，仅供toSQLType()调用。
func (m *mysql) int2SQLType(SQLType string, l int) string {
	if l > 0 {
		return SQLType + "(" + strconv.Itoa(l) + ")"
	}
	return SQLType
}

// 将一个无符号整数名称与长度拼接SQL语句，仅供toSQLType()调用。
func (m *mysql) uint2SQLType(SQLType string, l int) string {
	if l > 0 {
		return SQLType + "(" + strconv.Itoa(l) + ") UNSIGNED"
	}
	return SQLType + " UNSIGNED"
}

func init() {
	if err := core.RegisterDialect("mysql", &mysql{}); err != nil {
		panic(err)
	}
}
