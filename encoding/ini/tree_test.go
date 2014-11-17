// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ini

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/caixw/lib.go/assert"
)

type section1 struct { // 测试key2是否被正常忽略
	Key1 string `ini:"name(key1)"`
}

type section2 struct {
	Key1 string // 不使用struct tag
	Key2 int    `ini:"name(key2);set(CustomSetFunc);get(CustomGetFunc);" json:"abc"` // 自定义的转换函数
}

func (s *section2) CustomSetFunc(val string, v reflect.Value) error {
	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return err
	}

	num++
	v.SetInt(num)
	return nil
}

func (s *section2) CustomGetFunc(v reflect.Value) (string, error) {
	num := v.Int()
	num--
	return strconv.Itoa(int(num)), nil
}

type root struct {
	Key      string    `ini:"name(key)"`
	Section2 *section2 `ini:"name(section2)"`
	Section1 *section1 `ini:"name(section1)"`
}

func TestUnmarshal(t *testing.T) {
	str := `
    key=val
    [section1]
    key1=val1
    ;comment1  
    key2=val2
    [section2]
    key2=2
    Key1 =    val1
    `

	a := assert.New(t)

	v := &root{Section2: &section2{}, Section1: &section1{}}
	err := Unmarshal([]byte(str), v)
	a.NotError(err)
	a.Equal(v.Key, "val")
	a.Equal(v.Section1.Key1, "val1")
	a.Equal(v.Section2.Key1, "val1")
	a.Equal(v.Section2.Key2, 3)
}

func TestMarshal(t *testing.T) {
	obj := &root{
		Key: "val",
		Section1: &section1{
			Key1: "val1",
		},
		Section2: &section2{
			Key1: "val1",
			Key2: 3,
		},
	}

	str := `
    key=val
    [section1]
    key1=val1
    [section2]
    Key1=val1
    key2=2
    `

	// 分析obj的内容到bytes
	a := assert.New(t)
	val, err := Marshal(obj)
	a.NotError(err)
	a.NotNil(val)

	// 将分析出来的obj重新实例成root对象
	obj1 := &root{Section2: &section2{}, Section1: &section1{}}
	err = Unmarshal(val, obj1)
	a.NotError(err)

	// 将str实例成root对象
	obj2 := &root{Section2: &section2{}, Section1: &section1{}}
	err = Unmarshal([]byte(str), obj2)
	a.NotError(err)

	// 判断两个对象的值是否相等。
	a.Equal(obj1.Key, obj2.Key)
	a.Equal(obj1.Section1.Key1, obj2.Section1.Key1)
	a.Equal(obj1.Section2.Key2, obj2.Section2.Key2)
}
