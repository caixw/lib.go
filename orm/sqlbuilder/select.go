// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sqlbuilder

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/caixw/lib.go/orm/fetch"
)

var joinType = []string{" LEFT JOIN ", " RIGHT JOIN ", " INNER JOIN ", " FULL JOIN "}

// join功能
// typ join的类型： 0 LEFT JOIN; 1 RIGHT JOIN; 2 INNER JOIN; 3 FULL JOIN
func (s *SQL) joinOn(typ int, table string, on string) *SQL {
	s.join.WriteString(joinType[typ])
	s.join.WriteString(table)
	s.join.WriteString(" ON ")
	s.join.WriteString(on)

	return s
}

// left join ... on ...
func (s *SQL) LeftJoin(table, on string) *SQL {
	return s.joinOn(0, table, on)
}

// right join ... on ...
func (s *SQL) RightJoin(table, on string) *SQL {
	return s.joinOn(1, table, on)
}

// inner join ... on ...
func (s *SQL) InnerJoin(table, on string) *SQL {
	return s.joinOn(2, table, on)
}

// full join ... on ...
func (s *SQL) FullJoin(table, on string) *SQL {
	return s.joinOn(3, table, on)
}

// Order by 功能
// sort 排序方式: 1=asc,2=desc，其它值无效
func (s *SQL) orderBy(sort int, col string) *SQL {
	if sort != 1 && sort != 2 {
		s.errors = append(s.errors, errors.New("orderBy.sort参数错误，只能是1或是2"))
	}

	if s.order.Len() == 0 {
		s.order.WriteString("ORDER BY ")
	} else {
		s.order.WriteString(", ")
	}

	s.order.WriteString(col)
	switch sort {
	case 1:
		s.order.WriteString("ASC ")
	case 2:
		s.order.WriteString("DESC ")
	default:
		s.errors = append(s.errors, errors.New("无效的order值，只能是asc或是desc"))
	}

	return s
}

// order by asc
func (s *SQL) Asc(cols ...string) *SQL {
	for _, c := range cols {
		s.orderBy(1, c)
	}

	return s
}

// order by desc
func (s *SQL) Desc(cols ...string) *SQL {
	for _, c := range cols {
		s.orderBy(2, c)
	}

	return s
}

// limit ... offset ...
// offset值为0时，相当于limit N的效果。
func (s *SQL) Limit(limit, offset int) *SQL {
	s.limitSQL, s.limitArgs = s.db.Dialect().LimitSQL(limit, offset)
	return s
}

// 分页显示，调用Limit()实现，即Page()与Limit()方法会相互覆盖。
func (s *SQL) Page(start, size int) *SQL {
	if start < 1 {
		s.errors = append(s.errors, errors.New("start必须大于0"))
	}
	if size < 1 {
		s.errors = append(s.errors, errors.New("size必须大于0"))
	}

	start-- // 转到从0页开始
	return s.Limit(size, start*size)
}

func (s *SQL) selectSQL() string {
	s.buf.Reset()

	s.buf.WriteString("SELECT ")
	s.buf.WriteString(strings.Join(s.cols, ","))
	s.buf.WriteString(" FROM ")
	s.buf.WriteString(s.tableName)
	s.buf.WriteString(s.join.String())
	s.buf.WriteString(s.cond.String())  // where
	s.buf.WriteString(s.order.String()) // NOTE(caixw):mysql中若要limit，order字段是必须提供的
	s.buf.WriteString(s.limitSQL)

	return s.db.PrepareSQL(s.buf.String())
}

// 导出数据到sql.Rows
func (s *SQL) Query(args ...interface{}) (*sql.Rows, error) {
	if s.HasErrors() {
		return nil, Errors(s.errors)
	}

	if len(args) == 0 {
		// 与selectSQL中添加的顺序相同，where在limit之前
		args = append(s.condArgs, s.limitArgs)
	}

	return s.db.Query(s.selectSQL(), args...)
}

// 导出第一条数据到sql.Row
func (s *SQL) QueryRow(args ...interface{}) *sql.Row {
	if s.HasErrors() {
		panic("构建语句时发生错误信息")
	}

	if len(args) == 0 {
		// 与sqlString中添加的顺序相同，where在limit之前
		args = append(s.condArgs, s.limitArgs)
	}

	return s.db.QueryRow(s.selectSQL(), args...)
}

// 导出数据到map[string]interface{}
func (s *SQL) Fetch2Map(args ...interface{}) (map[string]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := fetch.Map(true, rows)
	if err != nil {
		return nil, err
	}

	return data[0], nil
}

// 导出所有数据到[]map[string]interface{}
func (s *SQL) Fetch2Maps(args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := fetch.Map(true, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 返回指定列的第一行内容
func (s *SQL) FetchColumn(col string, args ...interface{}) (interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := fetch.Column(true, col, rows)
	if err != nil {
		return nil, err
	}

	return data[0], nil
}

// 返回指定列的所有数据
func (s *SQL) FetchColumns(col string, args ...interface{}) ([]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data, err := fetch.Column(false, col, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 将当前select语句查询的数据导出到v中
func (s *SQL) Fetch(v interface{}, args ...interface{}) error {
	rows, err := s.Query(args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Map:
		if vv.Type().Key().Kind() != reflect.String {
			return errors.New("map的键名类型只能为string")
		}
		if vv.Type().Elem().Kind() != reflect.Interface {
			return errors.New("map的键值类型只能为interface{}")
		}

		mapped, err := s.Fetch2Map(args...)
		if err != nil {
			return err
		}
		vv.Set(reflect.ValueOf(mapped))
	case reflect.Slice:
		return fetch.Obj(vv.Interface(), rows)
	case reflect.Struct:
		return fetch.Obj(vv.Interface(), rows)
	default:
		return errors.New("只支持slice,map和struct指针")
	}
	return nil
}
