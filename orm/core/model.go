// Copyright 2014 by cAIxw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/caixw/lib.go/conv"
	"github.com/caixw/lib.go/encoding/tag"
)

type conType int

// 预定的约束类型，方便Model中使用。
const (
	none conType = iota
	index
	unique
	fk
	check
)

func (t conType) String() string {
	switch t {
	case none:
		return "<none>"
	case index:
		return "KEY INDEX"
	case unique:
		return "UNIQUE INDEX"
	case fk:
		return "FOREIGN KEY"
	case check:
		return "CHECK"
	default:
		return "<unknown>"
	}
}

// model缓存
var models = &modelsMap{items: map[reflect.Type]*Model{}}

type modelsMap struct {
	sync.Mutex
	items map[reflect.Type]*Model
}

// go本身不支持struct级别的struct tag，所以想要给一个struct
// 指定struct tag，只能通过一个函数返回一段描述信息。
type Metaer interface {
	// 表级别的数据。如表名，存储引擎等：
	//  "name:tbl;engine:myISAM;charset:utf-8"
	Meta() string
}

// 表示一个数据库的表模型。数据结构从字段和
// 字段的struct tag中分析得出。
type Model struct {
	Name string // 表的名称

	Cols          map[string]*Column     // 所有的列
	KeyIndexes    map[string][]*Column   // 索引列
	UniqueIndexes map[string][]*Column   // 唯一索引列
	FK            map[string]*ForeignKey // 外键
	PK            []*Column              // 主键
	AI            *AutoIncr              // 自增列
	Check         map[string]string      // Check
	Meta          map[string][]string    // 表级别的元数据，如存储引擎，字符集等。

	constraints map[string]conType // 约束名缓存
}

// 外键
type ForeignKey struct {
	Col                      *Column
	RefTableName, RefColName string
	UpdateRule, DeleteRule   string
}

// 自增列
type AutoIncr struct {
	Col         *Column
	Start, Step int // 起始和步长
}

// 列结构
type Column struct {
	model *Model

	Name     string // 数据库的字段名
	Len1     int
	Len2     int
	Nullable bool         // 是否可以为NULL
	GoType   reflect.Type // Go语言中的数据类型

	HasDefault bool
	Default    string // 默认值
}

// 当前列是否为自增列
func (c *Column) IsAI() bool {
	return (c.model.AI != nil) && (c.model.AI.Col == c)
}

// 从参数中获取Column的len1和len2变量。
// len(len1,len2)
func (c *Column) setLen(vals []string) (err error) {
	switch len(vals) {
	case 0:
	case 1:
		c.Len1, err = strconv.Atoi(vals[0])
	case 2:
		c.Len2, err = strconv.Atoi(vals[1])
	default:
		err = fmt.Errorf("[%v]字段的len属性指定了过多的参数:[%v]", c.Name, vals)
	}

	return
}

// 从vals中分析，得出Column.Nullable的值。
// nullable; or nullable(true);
func (c *Column) setNullable(vals []string) (err error) {
	if c.IsAI() {
		return fmt.Errorf("自增列[%v]不能为nullable", c.Name)
	}

	switch {
	case len(vals) == 0:
		c.Nullable = true
	case len(vals) == 1:
		if c.Nullable, err = conv.Bool(vals[0]); err != nil {
			return err
		}
	default:
		return fmt.Errorf("[%v]字段的nullable属性指定了太多的值:[%v]", c.Name, vals)
	}

	return nil
}

// 从一个obj声明一个Model实例。
// obj可以是一个struct实例或是指针。
func NewModel(obj interface{}) (*Model, error) {
	models.Lock()
	defer models.Unlock()

	rval := reflect.ValueOf(obj)
	if rval.Kind() == reflect.Ptr {
		rval = rval.Elem()
	}
	rtype := rval.Type()

	if rtype.Kind() != reflect.Struct {
		return nil, errors.New("obj参数只能是struct或是struct指针")
	}

	// 是否已经缓存的数组
	if m, found := models.items[rtype]; found {
		return m, nil
	}

	m := &Model{
		Cols:          map[string]*Column{},
		KeyIndexes:    map[string][]*Column{},
		UniqueIndexes: map[string][]*Column{},
		Name:          rtype.Name(),
		FK:            map[string]*ForeignKey{},
		Check:         map[string]string{},
		Meta:          map[string][]string{},
		constraints:   map[string]conType{},
	}

	if err := m.parseColumns(rval); err != nil {
		return nil, err
	}

	if err := m.parseMeta(obj); err != nil {
		return nil, err
	}

	models.items[rtype] = m
	return m, nil
}

// 将rval中的结构解析到m中。支持匿名字段
func (m *Model) parseColumns(rval reflect.Value) error {
	rtype := rval.Type()
	num := rtype.NumField()
	for i := 0; i < num; i++ {
		field := rtype.Field(i)

		if field.Anonymous {
			m.parseColumns(rval.Field(i))
			continue
		}

		if err := m.parseColumn(field); err != nil {
			return err
		}
	}

	return nil
}

// 分析一个字段。
func (m *Model) parseColumn(field reflect.StructField) (err error) {
	// 直接忽略以小写字母开头的字段
	if unicode.IsLower(rune(field.Name[0])) {
		return nil
	}

	tagTxt := field.Tag.Get("orm")

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
		case "name": // name(colname)
			if len(v) != 1 {
				return fmt.Errorf("name属性指定了太多的参数：[%v]", v)
			}
			col.Name = v[0]
		case "index":
			err = m.setIndex(col, v)
		case "pk":
			err = m.setPK(col, v)
		case "unique":
			err = m.setUnique(col, v)
		case "nullable":
			err = col.setNullable(v)
		case "ai":
			err = m.setAI(col, v)
		case "len":
			err = col.setLen(v)
		case "fk":
			err = m.setFK(col, v)
		case "default":
			err = m.setDefault(col, v)
		default:
			err = fmt.Errorf("未知的struct tag属性:[%v]", k)
		}

		if err != nil {
			return err
		}
	}
	m.Cols[col.Name] = col

	return nil
}

// 分析struct的meta接口数据。
func (m *Model) parseMeta(obj interface{}) error {
	meta, ok := obj.(Metaer)
	if !ok {
		return nil
	}

	tags := tag.Parse(meta.Meta())
	if len(tags) == 0 {
		return nil
	}

	for k, v := range tags {
		switch k {
		case "name":
			if len(v) != 1 {
				return fmt.Errorf("Meta接口的name属性指定了太多参数：[%v]", v)
			}

			m.Name = v[0]
		case "check":
			if len(v) != 2 {
				return fmt.Errorf("Meta接口的check属性的参数只能为2个，当前值为:[%v]", v)
			}

			if _, found := m.Check[v[0]]; found {
				return fmt.Errorf("已经存在相同名称[%v]的Check约束", v[0])
			}

			if typ := m.hasConstraint(v[0], check); typ != none {
				return fmt.Errorf("已经存在相同的约束名[%v]，位于[%v]中", v[0], typ)
			}

			m.constraints[v[0]] = check
			m.Check[v[0]] = v[1]
		default:
			m.Meta[k] = v
		}
	}

	return nil
}

// 通过vals设置字段的default属性
// default(5)
func (m *Model) setDefault(col *Column, vals []string) error {
	if m.AI.Col == col {
		return errors.New("自增列不能设置默认值")
	}

	if len(vals) != 1 {
		return fmt.Errorf("[%v]字段的default属性指定了太多的参数：[%v]", col.Name, vals)
	}

	col.HasDefault = true
	col.Default = vals[0]

	return nil
}

// 通过vals设置字段的index约束
// index(idx_name)
func (m *Model) setIndex(col *Column, vals []string) error {
	if len(vals) != 1 {
		return fmt.Errorf("[%v]字段的index属性指定了太多的参数:[%v]", col.Name, vals)
	}

	if typ := m.hasConstraint(vals[0], index); typ != none {
		return fmt.Errorf("已经存在相同的约束名[%v]，位于[%v]中", vals[0], typ)
	}

	m.constraints[vals[0]] = index
	m.KeyIndexes[vals[0]] = append(m.KeyIndexes[vals[0]], col)
	return nil
}

// 通过vals设置字段的primark key约束
// pk
func (m *Model) setPK(col *Column, vals []string) error {
	if len(vals) != 0 {
		return fmt.Errorf("[%v]字段的pk属性指定了太多的参数:[%v]", col.Name, vals)
	}

	if m.AI != nil {
		return fmt.Errorf("已经存在自增列，不需要再次指定主键")
	}

	m.PK = append(m.PK, col)
	return nil
}

// 通过vals设置字段的unique约束
// unique(unique_name)
func (m *Model) setUnique(col *Column, vals []string) error {
	if len(vals) != 1 {
		return fmt.Errorf("[%v]字段的unique属性只能带一个参数:[%v]", col.Name, vals)
	}

	if typ := m.hasConstraint(vals[0], unique); typ != none {
		return fmt.Errorf("已经存在相同的约束名[%v]，位于[%v]中", vals[0], typ)
	}

	m.constraints[vals[0]] = unique
	m.UniqueIndexes[vals[0]] = append(m.UniqueIndexes[vals[0]], col)

	return nil
}

// 通过vals设置字段的foregin key约束
// fk(fk_name,refTable,refColName,updateRule,deleteRule)
func (m *Model) setFK(col *Column, vals []string) error {
	if typ := m.hasConstraint(vals[0], fk); typ != none {
		return fmt.Errorf("已经存在相同的约束名[%v]，位于[%v]中", vals[0], typ)
	}

	if len(vals) < 3 {
		return errors.New("fk参数必须大于3个")
	}

	if _, found := m.FK[vals[0]]; found {
		return fmt.Errorf("重复的外键约束名:[%v]", vals[0])
	}

	fkInst := &ForeignKey{
		Col:          col,
		RefTableName: vals[1],
		RefColName:   vals[2],
	}

	if len(vals) > 3 { // 存在updateRule
		fkInst.UpdateRule = vals[3]
	}
	if len(vals) > 4 { // 存在deleteRule
		fkInst.DeleteRule = vals[4]
	}

	m.constraints[vals[0]] = fk
	m.FK[vals[0]] = fkInst
	return nil
}

// 通过vals设置Model的自增列。
// ai(colName,start,step)
func (m *Model) setAI(col *Column, vals []string) (err error) {
	if col.Nullable {
		return fmt.Errorf("自增列[%v]不能为nullable", col.Name)
	}

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
	if err != nil {
		return err
	}

	m.PK = append(m.PK[:0], col) // 去掉其它主键，将自增列设置为主键
	return nil
}

// 是否存在指定名称的约束名，name不区分大小写。
// 若已经存在返回表示该约束类型的常量，否则返回none。
func (m *Model) hasConstraint(name string, except conType) conType {
	// 约束名不区分大小写
	if typ, found := m.constraints[strings.ToLower(name)]; found && typ != except {
		return typ
	}

	return none
}

// 释放所有的Model缓存。
func FreeModels() {
	models.Lock()
	defer models.Unlock()

	models.items = map[reflect.Type]*Model{}
}
