// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package assert

import (
	"reflect"
	"time"
)

// 判断一个值是否为空(0, "", false, 空数组等)。
func IsEmpty(expr interface{}) bool {
	if expr == nil {
		return true
	}

	switch v := expr.(type) {
	case bool:
		return false == v
	case int:
		return 0 == v
	case int8:
		return 0 == v
	case int16:
		return 0 == v
	case int32:
		return 0 == v
	case int64:
		return 0 == v
	case uint:
		return 0 == v
	case uint8:
		return 0 == v
	case uint16:
		return 0 == v
	case uint32:
		return 0 == v
	case uint64:
		return 0 == v
	case string:
		return "" == v
	case time.Time:
		return v.IsZero()
	case *time.Time:
		return v.IsZero()
	}

	// 符合IsNil条件的，都为Empty
	ret := IsNil(expr)
	if ret {
		return true
	}

	v := reflect.ValueOf(expr)
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Chan:
		return 0 == v.Len()
	case reflect.Ptr:
		return false
	}

	return false
}

// 判断一个值是否为nil。
// 当特定类型的变量，已经声明，但还未赋值时，也将返回true
func IsNil(expr interface{}) bool {
	if nil == expr {
		return true
	}

	v := reflect.ValueOf(expr)
	k := v.Kind()

	if (k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.Interface ||
		k == reflect.Map ||
		k == reflect.Ptr ||
		k == reflect.Slice) &&
		v.IsNil() {
		return true
	}

	return false
}

// 判断两个值是否相等。
//
// 除了通过reflect.DeepEqual()判断值是否相等之外，一些类似
// 可转换的数值也能正确判断，比如int(5)与int64(5)能正确判断。
func IsEqual(v1, v2 interface{}) bool {
	if reflect.DeepEqual(v1, v2) {
		return true
	}

	vv1 := reflect.ValueOf(v1)
	vv2 := reflect.ValueOf(v2)

	// NOTE: 这里返回false，而不是true
	if !vv1.IsValid() || !vv2.IsValid() {
		return false
	}

	if vv1 == vv2 {
		return true
	}

	vv1Type := vv1.Type()

	// reflect.DeepEqual已经比较过以下类型的值，若是以下类型，直接返回false。
	// 在reflect.DeepEqual()中比较[]string{"1"},[]string{"2"}会返回false，但是以
	// 下的ConvertableTo中又可以转换， 进而又进行一次比较。所以此处必须过滤掉已经
	// 在reflect.DeepEqual()中已经处理过的值。
	switch vv1Type.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct, reflect.Ptr, reflect.Func, reflect.Interface:
		return false
	}

	vv2Type := vv2.Type()
	if vv1Type.ConvertibleTo(vv2Type) {
		return vv2.Interface() == vv1.Convert(vv2Type).Interface()
	} else if vv2Type.ConvertibleTo(vv1Type) {
		return vv1.Interface() == vv2.Convert(vv1Type).Interface()
	}

	return false
}
