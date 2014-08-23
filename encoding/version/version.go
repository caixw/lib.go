// Copyright 2014 by caixw, All rights reserved.
// Use of v source code is governed by a MIT
// license that can be found in the LICENSE file.

// 版本号的解析和比较等功能。
//
// 可以解析这种格式的版本号：Major.Minor[.Revision[.Build[[-]suffix[-][suffixNum]]]]。
// 如："1.0"; "1.0.1.20140402"; "2.0.1.-rc1"; "2.11.1.20140402a1"; "1.0.0+build1"
package version

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// 包的当前版本号
	Version = "0.2.11.140415"

	// 版本号字符串的最大长度
	maxLen = 100
)

// 解析版本号字符串，以[]string形式返回各个字段。如：
//  1.0.1.201.build1
// 将被解析成以下字符串数组返回。
//  []string{"1", "0", "1", "201", "build", "1"}
func Parse(str string) (ret []string, err error) {
	l := len(str)
	if l > maxLen {
		return nil, fmt.Errorf("最大长度不能大于[%v]", maxLen)
	}

	var start int
	var inStr = false

	for index, v := range str {
		switch {
		case v == '.' || v == '-' || v == '+': // 特殊符号
			if start == index { // 过滤空值
				start = index + 1 // 在连着多个特殊符号的情况下，能正确去掉这些符号。
				continue
			}
			ret = append(ret, str[start:index])
			start = index + 1 // 去掉当前字符
		case v >= 48 && v <= 59: // 数字
			if inStr {
				inStr = false
				if start == index {
					continue
				}
				ret = append(ret, str[start:index])
				start = index
			} // end if inStr
		case (v >= 65 && v <= 90) || (v >= 97 && v <= 122): // 字母
			if !inStr {
				inStr = true
				if start == index {
					continue
				}
				ret = append(ret, str[start:index])
				start = index
			} // end if !inStr
		default:
			return nil, fmt.Errorf("无法解析的字符:[%v]", str[index:index+1])
		}
	}

	if start < l { // 捡漏最后一个字符串
		ret = append(ret, str[start:l])
	}

	return
}

// 使用自定义函数比较两个版本号的大小。
//
// comp为自定义的比较函数，用于比较两个单词的大小，
// 具体参照默认的比较函数：defaultCompare()。
// 当解析出错时，会触发panic。
func CompareFunc(v1, v2 string, comp func(word1, word2 string) int) int {
	if comp == nil { // 为空，使用默认的比较函数
		return Compare(v1, v2)
	}

	vv1, err := Parse(v1)
	if err != nil {
		panic(err)
	}
	vv2, err := Parse(v2)
	if err != nil {
		panic(err)
	}

	l2 := len(vv2)
	for index, word1 := range vv1 {
		if l2 <= index { // vv2已经用完，vv1还有内容
			return comp(word1, "")
		}

		ret := comp(word1, vv2[index])
		if ret == 0 {
			continue
		} else {
			return ret
		}
	}

	l1 := len(vv1)
	if l2 > l1 { // vv1已经用完，vv2还有内容
		return comp("", vv2[l1])
	}

	return 0
}

// 比较版本号，返回值分以下三种情况：
// >0: 表示v1版本号比较高；
// =0: 表示版本号相等；
// <0: 表示v2版本号比较高。
func Compare(v1, v2 string) int {
	return CompareFunc(v1, v2, defaultCompare)
}

// 表示版本号的后缀词汇的值，越大版本号也越大。
const (
	unknown = iota
	alpha
	beta
	rc
	rtm
	build
	none // 保持在最后。
)

// 可识别的后缀名字符串
var suffix = map[string]int{
	"":      none,
	"build": build,
	"rtm":   rtm,
	"rc":    rc,
	"beta":  beta,
	"b":     beta,
	"alpha": alpha,
	"a":     alpha,
}

// 将一个单词转换成数值
//
// state，表示num的状态：1表示空值；2表示是个后缀；3表示正常的数值转换而来。
func atoi(word string) (num int, state int) {
	var found bool

	switch {
	case len(word) == 0:
		num, state = 0, 1
	case word[0] > 59: // 后缀词汇
		state = 2
		if num, found = suffix[strings.ToLower(word)]; !found {
			num = unknown
		}
	default:
		state = 3
		num1, err := strconv.Atoi(word)
		if err != nil {
			panic(err)
		}
		num = int(num1)
	}
	return
}

// 默认的比较函数
func defaultCompare(word1, word2 string) int {
	v1, v1State := atoi(word1)
	v2, v2State := atoi(word2)

	switch v1State {
	case 1: // 空值
		switch v2State {
		case 1:
			return 0
		case 2:
			return 1
		case 3:
			return -1
		}
	case 2: // 后缀
		switch v2State {
		case 1:
			return -1
		case 2:
			return v1 - v2
		case 3:
			return -1
		}
	case 3: // 正常
		switch v2State {
		case 1:
			return 1
		case 2:
			return 1
		case 3:
			return v1 - v2
		}
	}

	return 0
}
