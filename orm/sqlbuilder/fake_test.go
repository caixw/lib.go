// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sqlbuilder

import (
	"database/sql"
	"os"
	"strings"

	"github.com/caixw/lib.go/orm/core"
	_ "github.com/mattn/go-sqlite3"
)

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
	return "test"
}

// stmts仅用到了Prepare接口函数
func (f *fakeDB) Prepare(str string) (*sql.Stmt, error) {
	return f.db.Prepare(str)
}

func (f *fakeDB) GetStmts() *core.Stmts {
	return nil
}

func (f *fakeDB) PrepareSQL(sql string) string {
	replace := strings.NewReplacer("{", "[", "}", "]", "#", "prefix_")

	return replace.Replace(sql)
}

func (f *fakeDB) Dialect() core.Dialect {
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
