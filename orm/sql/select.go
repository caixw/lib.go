// Copyright 2014 by caixw, All rights reserved.
// Use of s source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"bytes"
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/caixw/lib.go/orm/internal"
)

// Select
//
//  sql := NewSelect(db)
//
//  sql.Select("username", "password", `"group"`).
//      From("table.users").
//      Where("id=?", 1)
//      And("username=?", 1)
//      Or(`"group"=?`, 1)
//      Fetch(tx)
type Select struct {
	whereExpr

	db        internal.DB
	table     string
	q         *bytes.Buffer
	join      *bytes.Buffer
	order     *bytes.Buffer
	cols      []string
	limitSQL  string
	limitArgs []interface{}
}

var _ SQLStringer = &Select{}
var _ Stmter = &Select{}
var _ Reseter = &Select{}
var _ Fetch = &Select{}

func NewSelect(db internal.DB) *Select {
	return &Select{
		whereExpr: whereExpr{},
		db:        db,
		q:         bytes.NewBufferString(""),
		// join: bytes.NewBufferString(""), 基本用不到的东西，就不要初始化了。
		order: bytes.NewBufferString(""),
		cols:  make([]string, 0),
	}
}

func (s *Select) Reset() {
	s.whereExpr.Reset()

	s.table = ""
	s.q.Reset()
	s.order.Reset()
	s.cols = s.cols[0:0]
	s.limitSQL = ""
	s.limitArgs = nil

	if s.join != nil {
		s.join.Reset()
	}
}

// 指定列名
func (s *Select) Columns(cols ...string) *Select {
	s.cols = append(s.cols, cols...)

	return s
}

// 指定表名
func (s *Select) Table(name string) *Select {
	s.table = s.db.ReplacePrefix(name)

	return s
}

// Select.Table的别名
func (s *Select) From(name string) *Select {
	return s.Table(name)
}

var joinType = []string{" LEFT JOIN ", " RIGHT JOIN ", " INNER JOIN ", " FULL JOIN "}

// join功能
// typ join的类型： 0 LEFT JOIN; 1 RIGHT JOIN; 2 INNER JOIN; 3 FULL JOIN
func (s *Select) joinOn(typ int, table string, on string) *Select {
	s.join.WriteString(joinType[typ])
	s.join.WriteString(s.db.ReplacePrefix(table))
	s.join.WriteString(" ON ")
	s.join.WriteString(on)

	return s
}

func (s *Select) LeftJoin(table, on string) *Select {
	return s.joinOn(0, table, on)
}

func (s *Select) RightJoin(table, on string) *Select {
	return s.joinOn(1, table, on)
}

// Order by 功能
// sort 排序方式: 1=asc,2=desc，其它值无效
func (s *Select) orderBy(sort int, col string) *Select {
	if sort != 1 && sort != 2 {
		panic("sort参数错误，只能是1或是2")
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
	}

	return s
}

func (s *Select) Asc(cols ...string) *Select {
	for _, c := range cols {
		s.orderBy(1, c)
	}

	return s
}

func (s *Select) Desc(cols ...string) *Select {
	for _, c := range cols {
		s.orderBy(2, c)
	}

	return s
}

func (s *Select) Limit(limit, offset int) *Select {
	s.limitSQL, s.limitArgs = s.db.Dialect().Limit(limit, offset)
	return s
}

// 分页显示，调用Limit()实现，即Page()与Limit()方法会相互覆盖。
// start显示的页面，起始页为1，小于1时将触发panic；
// size每页显示的条数，小于1时将触发panic；
func (s *Select) Page(start, size int) *Select {
	if start < 1 {
		panic("start必须大于0")
	}
	if size < 1 {
		panic("size必须大于0")
	}

	start-- // 转到从0页开始
	return s.Limit(size, start*size)
}

func (s *Select) SQLString(rebuild bool) string {
	if rebuild || s.q.Len() == 0 {
		s.q.Reset()

		s.q.WriteString("SELECT ")
		cols := s.db.ReplaceQuote(strings.Join(s.cols, ","))
		s.q.WriteString(cols)
		s.q.WriteString(" FROM ")
		s.q.WriteString(s.table)
		s.q.WriteString(s.join.String())
		s.q.WriteString(s.condString(s.db)) // where
		s.q.WriteString(s.order.String())   // NOTE(caixw):mysql中若要limit，order字段是必须提供的
		s.q.WriteString(s.limitSQL)
	}

	return s.q.String()
}

func (s *Select) Stmt(name string) (*sql.Stmt, error) {
	return s.db.AddSQLStmt(name, s.q.String())
}

func (s *Select) Query(args ...interface{}) (*sql.Rows, error) {
	if len(args) == 0 {
		// 与SQLString中添加的顺序相同，where在limit之前
		args = append(s.condArgs, s.limitArgs)
	}

	return s.db.Query(s.SQLString(false), args...)
}

// implement Fetch.QueryRow()
func (s *Select) QueryRow(args ...interface{}) *sql.Row {
	if len(args) == 0 {
		// 与SQLString中添加的顺序相同，where在limit之前
		args = append(s.condArgs, s.limitArgs)
	}

	return s.db.QueryRow(s.SQLString(false), args...)
}

func (s *Select) Fetch2Map(args ...interface{}) (map[string]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := Fetch2Maps(true, rows)
	if err != nil {
		return nil, err
	}

	return data[0], nil
}

func (s *Select) Fetch2Maps(args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := Fetch2Maps(true, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Select) FetchColumn(col string, args ...interface{}) (interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := FetchColumns(true, col, rows)
	if err != nil {
		return nil, err
	}

	return data[0], nil
}

func (s *Select) FetchColumns(col string, args ...interface{}) ([]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := FetchColumns(true, col, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Select) Fetch(v interface{}, args ...interface{}) error {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Slice:
		itemType := vv.Type().Elem() // 获取元素类型
		if itemType.Kind() != reflect.Map {
			return errors.New("只支持map[string]interface{}为元素的数组")
		}
		if itemType.Key().Kind() != reflect.String {
			return errors.New("map的键名类型只能为string")
		}
		if itemType.Elem().Kind() != reflect.Interface {
			return errors.New("map的键值类型只能为interface{]")
		}

		mapped, err := s.Fetch2Maps(args...)
		if err != nil {
			return err
		}

		vv.Set(reflect.ValueOf(mapped))
	case reflect.Map:
		if vv.Type().Key().Kind() != reflect.String {
			return errors.New("map的键名类型只能为string")
		}
		if vv.Type().Elem().Kind() != reflect.Interface {
			return errors.New("map的键值类型只能为interface{]")
		}

		mapped, err := s.Fetch2Map(args...)
		if err != nil {
			return err
		}
		vv.Set(reflect.ValueOf(mapped))
	case reflect.Struct:
		//
	default:
		return errors.New("只支持slice,map和struct指针")
	}
	return nil
}
