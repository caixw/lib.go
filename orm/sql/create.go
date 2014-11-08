// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"database/sql"
	"github.com/caixw/lib.go/orm/internal"
)

type Create struct {
	model *internal.Model
	table string
	db    internal.DB
}

var _ SQLStringer = &Create{}
var _ Execer = &Create{}

func NewCreate(db internal.DB) *Create {
	return &Create{db: db}
}

func (c *Create) Model(v interface{}) *Create {
	m, err := internal.NewModel(v)
	if err != nil {
		panic(err)
	}

	c.model = m

	return c
}

func (c *Create) Table(name string) *Create {
	c.table = c.db.ReplacePrefix(name)

	return c
}

func (c *Create) SQLString(rebuild bool) string {
	return ""
}

func (c *Create) Exec(args ...interface{}) (sql.Result, error) {
	return c.db.Exec(c.SQLString(false))
}
