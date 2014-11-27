// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// NOTE:测试需要用到mysql数据库中mysql.user表。

package core

import (
	"testing"

	"github.com/caixw/lib.go/assert"
	_ "github.com/go-sql-driver/mysql"
)

func TestStmtsAddSet(t *testing.T) {
	a := assert.New(t)

	db := &fakeDB{}

	s := NewStmts(db)
	a.NotNil(s)

	sql := "SELECT * FROM user WHERE 1"
	sqlStmt, err := db.Prepare(sql)
	a.NotError(err).NotNil(sqlStmt)

	a.True(s.Add("sql1", sqlStmt)).Equal(1, len(s.items))

	// 添加相同名称的sql
	a.False(s.Add("sql1", sqlStmt)).Equal(1, len(s.items))

	// 添加一个不同名称的sql
	a.True(s.Add("sql2", sqlStmt)).Equal(2, len(s.items))

	// 添加一个新的sql
	s.Set("sql3", sqlStmt)
	a.Equal(3, len(s.items))

	// 修改同名的sql
	s.Set("sql1", sqlStmt)
	a.Equal(3, len(s.items))

	// 查找存在的stmt
	stmt, found := s.Get("sql1")
	a.True(found).NotNil(stmt)

	// 查找不存在的stmt
	stmt, found = s.Get("sql4")
	a.False(found).Nil(stmt)

	// 释放所有的缓存
	s.Clear()
	a.Empty(s.items)
	stmt, found = s.Get("sql1")
	a.False(found).Nil(stmt)

	// 释放缓存之后，再次填充
	a.True(s.Add("sql4", sqlStmt)).Equal(1, len(s.items))

	// 查找已被释放的缓存
	stmt, found = s.Get("sql1")
	a.False(found).Nil(stmt)

	// 查找存在的缓存
	stmt, found = s.Get("sql4")
	a.True(found).NotNil(stmt)
}

func TestStmtsAddSetSQL(t *testing.T) {
	a := assert.New(t)

	s := NewStmts(&fakeDB{})
	a.NotNil(s)

	sql := "SELECT * FROM user WHERE 1"
	stmt, err := s.AddSQL("sql1", sql)
	a.NotError(err).
		NotNil(stmt).
		Equal(1, len(s.items))

	// 添加相同名称的sql
	stmt, err = s.AddSQL("sql1", sql)
	a.Error(err).
		Nil(stmt).
		Equal(1, len(s.items))

	// 添加一个不同名称的sql
	stmt, err = s.AddSQL("sql2", sql)
	a.NotError(err).
		NotNil(stmt).
		Equal(2, len(s.items))

	// 添加一个新的sql
	stmt, err = s.SetSQL("sql3", sql)
	a.NotError(err).
		NotNil(stmt).
		Equal(3, len(s.items))

	// 修改同名的sql
	stmt, err = s.SetSQL("sql1", sql)
	a.NotError(err).
		NotNil(stmt).
		Equal(3, len(s.items))

	// 查找存在的stmt
	stmt, found := s.Get("sql1")
	a.True(found).NotNil(stmt)

	// 查找不存在的stmt
	stmt, found = s.Get("sql4")
	a.False(found).Nil(stmt)

	// 释放所有的缓存
	s.Clear()
	a.Empty(s.items)
	stmt, found = s.Get("sql1")
	a.False(found).Nil(stmt)

	// 释放缓存之后，再次填充
	stmt, err = s.AddSQL("sql4", sql)
	a.NotError(err).
		NotNil(stmt).
		Equal(1, len(s.items))

	// 查找已被释放的缓存
	stmt, found = s.Get("sql1")
	a.False(found).Nil(stmt)

	// 查找存在的缓存
	stmt, found = s.Get("sql4")
	a.True(found).NotNil(stmt)

	s.Close()
	a.Nil(s.items)
}
