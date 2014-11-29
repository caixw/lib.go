// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 本文件主要包含了types.go文件中声明的两个接口的测试实例：
// Dialect:fakeDialect1,fakeDialect2;
// DB:fakeDB.

package core

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestConTypeString(t *testing.T) {
	a := assert.New(t)

	a.Equal("<none>", none.String()).
		Equal("KEY INDEX", fmt.Sprint(index)).
		Equal("UNIQUE INDEX", unique.String()).
		Equal("FOREIGN KEY", fk.String()).
		Equal("CHECK", check.String())

	var c1 conType
	a.Equal("<none>", c1.String())

	c1 = 100
	a.Equal("<unknown>", c1.String())
}

type dialectBase struct {
	/* data */
}

func (d *dialectBase) GetDBName(dataSource string) string {
	return ""
}
func (d *dialectBase) CreateTable(db DB, m *Model) error {
	return nil
}

func (d *dialectBase) LimitSQL(limit, offset int) (string, []interface{}) {
	return "", nil
}

func (d *dialectBase) SupportLastInsertId() bool {
	return true
}

// fakeDialect1
type fakeDialect1 struct {
	dialectBase
}

var _ Dialect = &fakeDialect1{}

func (t *fakeDialect1) QuoteStr() (string, string) {
	return "[", "]"
}

// fakeDialect2
type fakeDialect2 struct {
	dialectBase
	num int
}

var _ Dialect = &fakeDialect2{}

func (t *fakeDialect2) QuoteStr() (string, string) {
	return "{", "}"
}

// fakeDB
type fakeDB struct {
	db *sql.DB
}

func newFakeDB() (*fakeDB, error) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		return nil, err
	}

	return &fakeDB{
		db: db,
	}, nil
}

func (f *fakeDB) close() {
	f.db.Close()
	os.Remove("./test.db")
}

func (f *fakeDB) Name() string {
	return ""
}

// stmts仅用到了Prepare接口函数
func (f *fakeDB) Prepare(str string) (*sql.Stmt, error) {
	return f.db.Prepare(str)
}

func (f *fakeDB) GetStmts() *Stmts {
	return nil
}

func (f *fakeDB) ReplaceQuote(cols string) string {
	return ""
}

func (f *fakeDB) ReplacePrefix(cols string) string {
	return ""
}

func (f *fakeDB) Dialect() Dialect {
	return nil
}

func (f *fakeDB) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (f *fakeDB) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (f *fakeDB) QueryRow(sql string, args ...interface{}) *sql.Row {
	return nil
}
