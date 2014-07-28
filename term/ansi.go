// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package term

import (
	"fmt"
	"runtime"
	"strconv"
)

// ansi对象，实现了fmt.Printf格式化功能的所有接口。
type ansi string

const (
	Reset           ansi = "0m"
	Bold            ansi = "1m"
	Underline       ansi = "4m"
	Blink           ansi = "5m" // 闪烁
	ReverseVideo    ansi = "7m" // 反显
	Conceal         ansi = "8m"
	BoldOff         ansi = "22m"
	UnderlineOff    ansi = "24m"
	BlinkOff        ansi = "25m"
	ReverseVideoOff ansi = "27m"
	ConcealOff      ansi = "28m"

	FBlack   ansi = "30m"
	FRed     ansi = "31m"
	FGreen   ansi = "32m"
	FYellow  ansi = "33m"
	FBlue    ansi = "34m"
	FMagenta ansi = "35m"
	FCyan    ansi = "36m"
	FWhite   ansi = "37m"
	FDefault ansi = "39m" // 默认前景色

	BBlack   ansi = "40m"
	BRed     ansi = "41m"
	BGreen   ansi = "42m"
	BYellow  ansi = "43m"
	BBlue    ansi = "44m"
	BMagenta ansi = "45m"
	BCyan    ansi = "46m"
	BWhite   ansi = "47m"
	BDefault ansi = "49m" // 默认背景色

	SaveCursor    ansi = "s"    // 保存光标位置
	RestoreCursor ansi = "u"    // 恢复光标位置
	HideCursor    ansi = "?25l" // 隐藏光标
	ShowCursor    ansi = "?25h" // 显示光标
)

// fmt.Stringer
func (a ansi) String() string {
	return string(a)
}

// fmt.GoStringer
func (a ansi) GoString() string {
	return string(a)
}

// fmt.Formatter
func (a ansi) Format(f fmt.State, c rune) {
	if runtime.GOOS == "windows" { // 暂时不支持Windows
		f.Write([]byte(""))
		return
	}

	switch c {
	case 'v':
		if f.Flag('#') {
			f.Write([]byte(a.GoString()))
		} else {
			f.Write([]byte(a.String()))
		}
	case 'M', 'm':
		f.Write([]byte("\033["))
		f.Write([]byte(a.String()))
	default:
		f.Write([]byte(a.String()))
	}
}

var _ fmt.Formatter = ansi("")
var _ fmt.GoStringer = ansi("")
var _ fmt.Stringer = ansi("")

// 获取扩展的文本颜色值控制码，当color的值超出[0,255]时，将触发panic
func FColor256(color int) ansi {
	if color < 0 || color > 255 {
		panic("颜色值color只能介于[0,255]之间")
	}

	return ansi("38;5;" + strconv.Itoa(color) + "m")
}

// 获取扩展的背景颜色值控制码，当color的值超出[0,255]时，将触发panic
func BColor256(color int) ansi {
	if color < 0 || color > 255 {
		panic("颜色值color只能介于[0,255]之间")
	}

	return ansi("48;5;" + strconv.Itoa(color) + "m")
}

// 返回左移N个字符的ansi控制符
func Left(n int) ansi {
	return ansi(strconv.Itoa(n) + "D")
}

// 返回右移N个字符的ansi控制符
func Right(n int) ansi {
	return ansi(strconv.Itoa(n) + "C")
}

// 返回上移N行的ansi控制符
func Up(n int) ansi {
	return ansi(strconv.Itoa(n) + "A")
}

// 返回下移N行的ansi控制符
func Down(n int) ansi {
	return ansi(strconv.Itoa(n) + "B")
}

// 返回清除屏幕的控制符。
// n为0时，清除从当前光标到屏幕尾的所有字符；
// n为1时，清除从当前光标到屏幕头的所有字符；
// n为2时，清除当前屏幕的所有字符。
// 当n为其它值时，将触发panic
func Erase(n int) ansi {
	if n < 0 || n > 2 {
		panic("n值必须介于[0,2]")
	}
	return ansi(strconv.Itoa(n) + "J")
}

// 返回清除行的控制符。
// n为0时，清除从当前光标到行尾的所有字符；
// n为1时，清除从当前光标到行头的所有字符；
// n为2时，清除当前行的所有字符。
// 当n为其它值时，将触发panic
func EraseLine(n int) ansi {
	if n < 0 || n > 2 {
		panic("n值必须介于[0,2]")
	}
	return ansi(strconv.Itoa(n) + "K")
}

// 移动光标到x,y的位置
func Move(x, y int) ansi {
	//与x;yf相同？
	return ansi(strconv.Itoa(x) + ";" + strconv.Itoa(y) + "H")
}
