// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sqlbuilder

import (
	"bytes"
	"errors"
	"strings"
)

const (
	and = iota
	or
)

// 重置所有状态为初始值。
func (w *SQL) resetWhere() {
	w.cond.Reset()
	w.condArgs = w.condArgs[:0]
}

// SQL.And()的别名
func (s *SQL) Where(cond string, args ...interface{}) *SQL {
	return s.And(cond, args...)
}

func (s *SQL) And(cond string, args ...interface{}) *SQL {
	return s.build(and, cond, args...)
}

func (s *SQL) Or(cond string, args ...interface{}) *SQL {
	return s.build(or, cond, args...)
}

func (s *SQL) AndBetween(col string, start, end interface{}) *SQL {
	return s.between(and, col, start, end)
}

func (s *SQL) OrBetween(col string, start, end interface{}) *SQL {
	return s.between(or, col, start, end)
}

// SQL.AndBetween()的别名
func (s *SQL) Between(col string, start, end interface{}) *SQL {
	return s.AndBetween(col, start, end)
}

func (s *SQL) AndIn(col string, args ...interface{}) *SQL {
	return s.in(and, col, args...)
}

func (s *SQL) OrIn(col string, args ...interface{}) *SQL {
	return s.in(or, col, args...)
}

// SQL.AndIn()的别名
func (s *SQL) In(col string, args ...interface{}) *SQL {
	return s.AndIn(col, args...)
}

func (s *SQL) AndIsNull(col string) *SQL {
	return s.isNull(and, col)
}
func (s *SQL) OrIsNull(col string) *SQL {
	return s.isNull(or, col)
}

// SQL.AndIsNull()的别名
func (s *SQL) IsNull(col string) *SQL {
	return s.AndIsNull(col)
}

func (s *SQL) AndIsNotNull(col string) *SQL {
	return s.isNotNull(and, col)
}

func (s *SQL) OrIsNotNull(col string) *SQL {
	return s.isNotNull(or, col)
}

// SQL.AndIsNotNull()的别名
func (s *SQL) IsNotNull(col string) *SQL {
	return s.AndIsNull(col)
}

// 所有SQL子句的构建，最终都调用此方法来写入实例中。
// op 与前一个语句的连接符号，可以是and或是or常量；
// cond 条件语句，值只能是占位符，不能直接写值；
// condArgs 占位符对应的值。
//  w := newSQL(...)
//  w.build(and, "username=='abc'") // 错误：不能使用abc，只能使用？占位符。
//  w.build(and, "username=?", "abc") // 正确，将转换成: and username='abc'
func (s *SQL) build(op int, cond string, args ...interface{}) *SQL {
	switch {
	case s.cond.Len() == 0:
		s.cond.WriteString(" WHERE(")
	case op == and:
		s.cond.WriteString(" AND(")
	case op == or:
		s.cond.WriteString(" OR(")
	default:
		s.errors = append(s.errors, errors.New("无效的op操作符"))
	}

	s.cond.WriteString(cond)
	s.cond.WriteByte(')')

	s.condArgs = append(s.condArgs, args...)

	return s
}

// SQL col in(v1,v2)语句的实现函数，供andIn()和orIn()函数调用。
func (s *SQL) in(op int, col string, args ...interface{}) *SQL {
	if len(args) <= 0 {
		s.errors = append(s.errors, errors.New("condArgs参数不能为空"))
		return s
	}

	cond := bytes.NewBufferString(col)
	cond.WriteString(" IN(")
	cond.WriteString(strings.Repeat("?,", len(s.condArgs)))
	cond.Truncate(cond.Len() - 1) // 去掉最后的逗号
	cond.WriteByte(')')

	return s.build(op, cond.String(), s.condArgs...)
}

// SQL col between start and end 语句的实现函数，供andBetween()和orBetween()调用。
func (s *SQL) between(op int, col string, start, end interface{}) *SQL {
	return s.build(op, col+" BETWEEN ? AND ?", start, end)
}

// SQL col is null 语句的实现函数，供andIsNull()和orIsNull()调用。
func (w *SQL) isNull(op int, col string) *SQL {
	return w.build(op, col+" IS NULL")
}

// SQL col is not null 语句的实现函数，供andIsNotNull()和orIsNotNull()调用。
func (w *SQL) isNotNull(op int, col string) *SQL {
	return w.build(op, col+" IS NOT NULL")
}
