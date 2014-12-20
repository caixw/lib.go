// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"testing"
)

func TestIsEqual(t *testing.T) {
	eq := func(v1, v2 interface{}) {
		if !IsEqual(v1, v2) {
			t.Errorf("eq:[%v]!=[%v]", v1, v2)
		}
	}

	neq := func(v1, v2 interface{}) {
		if IsEqual(v1, v2) {
			t.Errorf("eq:[%v]==[%v]", v1, v2)
		}
	}

	eq([]byte("abc"), "abc")
	eq("abc", []byte("abc"))

	eq([]byte("中文abc"), "中文abc")
	eq("中文abc", []byte("中文abc"))

	eq([]rune("中文abc"), "中文abc")
	eq("中文abc", []rune("中文abc"))

	eq(5, 5.0)
	eq(int8(5), 5)
	eq(5, int8(5))
	eq(float64(5), int8(5))
	eq([]int{1, 2, 3}, []int{1, 2, 3})
	eq([]int{1, 2, 3}, []int8{1, 2, 3})
	eq([]float32{1, 2.0, 3}, []int8{1, 2, 3})
	eq([]float32{1, 2.0, 3}, []float64{1, 2, 3})

	// 比较两个元素类型可相互转换的数组
	eq(
		[][]int{
			[]int{1, 2},
			[]int{3, 4},
		},
		[][]int8{
			[]int8{1, 2},
			[]int8{3, 4},
		},
	)

	// 比较两个元素类型可转换的map
	eq(
		[]map[int]int{
			map[int]int{1: 1, 2: 2},
			map[int]int{3: 3, 4: 4},
		},
		[]map[int]int8{
			map[int]int8{1: 1, 2: 2},
			map[int]int8{3: 3, 4: 4},
		},
	)
	eq(map[string]int{"1": 1, "2": 2}, map[string]int8{"1": 1, "2": 2})

	// 比较两个元素类型可转换的map
	eq(
		map[int]string{
			1: "1",
			2: "2",
		},
		map[int][]byte{
			1: []byte("1"),
			2: []byte("2"),
		},
	)

	// array 对比
	eq([2]int{1, 2}, [2]int{1, 2})
	eq([2]int{9, 3}, [2]int8{9, 3})
	eq([2]int8{1, 4}, [2]int{1, 4})
	eq([2]int{1, 5}, []int8{1, 5})

	neq(map[int]int{1: 1, 2: 2}, map[int8]int{1: 1, 2: 2})
	neq([]int{1, 2, 3}, []int{3, 2, 1})
	neq("5", 5)
	neq(true, "true")
	neq(true, 1)
	neq(true, "1")
}

func TestIsEmpty(t *testing.T) {
	if IsEmpty([]string{""}) {
		t.Error("IsEmpty([]string{\"\"})")
	}

	if !IsEmpty([]string{}) {
		t.Error("IsEmpty([]string{})")
	}

	if !IsEmpty([]int{}) {
		t.Error("IsEmpty([]int{})")
	}

	if !IsEmpty(map[string]int{}) {
		t.Error("IsEmpty(map[string]int{})")
	}

	if !IsEmpty(0) {
		t.Error("IsEmpty(0)")
	}

	if !IsEmpty("") {
		t.Error("IsEmpty(``)")
	}
}

func TestIsNil(t *testing.T) {
	if !IsNil(nil) {
		t.Error("IsNil(nil)")
	}

	var v1 []int
	if !IsNil(v1) {
		t.Error("IsNil(v1)")
	}

	var v2 map[string]string
	if !IsNil(v2) {
		t.Error("IsNil(v2)")
	}
}

func TestHasPanic(t *testing.T) {
	f1 := func() {
		panic("panic")
	}

	if has, _ := HasPanic(f1); !has {
		t.Error("f1未发生panic")
	}

	f2 := func() {
		f1()
	}

	if has, msg := HasPanic(f2); !has {
		t.Error("f2未发生panic")
	} else if msg != "panic" {
		t.Errorf("f2发生了panic，但返回信息不正确，应为[panic]，但其实返回了%v", msg)
	}

	f3 := func() {
		defer func() {
			if msg := recover(); msg != nil {
				t.Logf("TestHasPanic.f3 recover msg:[%v]", msg)
			}
		}()

		f1()
	}

	if has, msg := HasPanic(f3); has {
		t.Error("f3发生了panic，其信息为:[%v]", msg)
	}

	f4 := func() {
		//todo
	}

	if has, msg := HasPanic(f4); has {
		t.Error("f4发生panic，其信息为[%v]", msg)
	}
}

func TestIsContains(t *testing.T) {
	fn := func(result bool, container, item interface{}) {
		if result != IsContains(container, item) {
			t.Errorf("%v == (IsContains%v, %v)出错\n", result, container, item)
		}
	}

	fn(false, nil, nil)

	fn(true, "abc", "a")
	fn(true, "abc", 'a')       // string vs byte
	fn(true, "abc", rune('a')) // string vs rune
	fn(true, "abc", "c")
	fn(true, "abc", "bc")

	fn(true, "中文a", "中")
	fn(true, "中文a", "a")
	fn(true, "中文a", '中')

	fn(true, []int{1, 2, 3}, 1)
	fn(true, []int{1, 2, 3}, int8(3))
	fn(true, []int{1, 2, 4}, []int{1, 2})
	fn(true, []interface{}{[]int{1, 2}, 5, 6}, []int8{1, 2})
	fn(true, []interface{}{[]int{1, 2}, 5, 6}, 5)

	fn(true, map[string]int{"1": 1, "2": 2}, map[string]int8{"1": 1})
	fn(true,
		map[string][]int{
			"1": []int{1, 2, 3},
			"2": []int{4, 5, 6},
		},
		map[string][]int8{
			"1": []int8{1, 2, 3},
			"2": []int8{4, 5, 6},
		},
	)

	fn(false, map[string]int{}, nil)
	fn(false, map[string]int{"1": 1, "2": 2}, map[string]int8{})
	fn(false, []int{1, 2, 3}, nil)
	fn(false, []int{1, 2, 3}, []int8{1, 3})
}

func TestStyleString(t *testing.T) {
	fn := func(style int, wont string) {
		styleStr := styleString(style)
		if styleStr != wont {
			t.Errorf("style:[%v]；期望值:[%v]；实际返回值[%v]", style, wont, styleStr)
		}
	}

	fn(StyleCase, "StyleCase")
	fn(StyleCase|StyleSpace, "StyleSpace | StyleCase")
	fn(StyleStrit|StyleSpace, "StyleStrit")
	fn(StyleTrim|StyleSpace|StyleCase, "StyleTrim | StyleSpace | StyleCase")
	fn(20, "<invalid style:20>")
}

func TestStringIsEqual(t *testing.T) {
	// 判断相等
	eq := func(s1, s2 string, style int) {
		if !StringIsEqual(s1, s2, style) {
			t.Errorf("s1:[%v] != s2[%v], style=[%v]", s1, s2, styleString(style))
		}
	}

	// 判断不相等
	neq := func(s1, s2 string, style int) {
		if StringIsEqual(s1, s2, style) {
			t.Errorf("s1:[%v] != s2[%v], style=[%v]", s1, s2, styleString(style))
		}
	}

	eq("abc", "abc", StyleStrit)
	eq("abc", "Abc", StyleCase)
	eq("abc", " abc", StyleTrim)
	eq("ab c", "ab\tc", StyleSpace)
	eq("ab c", "ab\t \tc", StyleSpace)
	eq("aB c", "ab\tc", StyleSpace|StyleCase)

	neq("Abc", "abc", StyleStrit)
	neq("abc", "aBc", StyleStrit|StyleCase) // StyleStit会忽略其它方式
}
