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

// create table语句中pk约束的语句
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

// create table语句中的unique约束部分的语句。
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

// create table语句中fk的约束部分的语句
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

	if len(fk.UpdateRule) > 0 {
		buf.WriteString(" ON UPDATE ")
		buf.WriteString(fk.UpdateRule)
	}

	if len(fk.DeleteRule) > 0 {
		buf.WriteString(" ON DELETE ")
		buf.WriteString(fk.DeleteRule)
	}
}

// create table语句中check约束部分的语句
func createCheckSQL(buf *bytes.Buffer, expr, chkName string, d base) {
	//CONSTRAINT chk_name CHECK (id>0 AND username='admin')
	buf.WriteString(" CONSTRAINT ")

	d.quote(buf, chkName)

	buf.WriteString(" CHECK(")
	buf.WriteString(expr)
	buf.WriteByte(')')
}

// 添加标准的索引约束：pk,unique,foreign key
// 一些非标准的索引需要各个Dialect自己去实现：如mysql的KEY索引
func addIndexes(db core.DB, model *core.Model, d base) error {
	// ALTER TABLE语句的公共语句部分，可以重复利用：
	// ALTER TABLE table_name ADD CONSTRAINT
	buf := bytes.NewBufferString("ALTER TABLE ")
	d.quote(buf, model.Name)
	buf.WriteString(" ADD CONSTRAINT ")
	size := buf.Len()

	// ALTER TABLE tbname ADD CONSTRAINT pk PRIMARY KEY
	buf.WriteString("pk PRIMARY KEY(")
	for _, col := range model.PK {
		buf.WriteString(col.Name)
		buf.WriteByte(',')
	}
	buf.UnreadByte()
	buf.WriteByte(')')
	if _, err := db.Exec(buf.String()); err != nil {
		return err
	}

	// ALTER TABLE tbname ADD CONSTRAINT uniquteName unique(...)
	for name, cols := range model.UniqueIndexes {
		buf.Truncate(size)
		d.quote(buf, name)
		buf.WriteString(" UNIQUE(")
		for _, col := range cols {
			buf.WriteString(col.Name)
			buf.WriteByte(',')
		}
		buf.UnreadByte()
		buf.WriteByte(')')

		if _, err := db.Exec(buf.String()); err != nil {
			return err
		}
	}

	// fk ALTER TABLE tbname ADD CONSTRAINT fkname FOREIGN KEY (col) REFERENCES tbl(tblcol)
	for name, fk := range model.FK {
		buf.Truncate(size)
		d.quote(buf, name)
		buf.WriteString(" FOREIGN KEY(")
		d.quote(buf, fk.Col.Name)
		buf.WriteByte(')')

		buf.WriteString(" REFERENCES ")
		buf.WriteString(fk.RefTableName)
		buf.WriteByte('(')
		buf.WriteString(fk.RefColName)
		buf.WriteByte(')')

		if _, err := db.Exec(buf.String()); err != nil {
			return err
		}
	}

	// chk ALTER TABLE tblname ADD CONSTRAINT chkName CHECK (id>0 AND city='abc')
	for name, expr := range model.Check {
		buf.Truncate(size)
		d.quote(buf, name) // checkName
		buf.WriteString(" CHECK(")
		buf.WriteString(expr)
		buf.WriteByte(')')

		if _, err := db.Exec(buf.String()); err != nil {
			return err
		}
	}

	return nil
}
