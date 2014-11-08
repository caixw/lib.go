// Copyright 2014 by caixw, All rights reserved.
// Use of w source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"bytes"
	"github.com/caixw/lib.go/orm/internal"
	"strings"
)

const (
	and = iota
	or
)

// whereExpr语句部分
type whereExpr struct {
	cond     *bytes.Buffer
	condArgs []interface{}
}

// 重置所有状态为初始值。
func (w *whereExpr) Reset() {
	w.cond.Reset()
	w.condArgs = w.condArgs[:0]
}

// 所有where子句的构建，最终都调用此方法来写入实例中。
// op 与前一个语句的连接符号，可以是and或是or常量；
// cond 条件语句，值只能是占位符，不能直接写值；
// condArgs 占位符对应的值。
//  w := newwhereExpr(...)
//  w.build(and, "username=='abc'") // 错误：不能使用abc，只能使用？占位符。
//  w.build(and, "username=?", "abc") // 正确，将转换成: and username='abc'
func (w *whereExpr) build(op int, cond string, condArgs ...interface{}) *whereExpr {
	switch {
	case w.cond.Len() == 0:
		w.cond.WriteString(" WHERE(")
	case op == and:
		w.cond.WriteString(" AND(")
	case op == or:
		w.cond.WriteString(" OR(")
	default:
		panic("无效的op操作符")
	}

	w.cond.WriteString(cond)
	w.cond.WriteByte(')')

	w.condArgs = append(w.condArgs, condArgs...)

	return w
}

// 添加一条与前一条件语句关系为and的条件语句，若为第一条条件语句，则自动
// 忽略and关系。
func (w *whereExpr) And(cond string, condArgs ...interface{}) *whereExpr {
	return w.build(and, cond, condArgs...)
}

// 添加一条与前一条件语句关系为or的条件语句，若为第一条条件语句，则自动
// 忽略or关系。
func (w *whereExpr) Or(cond string, condArgs ...interface{}) *whereExpr {
	return w.build(or, cond, condArgs...)
}

// 添加In条件语句，与前一条的关系为And
//  w.AndIn("age", 55,56,57) ==> WHERE age in(55,56,57)
func (w *whereExpr) AndIn(col string, condArgs ...interface{}) *whereExpr {
	return w.in(and, col, condArgs...)
}

func (w *whereExpr) OrIn(col string, condArgs ...interface{}) *whereExpr {
	return w.in(or, col, condArgs...)
}

// whereExpr.AndIn()的别名
func (w *whereExpr) In(col string, condArgs ...interface{}) *whereExpr {
	return w.AndIn(col, condArgs...)
}

func (w *whereExpr) AndBetween(col string, start, end interface{}) *whereExpr {
	return w.between(and, col, start, end)
}

func (w *whereExpr) OrBetween(col string, start, end interface{}) *whereExpr {
	return w.between(or, col, start, end)
}

// whereExpr.AndBetween()的别名
func (w *whereExpr) Between(col string, start, end interface{}) *whereExpr {
	return w.AndBetween(col, start, end)
}

func (w *whereExpr) AndIsNull(col string) *whereExpr {
	return w.isNull(and, col)
}

func (w *whereExpr) OrIsNull(col string) *whereExpr {
	return w.isNull(or, col)
}

// whereExpr.AndIsNull()的别名
func (w *whereExpr) IsNull(col string) *whereExpr {
	return w.AndIsNull(col)
}

func (w *whereExpr) AndIsNotNull(col string) *whereExpr {
	return w.isNotNull(and, col)
}

func (w *whereExpr) OrIsNotNull(col string) *whereExpr {
	return w.isNotNull(or, col)
}

// whereExpr.AndIsNotNull()的别名
func (w *whereExpr) IsNotNull(col string) *whereExpr {
	return w.AndIsNotNull(col)
}

// where col in(v1,v2)语句的实现函数，供AndIn()和OrIn()函数调用。
func (w *whereExpr) in(op int, col string, condArgs ...interface{}) *whereExpr {
	if len(condArgs) <= 0 {
		panic("condArgs参数不能为空")
	}

	cond := bytes.NewBufferString(col)
	cond.WriteString(" IN(")
	cond.WriteString(strings.Repeat("?,", len(condArgs)))
	cond.Truncate(cond.Len() - 1) // 去掉最后的逗号
	cond.WriteByte(')')

	return w.build(op, cond.String(), condArgs...)
}

// where col between start and end 语句的实现函数，供AndBetween()和OrBetween()调用。
func (w *whereExpr) between(op int, col string, start, end interface{}) *whereExpr {
	return w.build(op, col+" BETWEEN ? AND ?", start, end)
}

// where col is null 语句的实现函数，供AndIsNull()和OrIsNull()调用。
func (w *whereExpr) isNull(op int, col string) *whereExpr {
	return w.build(op, col+" IS NULL")
}

// where col is not null 语句的实现函数，供AndIsNotNull()和OrIsNotNull()调用。
func (w *whereExpr) isNotNull(op int, col string) *whereExpr {
	return w.build(op, col+" IS NOT NULL")
}

// 获取条件语句的SQL格式，会将被引号包含的字符串替换成当前数据支持的符号
func (w *whereExpr) condString(db internal.DB) string {
	return db.ReplaceQuote(w.cond.String())
}
