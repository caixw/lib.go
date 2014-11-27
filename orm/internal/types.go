// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package internal

import (
	"bytes"

	"github.com/caixw/lib.go/orm/core"
)

// 对core.Dialect接口的扩展，实现一些通用的接口
type base interface {
	core.Dialect

	// 将sql用core.Dialect.QuoteStr()的字符串引用，并写入buf
	quote(buf *bytes.Buffer, sql string)

	// 将col转换成sql类型，并写入buf中。
	sqlType(buf *bytes.Buffer, col *core.Column)
}
