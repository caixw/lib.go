// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"testing"
)

func TestIsEqual(t *testing.T) {
	if IsEqual("5", 5) {
		t.Error(`IsEqual("5"==5)`)
	}

	if !IsEqual(5, 5.0) {
		t.Error(`IsEqual(5!=5.0)`)
	}
}

func TestIsEmpty(t *testing.T) {
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
