// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"testing"

	"github.com/caixw/lib.go/assert"
	_ "github.com/caixw/lib.go/orm/dialect/test"
)

func TestDBReplaceQuote(t *testing.T) {
	a := assert.New(t)

	//CloseAll() // 先关闭，方便后面随便创建
	d, err := newDB("fakedb1", "root:@/", "test_")
	a.NotError(err).NotNil(d)

	fn := func(sql, wont string) {
		str := d.ReplaceQuote(sql)
		a.Equal(str, wont)
	}

	fn(`"abc"`, "[abc]")
	fn(`"abc".id`, "[abc].id")
	fn(`"abc"."id"`, "[abc].[id]")
	fn(`"abc"."id" as "uid","username"`, "[abc].[id] as [uid],[username]")
	fn(`"abc".*`, "[abc].*")

	fn(`"abc"`, "[abc]")
	fn(`"table"."id"`, "[table].[id]")
	fn(`"table"."id" as "uid"`, "[table].[id] as [uid]")

	fn(`WHERE "a"<>NULL`, "WHERE [a]<>NULL")
	fn(`WHERE 5<"a"`, "WHERE 5<[a]")
	fn(`and "a"<>"b"`, "and [a]<>[b]")
	fn(`WHERE "a"=5 and "b" is NULL`, "WHERE [a]=5 and [b] is NULL")
}

func TestDBReplaceTable(t *testing.T) {
	a := assert.New(t)

	//CloseAll() // 先关闭，方便后面随便创建
	d, err := newDB("fakedb1", "root:@/", "test_")
	a.NotError(err).NotNil(d)

	a.Equal(d.ReplacePrefix("table.user.id"), "test_user.id")
	a.Equal(d.ReplacePrefix("user.id"), "user.id")
	a.Equal(d.ReplacePrefix("table_user.id"), "table_user.id")
	a.Equal(d.ReplacePrefix("table_user.table"), "table_user.table")
}
