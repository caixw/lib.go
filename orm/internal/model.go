// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package internal

import (
	"errors"
	"reflect"
	"strconv"
	"unicode"

	"github.com/caixw/lib.go/conv"
	"github.com/caixw/lib.go/encoding/tag"
)

// 在struct tag中的表示名称。
const structTag = "orm"

// Model 表示一个对象的数据库模型。
type Model struct {
	Cols        map[string]*Column   // 所有的列
	Index       map[string][]*Column // 索引列
	UniqueIndex map[string][]*Column // 唯一索引列
	PK          []*Column            // 主键
	AI          *AutoIncr            // 自增列
	FK          *ForeignKey          // 外键
}

// 列结构
type Column struct {
	Name     string // 数据库的字段名
	Len1     int
	Len2     int
	Nullable bool         // 是否可以为NULL
	GoType   reflect.Type // Go语言中的数据类型
	DefVal   string       // 默认值
}

// 自增列
type AutoIncr struct {
	Col         *Column
	Start, Step int // 起始和步长
}

// 外键
type ForeignKey struct {
	Col                *Column // 当前表中的字符
	TableName, ColName string  // 关联的表和字符
}

// 从一个obj声明一个Model实例
func NewModel(obj interface{}) (*Model, error) {
	rtype := reflect.TypeOf(obj)
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	if rtype.Kind() != reflect.Struct {
		return nil, errors.New("obj参数只能是struct或是struct指针")
	}

	model := &Model{Cols: make(map[string]*Column)}

	num := rtype.NumField()
	for i := 0; i < num; i++ {
		if err := parseColumn(rtype.Field(i), model); err != nil {
			return nil, err
		}
	}

	return model, nil
}

// 分析一个字段。
func parseColumn(field reflect.StructField, model *Model) error {
	// 直接忽略以小写字母开头的字段
	if unicode.IsLower(rune(field.Name[0])) {
		return nil
	}

	tagTxt := field.Tag.Get(structTag)

	// 没有附加的struct tag，直接取得几个关键信息返回。
	if len(tagTxt) == 0 {
		model.Cols[field.Name] = &Column{GoType: field.Type, Name: field.Name}
		return nil
	}

	// 以-开头，表示忽略此字段。要确保struct tag最少有一个字符，
	// 所以要上面len(tagTxt) == 0的判断之后。
	if tagTxt[0] == '-' {
		return nil
	}

	col := &Column{GoType: field.Type, Name: field.Name}
	tags := tag.Parse(tagTxt)
	for k, v := range tags {
		switch k {
		case "name":
			col.Name = v[0]
		case "index":
			if model.Index == nil {
				model.Index = make(map[string][]*Column)
			}
			model.Index[v[0]] = append(model.Index[v[0]], col)
		case "pk":
			model.PK = append(model.PK, col)
		case "unique":
			if model.UniqueIndex == nil {
				model.UniqueIndex = make(map[string][]*Column)
			}
			model.UniqueIndex[v[0]] = append(model.UniqueIndex[v[0]], col)
		case "fk":
			model.FK = &ForeignKey{Col: col, ColName: v[1], TableName: v[0]}
		case "nullable":
			if len(v) == 0 {
				col.Nullable = true
			} else {
				col.Nullable = conv.MustBool(v[0], false)
			}
		case "ai":
			if err := model.setAI(col, v); err != nil {
				return err
			}
		case "len":
			if err := col.setLen(v); err != nil {
				return err
			}
		case "default":
			// 默认值
		}
	}
	model.Cols[col.Name] = col

	return nil
}

// 修改或是添加model的AI变量。
func (m *Model) setAI(col *Column, vals []string) (err error) {
	m.AI = &AutoIncr{Col: col, Start: 1, Step: 1}

	switch len(vals) {
	case 0:
	case 1:
		m.AI.Start, err = strconv.Atoi(vals[0])
	case 2:
		m.AI.Step, err = strconv.Atoi(vals[1])
	default:
		err = errors.New("ai标签有过多的参数")
	}

	return
}

// 从参数中获取Column的Len1和Len2变量。
func (c *Column) setLen(vals []string) (err error) {
	switch len(vals) {
	case 0:
	case 1:
		c.Len1, err = strconv.Atoi(vals[0])
	case 2:
		c.Len2, err = strconv.Atoi(vals[1])
	default:
		err = errors.New("len标签有过多的参数")
	}

	return
}
