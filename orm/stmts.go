// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"database/sql"
	"fmt"
	"sync"
)

// 缓存sql.Stmt实例，方便再次调用。
type Stmts struct {
	sync.Mutex
	items map[string]*sql.Stmt
	db    *db
}

func newStmts(db *db) *Stmts {
	return &Stmts{
		items: make(map[string]*sql.Stmt),
		db:    db,
	}
}

func (s *Stmts) AddSql(name, sql string) (*sql.Stmt, error) {
	s.Lock()
	defer s.Unlock()

	if _, found := s.items[name]; found {
		return nil, fmt.Errorf("该名称[%v]的stmt已经存在", name)
	}

	stmt, err := s.db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	s.items[name] = stmt
	return stmt, nil
}

func (s *Stmts) SetSql(name, sql string) (*sql.Stmt, error) {
	stmt, err := s.db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	s.Lock()
	defer s.Unlock()

	s.items[name] = stmt
	return stmt, nil
}

// 添加一个sql.Stmt，若已经存在相同名称的，则返回false。
// 若要修改已存在的sql.Stmt，请使用Set()函数
func (s *Stmts) Add(name string, stmt *sql.Stmt) bool {
	s.Lock()
	defer s.Unlock()

	if _, found := s.items[name]; found {
		return false
	}

	s.items[name] = stmt
	return true
}

// 添加或是修改sql.Stmt
func (s *Stmts) Set(name string, stmt *sql.Stmt) {
	s.Lock()
	defer s.Unlock()

	s.items[name] = stmt
}

// 查找指定名称的sql.Stmt实例。若不存在，返回nil,false
func (s *Stmts) Get(name string) (stmt *sql.Stmt, found bool) {
	stmt, found = s.items[name]
	return
}

// 释放所有缓存空间。
func (s *Stmts) free() {
	s.items = nil
}
