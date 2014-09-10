// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ini

import (
	"testing"

	"github.com/caixw/lib.go/assert"
)

func TestReader(t *testing.T) {
	a := assert.New(t)
	str := `
    [section1]
    key =    val
    ;comment1  
    ### comment2
    key2=val2
    `
	sectionVals := map[string]string{
		"section1": "",
	}
	commentVals := map[string]string{
		"## comment2\n": "",
		"comment1  \n":  "",
	}
	elementVals := map[string]string{
		"key":  "val",
		"key2": "val2",
	}

	r := NewReaderString(str)
	a.NotNil(r)

LOOP:
	for {
		token, err := r.Token()
		a.NotError(err)
		switch token.Type {
		case EOF:
			break LOOP
		case Comment:
			val := token.Value
			_, found := commentVals[val]
			a.True(found, "实际值为：[%v]", val)
		case Element:
			k, v := token.Key, token.Value
			vv, found := elementVals[k]
			a.True(found)
			a.Equal(vv, v)
		case Section:
			val := token.Value
			_, found := sectionVals[val]
			a.True(found)
		default:
			t.Error("未知的类型")
			break LOOP
		}
	}
}

func TestUnmarshalMap(t *testing.T) {
	str := `
    nosectionkey=nosectionval
    [section]
    skey=sval
    [section1]
    key =    val
    ;comment1  
    ### comment2
    key2=val2
    `
	a := assert.New(t)

	// 不带section参数
	v1 := map[string]interface{}{
		"nosectionkey": "nosectionval",
		"section": map[string]interface{}{
			"skey": "sval",
		},
		"section1": map[string]interface{}{
			"key":  "val",
			"key2": "val2",
		},
	}
	m, err := UnmarshalMap([]byte(str), "")
	a.NotError(err)
	a.Equal(m, v1)

	// 带section参数
	v2 := map[string]interface{}{
		"skey": "sval",
	}
	m, err = UnmarshalMap([]byte(str), "section")
	a.NotError(err)
	a.Equal(m, v2)
}

func TestUnmarshal(t *testing.T) {
	str := `
    nosectionkey=nosectionval
    [section1]
    skey=sval
    key =    val
    ;comment1  
    key2=val2
    [section2]
    key2=val2
    `
	a := assert.New(t)

	type section struct {
		Key2 string `ini:"name:key2";json:"abc"`
	}

	type s1 struct {
		Key string `ini:"name:nosectionkey"`
		//Section1 map[string]interface{} `ini:"Name:section1"`
		Section2 *section `ini:"name:section2"`
	}

	sv := &s1{Section2: &section{}}
	err := Unmarshal([]byte(str), sv)
	a.NotError(err)
	a.Equal(sv.Key, "nosectionval")
	//a.Equal(sv.Section1["key"], "val")
	a.Equal(sv.Section2.Key2, "val2")
}
