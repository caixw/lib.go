// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

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

// 将rows中的数据导出到obj中。obj可以是struct或是struct组成的数组。
func Fetch2Objs(obj interface{}, rows *sql.Rows) (err error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	var mapped []map[string]interface{}
	switch val.Kind() {
	case reflect.Slice: // 导出到一组struct对象数组中
		itemType := val.Type().Elem()
		if itemType.Kind() == reflect.Ptr {
			itemType = itemType.Elem()
		}
		// 判断数组元素的类型是否为struct
		if itemType.Kind() != reflect.Struct {
			return fmt.Errorf("元素类型只能为reflect.Struct或是struct指针，当前为[%v]", itemType.Kind())
		}

		// 先导出数据到map中
		mapped, err = Fetch2Maps(false, rows)
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
	case reflect.Struct: // 导出到一个struct对象中
		mapped, err = Fetch2Maps(true, rows)
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

	default:
		return errors.New("只支持Slice和Struct指针")
	}

	return nil
}

// 将rows中的数据导出到map[string]interface{}中。
// 若once值为true，则只导出第一条数据。返回的map长度为1
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
// 功能上与Fetch2Maps()上一样，但map的值固定为string，方便特殊情况下使用。
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

// 导出rows中某列的数据。
// once若为true，则只导出第一条数据的指定列。
// colName 指定需要导出的列名，若不指定了不存在的名称，返回error；
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

// 导出rows中某列的所有数据。功能同FetchColumns()，除了返回的是字符串数组以外。
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
