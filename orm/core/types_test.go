// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 声明了一些测试用的虚假类：
// - fakeDB实现了DB接口的类，内部调用sqlite3实现。
// - fake1 fakeDriver1注册的数据库实例，与fakeDialect1组成一对。
// - fake2 fakeDriver2注册的数据库实例，与fakeDialect2组成一对。

package core

import (
	"database/sql"
	"database/sql/driver"
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

type dialectBase struct{}

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

// fakeDriver2 对应fakeDialect2
type fakeDriver1 struct{}

func (f *fakeDriver1) Open(arg string) (driver.Conn, error) {
	return nil, nil
}

func init() {
	sql.Register("fake1", &fakeDriver1{})
}

// fakeDriver2 对应fakeDialect2
type fakeDriver2 struct{}

func (f *fakeDriver2) Open(arg string) (driver.Conn, error) {
	return nil, nil
}

func init() {
	sql.Register("fake2", &fakeDriver2{})
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
