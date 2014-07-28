// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package term

import (
	"fmt"
	"io"
)

// ansi控制器的io.Writer接口
//
//  a := NewWriter(os.Stdout)
//  a.Left(5).ClearLine(2).SGR(term.SGRFRed,term.SGRBGreen).Print("abc")
//  fmt.Fprintf(a,
type AnsiWriter struct {
	w io.Writer
}

func NewWriter(w io.Writer) *AnsiWriter {
	return &AnsiWriter{w: w}
}

// implements io.Writer
func (a *AnsiWriter) Write(b []byte) (int, error) {
	return a.w.Write(b)
}

var _ io.Writer = &AnsiWriter{}

// 向io.Writer写入ansi控制器
func (a *AnsiWriter) WriteAnsi(code ansi) *AnsiWriter {
	if _, err := fmt.Fprintf(a.w, "%m", code); err != nil {
		panic(err)
	}
	return a
}

// 左移n个字符光标
func (a *AnsiWriter) Left(n int) *AnsiWriter {
	return a.WriteAnsi(Left(n))
}

// 右移n个字符光标
func (a *AnsiWriter) Right(n int) *AnsiWriter {
	return a.WriteAnsi(Right(n))
}

// 上移n行光标
func (a *AnsiWriter) Up(n int) *AnsiWriter {
	return a.WriteAnsi(Up(n))
}

// 下移n行光标
func (a *AnsiWriter) Down(n int) *AnsiWriter {
	return a.WriteAnsi(Down(n))
}

// 清除屏幕。
// n为0时，清除从当前光标到屏幕尾的所有字符；
// n为1时，清除从当前光标到屏幕头的所有字符；
// n为2时，清除当前屏幕的所有字符。
// 当n为其它值时，将触发panic
func (a *AnsiWriter) Erase(n int) *AnsiWriter {
	return a.WriteAnsi(Erase(n))
}

// 清除行。
// n为0时，清除从当前光标到行尾的所有字符；
// n为1时，清除从当前光标到行头的所有字符；
// n为2时，清除当前行的所有字符。
// 当n为其它值时，将触发panic
func (a *AnsiWriter) EraseLine(n int) *AnsiWriter {
	return a.WriteAnsi(EraseLine(n))
}

// 移动光标到x,y的位置
func (a *AnsiWriter) Move(x, y int) *AnsiWriter {
	return a.WriteAnsi(Move(x, y))
}

func (a *AnsiWriter) SaveCursor() *AnsiWriter {
	return a.WriteAnsi(SaveCursor)
}

func (a *AnsiWriter) RestoreCursor() *AnsiWriter {
	return a.WriteAnsi(RestoreCursor)
}

func (a *AnsiWriter) HideCursor() *AnsiWriter {
	return a.WriteAnsi(HideCursor)
}

func (a *AnsiWriter) ShowCursor() *AnsiWriter {
	return a.WriteAnsi(ShowCursor)
}

func (a *AnsiWriter) SGR(sgr ...string) *AnsiWriter {
	a.WriteAnsi(SGR(sgr...))
	return a
}

func (a *AnsiWriter) FColor256(color int) *AnsiWriter {
	return a.WriteAnsi(FColor256(color))
}

func (a *AnsiWriter) BColor256(color int) *AnsiWriter {
	return a.WriteAnsi(BColor256(color))
}

func (a *AnsiWriter) Color256(f, b int) *AnsiWriter {
	a.WriteAnsi(FColor256(f))
	return a.WriteAnsi(BColor256(b))
}

func (a *AnsiWriter) Printf(format string, args ...interface{}) *AnsiWriter {
	if _, err := fmt.Fprintf(a.w, format, args...); err != nil {
		panic(err)
	}
	return a
}

func (a *AnsiWriter) Print(args ...interface{}) *AnsiWriter {
	if _, err := fmt.Fprint(a.w, args...); err != nil {
		panic(err)
	}

	return a
}

func (a *AnsiWriter) Println(args ...interface{}) *AnsiWriter {
	if _, err := fmt.Fprintln(a.w, args...); err != nil {
		panic(err)
	}
	return a
}
