// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package test

import (
	"github.com/caixw/lib.go/orm/core"
)

type base struct {
	/* data */
}

func (t base) GetDBName(dataSource string) string {
	return ""
}

func (t *base) CreateTable(db core.DB, m *core.Model) error {
	return nil
}

func (m *base) SupportLastInsertId() bool {
	return true
}
