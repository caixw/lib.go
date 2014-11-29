// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

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

// fakeDialect1
type fakeDialect1 struct {
}

var _ Dialect = &fakeDialect1{}

func (t *fakeDialect1) GetDBName(dataSource string) string {
	return ""
}

func (t *fakeDialect1) QuoteStr() (string, string) {
	return "[", "]"
}

func (t *fakeDialect1) CreateTable(db DB, m *Model) error {
	return nil
}

func (t *fakeDialect1) LimitSQL(limit, offset int) (string, []interface{}) {
	return "", nil
}

func (t *fakeDialect1) SupportLastInsertId() bool {
	return true
}

// fakeDialect2
type fakeDialect2 struct {
	num int
}

var _ Dialect = &fakeDialect2{}

func (t *fakeDialect2) GetDBName(dataSource string) string {
	return ""
}

func (t *fakeDialect2) QuoteStr() (string, string) {
	return "{", "}"
}

func (t *fakeDialect2) CreateTable(db DB, m *Model) error {
	return nil
}

func (t *fakeDialect2) LimitSQL(limit, offset int) (string, []interface{}) {
	return "", nil
}

func (t *fakeDialect2) SupportLastInsertId() bool {
	return true
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
