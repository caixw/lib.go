// Copyright 2014 by cAIxw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unicode"

	"github.com/caixw/lib.go/conv"
	"github.com/caixw/lib.go/encoding/tag"
)

// go本身不支持struct级别的struct tag，所以要给一个struct
// 指定struct tag，只能通过一个函数返回一段描述信息。
type Metaer interface {
	// 表级别的数据。如表名，存储引擎等：
	//  "name:tbl;engine:myISAM;charset:utf-8"
	Meta() string
}

const (
	structTag = "orm" // 在struct tag中的表示名称。
)

// Model 从struct tag中初始化的数据表模型。
type Model struct {
	Name string

	Cols          map[string]*Column     // 所有的列
	KeyIndexes    map[string][]*Column   // 索引列
	UniqueIndexes map[string][]*Column   // 唯一索引列
	FK            map[string]*ForeignKey // 外键
	PK            []*Column              // 主键
	AI            *AutoIncr              // 自增列
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

// 列结构
type Column struct {
	model *Model

	Name     string // 数据库的字段名
	Len1     int
	Len2     int
	Nullable bool         // 是否可以为NULL
	GoType   reflect.Type // Go语言中的数据类型
	DefVal   string       // 默认值
}

// 当前列是否为自增列
func (c *Column) IsAI() bool {
	return (c.model.AI != nil) && (c.model.AI.Col == c)
}

// 从参数中获取Column的len1和len2变量。
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

// 从一个obj声明一个Model实例
func NewModel(obj interface{}) (*Model, error) {
	rtype := reflect.TypeOf(obj)
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	if rtype.Kind() != reflect.Struct {
		return nil, errors.New("obj参数只能是struct或是struct指针")
	}

	m := &Model{
		Cols:          map[string]*Column{},
		KeyIndexes:    map[string][]*Column{},
		UniqueIndexes: map[string][]*Column{},
		FK:            map[string]*ForeignKey{},
		Name:          rtype.Name(),
	}

	num := rtype.NumField()
	for i := 0; i < num; i++ {
		if err := m.parseColumn(rtype.Field(i)); err != nil {
			return nil, err
		}
	}

	// 分析meta数据
	meta, ok := obj.(Metaer)
	if !ok {
		return m, nil
	}
	metas := tag.Parse(meta.Meta())
	if len(metas) == 0 {
		return m, nil
	}
	for k, v := range metas {
		switch k {
		case "name":
			m.Name = v[0]
		default:
		}
	}

	return m, nil
}

// 分析一个字段。
func (m *Model) parseColumn(field reflect.StructField) error {
	// 直接忽略以小写字母开头的字段
	if unicode.IsLower(rune(field.Name[0])) {
		return nil
	}

	tagTxt := field.Tag.Get(structTag)

	// 没有附加的struct tag，直接取得几个关键信息返回。
	if len(tagTxt) == 0 {
		m.Cols[field.Name] = &Column{GoType: field.Type, Name: field.Name, model: m}
		return nil
	}

	// 以-开头，表示忽略此字段。要确保struct tag最少有一个字符，
	// 所以要上面len(tagTxt) == 0的判断之后。
	if tagTxt[0] == '-' {
		return nil
	}

	col := &Column{GoType: field.Type, Name: field.Name, model: m}
	tags := tag.Parse(tagTxt)
	for k, v := range tags {
		switch k {
		case "name":
			col.Name = v[0]
		case "index":
			m.KeyIndexes[v[0]] = append(m.KeyIndexes[v[0]], col)
		case "pk":
			if m.AI != nil { // 若存在自增列，则不理其它主键设置
				continue
			}
			m.PK = append(m.PK, col)
		case "unique":
			m.UniqueIndexes[v[0]] = append(m.UniqueIndexes[v[0]], col)
		case "fk":
			m.FK[v[0]] = &ForeignKey{Col: col, ColName: v[2], TableName: v[1]}
		case "nullable":
			if col.IsAI() {
				panic(fmt.Sprintf("自增列[%v]不能为nullable", col.Name))
			}

			if len(v) == 0 {
				col.Nullable = true
			} else {
				col.Nullable = conv.MustBool(v[0], false)
			}
		case "ai":
			if col.Nullable {
				panic(fmt.Sprintf("自增列[%v]不能为nullable", col.Name))
			}
			if err := m.setAI(col, v); err != nil {
				return err
			}

			m.PK = append(m.PK[:0], col) // 则去掉其它主键，将自增列设置为主键
		case "len":
			if err := col.setLen(v); err != nil {
				return err
			}
		case "default":
			col.DefVal = v[0]
		default:
		}
	}
	m.Cols[col.Name] = col

	return nil
}

// 修改或是添加Model的AI变量。
func (m *Model) setAI(col *Column, vals []string) (err error) {
	switch col.GoType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	default:
		return errors.New("自增列只能是整数类型")
	}

	m.AI = &AutoIncr{Col: col, Start: 1, Step: 1}

	switch len(vals) {
	case 0:
	case 1:
		m.AI.Start, err = strconv.Atoi(vals[0])
	case 2:
		m.AI.Step, err = strconv.Atoi(vals[1])
	default:
		err = errors.New("AI标签有过多的参数")
	}

	return
}
