// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package util

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"unicode"

	"github.com/caixw/lib.go/conv"
	"github.com/caixw/lib.go/encoding/tag"
)

// 将v转换成map[string]reflect.Value形式，其中键名为对象的字段名，
// 键值为字段的值。支持匿名字段，不会转换不可导出(小写字母开头)的
// 字段，也不会转换struct tag以-开头的字段。
func parseObj(v reflect.Value, ret *map[string]reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("v参数的类型只能是reflect.Struct或是struct的指针,当前为[%v]", v.Kind())
	}

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
			continue
		}

		if tagTxt[0] == '-' { // 该字段被标记为忽略
			continue
		}

		name, found := tag.Get(tagTxt, "name")
		if !found { // 没有指定name属性，继续使用字段名
			if unicode.IsUpper(rune(field.Name[0])) {
				(*ret)[field.Name] = v.Field(i)
			}
			continue
		}
		(*ret)[name[0]] = v.Field(i)
	} // end for
	return nil
}

// 将rows中的一条记录写入到val中，必须保证val的类型为reflect.Struct。
// 仅供Fetch2Objs调用。
func fetchOnceObj(val reflect.Value, rows *sql.Rows) error {
	mapped, err := Fetch2Maps(true, rows)
	if err != nil {
		return err
	}

	objItem := make(map[string]reflect.Value, 0)
	if err = parseObj(val, &objItem); err != nil {
		return err
	}

	for index, item := range objItem {
		v, found := mapped[0][index]
		if !found {
			continue
		}
		if err = conv.To(v, item); err != nil {
			return err
		}
	}

	return nil
}

// 将rows中的记录按obj的长度数量导出到obj中。
// val的类型必须是reflect.Slice或是reflect.Array.
func fetchObjToFixedSlice(val reflect.Value, rows *sql.Rows) error {
	itemType := val.Type().Elem()
	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}
	// 判断数组元素的类型是否为struct
	if itemType.Kind() != reflect.Struct {
		return fmt.Errorf("元素类型只能为reflect.Struct或是struct指针，当前为[%v]", itemType.Kind())
	}

	// 先导出数据到map中
	mapped, err := Fetch2Maps(false, rows)
	if err != nil {
		return err
	}

	l := len(mapped)
	if l > val.Len() {
		l = val.Len()
	}

	for i := 0; i < l; i++ {
		objItem := make(map[string]reflect.Value, 0)
		if err = parseObj(val.Index(i), &objItem); err != nil {
			return err
		}
		for index, item := range objItem {
			v, found := mapped[i][index]
			if !found {
				continue
			}
			if err = conv.To(v, item); err != nil {
				return err
			}
		} // end for objItem
	}

	return nil
}

// 将rows中的所有记录导出到val中，val必须为slice的指针。
// 若val的长度不够，会根据rows中的长度调整。
func fetchObjToSlice(val reflect.Value, rows *sql.Rows) error {
	elem := val.Elem()

	itemType := elem.Type().Elem()
	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}
	// 判断数组元素的类型是否为struct
	if itemType.Kind() != reflect.Struct {
		return fmt.Errorf("元素类型只能为reflect.Struct或是struct指针，当前为[%v]", itemType.Kind())
	}

	// 先导出数据到map中
	mapped, err := Fetch2Maps(false, rows)
	if err != nil {
		return err
	}

	// 使elem表示的数组长度最起码和mapped一样。
	size := len(mapped) - elem.Len()
	if size > 0 {
		for i := 0; i < size; i++ {
			elem = reflect.Append(elem, reflect.New(itemType))
		}
		val.Elem().Set(elem)
	}

	for i := 0; i < len(mapped); i++ {
		objItem := make(map[string]reflect.Value, 0)
		if err = parseObj(elem.Index(i), &objItem); err != nil {
			return err
		}

		for index, item := range objItem {
			e, found := mapped[i][index]
			if !found {
				continue
			}
			if err = conv.To(e, item); err != nil {
				return err
			}
		} // end for objItem
	}

	return nil
}

// 将rows中的数据导出到obj中。obj只有在类型为slice指针时，才有可能
// 随着rows的长度变化，否则其长度是固定的，具体可以为以下四种类型：
//
// struct指针：
// 将rows中的第一条记录转换成obj对象。
//
// struct array指针或是struct slice:
// 将rows中的len(obj)条记录导出到obj对象中；若rows中的数量不足，
// 则obj尾部的元素保存原来的值。
//
// struct slice指针：
// 将rows中的所有记录依次写入obj中。若rows中的记录比len(obj)要长，
// 则会增长obj的长度以适应rows的所有记录。
//
// struct可以在struct tag中用name指定字段名称，或是以减号(-)开头
// 表示忽略该字段的导出：
//  type user struct {
//      ID    int `orm:"name(id)"`  // 对应rows中的id字段，而不是ID。
//      age   int `orm:"name(Age)"` // 小写不会被导出。
//      Count int `orm:"-"`         // 不会匹配与该字段对应的列。
//  }
func Fetch2Objs(obj interface{}, rows *sql.Rows) (err error) {
	val := reflect.ValueOf(obj)

	switch val.Kind() {
	case reflect.Ptr:
		elem := val.Elem()
		switch elem.Kind() {
		case reflect.Slice: // slice指针，可以增长
			return fetchObjToSlice(val, rows)
		case reflect.Array: // 数组指针，只能按其大小导出
			return fetchObjToFixedSlice(elem, rows)
		case reflect.Struct: // 结构指针，只能导出一个
			return fetchOnceObj(elem, rows)
		default:
			return fmt.Errorf("不允许的数据类型：[%v]", val.Kind())
		}
	case reflect.Slice: // slice只能按其大小导出。
		return fetchObjToFixedSlice(val, rows)
	default:
		return fmt.Errorf("不允许的数据类型：[%v]", val.Kind())
	}
	return nil
}

// 将rows中的所有或一行数据导出到map[string]interface{}中。
// 若once值为true，则只导出第一条数据。
func Fetch2Maps(once bool, rows *sql.Rows) ([]map[string]interface{}, error) {
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

// 将rows中的数据导出到一个map[string]string中。
// 功能上与Fetch2Maps()上一样，但map的键值固定为string。
func Fetch2MapsString(once bool, rows *sql.Rows) (data []map[string]string, err error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	buf := make([]interface{}, len(cols))
	for k, _ := range buf {
		var val string
		buf[k] = &val
	}

	for rows.Next() {
		if err = rows.Scan(buf...); err != nil {
			return nil, err
		}

		line := make(map[string]string, len(cols))
		for i, v := range cols {
			line[v] = *(buf[i].(*string))
		}

		data = append(data, line)

		if once {
			return data, nil
		}
	}
	return data, nil
}

// 导出rows中某列的所有或一行数据。
// once若为true，则只导出第一条数据的指定列。
// colName指定需要导出的列名，若不指定了不存在的名称，返回error；
func FetchColumns(once bool, colName string, rows *sql.Rows) ([]interface{}, error) {
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
		value := reflect.Indirect(reflect.ValueOf(buff[index]))
		data = append(data, value.Interface())
		if once {
			return data, nil
		}
	}

	return data, nil
}

// 导出rows中某列的所有或是一行数据。
// 除了返回的为[]string以外，其它功能同FetchColumns()。
func FetchColumnsString(once bool, colName string, rows *sql.Rows) ([]string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	index := -1 // 该列名在所rows.Columns()中的索引号
	buff := make([]interface{}, len(cols))
	for i, v := range cols {
		var value string
		buff[i] = &value

		if colName == v { // 获取index的值
			index = i
		}
		// TODO(caixw) 用不到的列，直接赋值为nil，性能上会不会有所提升?
	}

	if index == -1 {
		return nil, fmt.Errorf("指定的名[%v]不存在", colName)
	}

	var data []string
	for rows.Next() {
		if err := rows.Scan(buff...); err != nil {
			return nil, err
		}
		data = append(data, *(buff[index].(*string)))
		if once {
			return data, nil
		}
	}

	return data, nil
}
