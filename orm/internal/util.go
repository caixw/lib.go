// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package internal

import (
	"bytes"

	"github.com/caixw/lib.go/orm/core"
)

// mysq系列数据库分页语法的实现。支持以下数据库：
// MySQL, H2, HSQLDB, Postgres, SQLite3
func mysqlLimitSQL(limit, offset int) (string, []interface{}) {
	return " LIMIT ? OFFSET ? ", []interface{}{limit, offset}
}

// oracle系列数据库分页语法的实现。支持以下数据库：
// Derby, SQL Server 2012, Oracle 12c, the SQL 2008 standard
func oracleLimitSQL(limit, offset int) (string, []interface{}) {
	return " OFFSET ? ROWS FETCH NEXT ? ROWS ONLY ", []interface{}{offset, limit}
}

// 用于产生在createTable中使用的普通列信息表达式，不包含autoincrement的关键字。
func createColSQL(buf *bytes.Buffer, col *core.Column, d base) {
	// col_name VARCHAR(100) NOT NULL DEFAULT 'abc'
	d.quote(buf, col.Name)
	buf.WriteByte(' ')

	// 写入字段类型
	d.sqlType(buf, col)

	if !col.Nullable {
		buf.WriteString(" NOT NULL")
	}

	if col.HasDefault {
		buf.WriteString(" DEFAULT '")
		buf.WriteString(col.Default)
		buf.WriteByte('\'')
	}
}

// 产生pk语句
func createPKSQL(buf *bytes.Buffer, cols []*core.Column, pkName string, d base) {
	//CONSTRAINT pk_name PRIMARY KEY (id,lastName)
	buf.WriteString(" CONSTRAINT ")
	d.quote(buf, pkName)
	buf.WriteString(" PRIMARY KEY(")
	for _, col := range cols {
		d.quote(buf, col.Name)
		buf.WriteByte(',')
	}
	buf.UnreadByte() // 去掉最后一个逗号

	buf.WriteByte(')')
}

// create table语句中的索引约束部分的语句。包括pk,unique
func createUniqueSQL(buf *bytes.Buffer, cols []*core.Column, indexName string, d base) {
	//CONSTRAINT unique_name UNIQUE (id,lastName)
	buf.WriteString(" CONSTRAINT ")
	d.quote(buf, indexName)
	buf.WriteString(" UNIQUE(")
	for _, col := range cols {
		d.quote(buf, col.Name)
		buf.WriteByte(',')
	}
	buf.UnreadByte() // 去掉最后一个逗号

	buf.WriteByte(')')
}

func createFKSQL(buf *bytes.Buffer, fk *core.ForeignKey, fkName string, d base) {
	//CONSTRAINT fk_name FOREIGN KEY (id) REFERENCES user(id)
	buf.WriteString(" CONSTRAINT ")
	d.quote(buf, fkName)

	buf.WriteString(" FOREIGN KEY(")
	d.quote(buf, fk.Col.Name)

	buf.WriteString(") REFERENCES ")
	d.quote(buf, fk.RefTableName)

	buf.WriteByte('(')
	d.quote(buf, fk.RefColName)
	buf.WriteByte(')')
}

func createCheckSQL(buf *bytes.Buffer, expr, chkName string, d base) {
	//CONSTRAINT chk_name CHECK (id>0 AND username='admin')
	buf.WriteString(" CONSTRAINT ")

	d.quote(buf, chkName)

	buf.WriteString(" CHECK(")
	buf.WriteString(expr)
	buf.WriteByte(')')
}

func upgradeCols(model *core.Model, d base, db core.DB) error {
	dbCols, err := d.getCols(db, model.Name)
	if err != nil {
		return err
	}

	// 转换成map，仅用到键名，键值一律置空
	dbColsMap := make(map[string]interface{}, len(dbCols))
	for _, col := range dbCols {
		dbColsMap[col] = nil
	}

	buf := bytes.NewBufferString("ALTER TABLE ")
	d.quote(buf, model.Name)
	size := buf.Len()

	// 将model中的列信息作用于数据库中的表，
	// 并将过滤dbCols中的列，只剩下不存在于model中的字段。
	for colName, col := range model.Cols {
		buf.Truncate(size)

		if _, found := dbColsMap[colName]; !found {
			buf.WriteString(" ADD ")
		} else {
			buf.WriteString(" ALTER COLUMN ")
			delete(dbColsMap, colName)
		}

		createColSQL(buf, col, d)

		if _, err := db.Exec(buf.String()); err != nil {
			return err
		}
	}

	if len(dbCols) == 0 {
		return nil
	}

	// 删除已经不存在于model中的字段。
	buf.Truncate(size)
	buf.WriteString(" DROP COLUMN ")
	size = buf.Len()
	for name, _ := range dbColsMap {
		buf.Truncate(size)
		buf.WriteString(name)
		if _, err := db.Exec(buf.String()); err != nil {
			return err
		}
	}

	return nil
}
