// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package orm

import (
	"github.com/caixw/lib.go/orm/core"
)

// 判断接口继承
var _ core.DB = &Engine{}
var _ core.DB = &Tx{}
