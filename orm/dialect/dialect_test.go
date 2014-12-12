// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// 声明了一些测试用的虚假类：
// - fakeDB实现了DB接口的类，内部调用sqlite3实现。
// - fake1 fakeDriver1注册的数据库实例，与fakeDialect1组成一对。
// - fake2 fakeDriver2注册的数据库实例，与fakeDialect2组成一对。

package dialect

import (
	"database/sql"
	"database/sql/driver"
	"os"
	"testing"

	"github.com/caixw/lib.go/assert"
	"github.com/caixw/lib.go/orm/core"
)

func TestIsRegistedDriver(t *testing.T) {
	a := assert.New(t)

	a.True(isRegistedDriver("fake1"))
	a.False(isRegistedDriver("abcdeg"))
}

func TestDialects(t *testing.T) {
	a := assert.New(t)

	clear()
	a.Empty(dialects.items)

	err := Register("fake1", &fakeDialect1{})
	a.NotError(err).
		True(IsRegisted("fake1"))

	// 注册一个相同名称的
	err = Register("fake1", &fakeDialect2{})
	a.Error(err)                    // 注册失败
	a.Equal(1, len(dialects.items)) // 数量还是1，注册没有成功

	// 再注册一个名称不相同的
	err = Register("fake2", &fakeDialect2{})
	a.NotError(err)
	a.Equal(2, len(dialects.items))

	// 注册类型相同，但名称不同的实例
	err = Register("fake3", &fakeDialect2{num: 2})
	a.Error(err)                    // 注册失败
	a.Equal(2, len(dialects.items)) // 数量还是2，注册没有成功

	// 清空
	clear()
	a.Empty(dialects.items)
}

type dialectBase struct{}

func (d *dialectBase) GetDBName(dataSource string) string {
	return ""
}
func (d *dialectBase) CreateTable(db core.DB, m *core.Model) error {
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

var _ core.Dialect = &fakeDialect1{}

func (t *fakeDialect1) QuoteStr() (string, string) {
	return "[", "]"
}

// fakeDialect2
type fakeDialect2 struct {
	dialectBase
	num int
}

var _ core.Dialect = &fakeDialect2{}

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
