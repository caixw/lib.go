// Copyright 2014 by caixw, All rights reserved.
// Use of s source code is governed by a MIT
// license that can be found in the LICENSE file.

package sql

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/caixw/lib.go/encoding/tag"
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

// left join ... on ...
func (s *Select) LeftJoin(table, on string) *Select {
	return s.joinOn(0, table, on)
}

// right join ... on ...
func (s *Select) RightJoin(table, on string) *Select {
	return s.joinOn(1, table, on)
}

// inner join ... on ...
func (s *Select) InnerJoin(table, on string) *Select {
	return s.joinOn(2, table, on)
}

// full join ... on ...
func (s *Select) FullJoin(table, on string) *Select {
	return s.joinOn(3, table, on)
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
	default:
		panic("无效的order值，只能是asc或是desc")
	}

	return s
}

// order by asc
func (s *Select) Asc(cols ...string) *Select {
	for _, c := range cols {
		s.orderBy(1, c)
	}

	return s
}

// order by desc
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

// 将当前实例转换成SQL，所有的变量都以？代替。
// rebuild是否重新产生SQL，一般情况下，只有在调用SQLString()方法之后，
// 如果再修改了内容，需要再次调用SQLString()方法，并传递true值。
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

// 将当前语句预编译并缓存到stmts中，方便之后再次使用。
func (s *Select) Stmt(name string) (*sql.Stmt, error) {
	return s.db.AddSQLStmt(name, s.q.String())
}

// implement Fetch.Query()
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

// implement Fetch.Fetch2Map()
func (s *Select) Fetch2Map(args ...interface{}) (map[string]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := fetch2Maps(true, rows)
	if err != nil {
		return nil, err
	}

	return data[0], nil
}

// implement Fetch.Fetch2Maps()
func (s *Select) Fetch2Maps(args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := fetch2Maps(true, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// implement Fetch.FetchColumn()
func (s *Select) FetchColumn(col string, args ...interface{}) (interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := fetchColumns(true, col, rows)
	if err != nil {
		return nil, err
	}

	return data[0], nil
}

// implement Fetch.FetchColumns()
func (s *Select) FetchColumns(col string, args ...interface{}) ([]interface{}, error) {
	rows, err := s.Query(args...)
	if err != nil {
		return nil, err
	}

	data, err := fetchColumns(true, col, rows)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// implement Fetch.Fetch()
func (s *Select) Fetch(v interface{}, args ...interface{}) error {
	rows, err := s.Query(args...)
	if err != nil {
		return err
	}

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
			return errors.New("map的键值类型只能为interface{]")
		}

		mapped, err := s.Fetch2Map(args...)
		if err != nil {
			return err
		}
		vv.Set(reflect.ValueOf(mapped))
	case reflect.Slice:
		return fetch2Objs(vv.Interface(), rows)
	case reflect.Struct:
		return fetch2Objs(vv.Interface(), rows)
	default:
		return errors.New("只支持slice,map和struct指针")
	}
	return nil
}

// 将v分解成map[string]reflect.Value形式，其中键名为对象的字段名，
// 键值为字段的值，支持匿名字段，不会导出不可导出(小写字母开头)的
// 字段，也不会导出struct tag以-开头的字段。
//
// 假定v.Kind()==reflect.Struct，不再进行判断
func parseObj(v reflect.Value, ret *map[string]reflect.Value) {
	vt := v.Type()
	num := vt.NumField()
	for i := 0; i < num; i++ {
		field := vt.Field(i)

		if field.Anonymous { // 匿名对象
			parseObj(v.Field(i), ret)
		}

		tagTxt := field.Tag.Get("orm")
		if len(tagTxt) == 0 { // 不存在Struct tag
			if unicode.IsUpper(rune(field.Name[0])) {
				(*ret)[field.Name] = v.Field(i)
			}
		}

		if tagTxt[0] == '-' { // 该字段被标记为忽略
			continue
		}

		name, found := tag.Get(tagTxt, "name")
		if !found { // 没有指定name属性，继续使用字段名
			if unicode.IsUpper(rune(field.Name[0])) {
				(*ret)[field.Name] = v.Field(i)
			}
		}
		(*ret)[name[0]] = v.Field(i)
	} // end for
}

func fetch2Objs(obj interface{}, rows *sql.Rows) (err error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var mapped []map[string]interface{}
	switch val.Kind() {
	case reflect.Slice: // 导出到一组struct对象数组中
		itemType := val.Type().Elem()
		// 判断数组元素的类型是否为struct
		if itemType.Kind() != reflect.Struct {
			return fmt.Errorf("元素类型只能为reflect.Struct，当前为[%v]", itemType.Kind())
		}

		// 先导出数据到map中
		mapped, err = fetch2Maps(false, rows)
		if err != nil {
			return err
		}

		// 使val表示的数组长度最起码和mapped一样。
		size := len(mapped) - val.Len()
		for i := 0; i < size; i++ {
			val = reflect.Append(val, reflect.New(itemType))
		}

		for i := 0; i < len(mapped); i++ {
			objItem := make(map[string]reflect.Value, 0)
			parseObj(val.Index(i), &objItem)
			for index, item := range objItem {
				if v, found := mapped[i][index]; found {
					item.Set(reflect.ValueOf(v))
				}
			}

		}
	case reflect.Struct: // 导出到一个struct对象中
		mapped, err = fetch2Maps(true, rows)
		if err != nil {
			return err
		}
		objItem := make(map[string]reflect.Value, 0)
		parseObj(val, &objItem)
		for index, item := range objItem {
			if v, found := mapped[0][index]; found {
				item.Set(reflect.ValueOf(v))
			}
		}

	default:
		return errors.New("只支持Slice和Struct指针")
	}

	return nil
}

// 将rows中的数据导出到map中
//
// once 是否只查询一条记录，若为true，则返回长度为1的slice
// rows 查询的结果
// 对外公开，方便db包调用
func fetch2Maps(once bool, rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 临时缓存，用于保存从rows中读取出来的一行。
	buff := make([]interface{}, len(cols))
	for i, _ := range cols {
		var value interface{}
		buff[i] = &value
	}

	var data []map[string]interface{}
	for rows.Next() {
		if err := rows.Scan(buff...); err != nil {
			return nil, err
		}

		line := make(map[string]interface{}, len(cols))
		for i, v := range cols {
			if buff[i] == nil {
				continue
			}
			value := reflect.Indirect(reflect.ValueOf(buff[i]))
			line[v] = value.Interface()
		}

		data = append(data, line)
		if once {
			return data, nil
		}
	}

	return data, nil
}

// 导出某列的所有数据
//
// colName 该列的名称，若不指定了不存在的名称，返回error
func fetchColumns(once bool, colName string, rows *sql.Rows) ([]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	index := -1 // 该列名在所rows.Columns()中的索引号
	buff := make([]interface{}, len(cols))
	for i, v := range cols {
		var value interface{}
		buff[i] = &value

		if colName == v { // 获取index的值
			index = i
		}
	}

	if index == -1 {
		return nil, errors.New("指定的名不存在")
	}

	var data []interface{}
	for rows.Next() {
		if err := rows.Scan(buff...); err != nil {
			return nil, err
		}
		data = append(data, buff[index])
		if once {
			return data, nil
		}
	}

	return data, nil
}
