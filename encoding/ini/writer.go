// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// 用于输出ini内容到指定的io.Writer。
//
// 内容并不是实时写入io.Writer的，需要调用Writer.Flush()
// 才会真正地写入到io.Writer流中。
type Writer struct {
	buf    *bufio.Writer
	line   int
	symbol byte
}

// 从一个io.Writer初始化Writer
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		buf:    bufio.NewWriter(w),
		symbol: '#',
	}
}

// 设置注释符号，默认为`#`
//
// 有数据写入之后，不能再更改，否则会触发panic
func (w *Writer) SetCommentSymbol(symbol byte) *Writer {
	if w.line != 0 {
		panic("已经有数据写入，不能再次更改symbol")
	}

	w.symbol = symbol
	return w
}

// 添加一个新的空行。
func (w *Writer) NewLine() *Writer {
	w.buf.WriteByte('\n')
	w.line++
	return w
}

// 添加section，没有嵌套功能，添加一个新的Section，意味着前一个section的结束。
func (w *Writer) AddSection(section string) *Writer {
	w.buf.WriteByte('[')
	w.buf.WriteString(section)
	w.buf.WriteByte(']')

	return w.NewLine()
}

// 添加一个键值对。
func (w *Writer) AddElement(key, val string) *Writer {
	w.buf.WriteString(key)
	w.buf.WriteByte('=')
	w.buf.WriteString(val)

	return w.NewLine()
}

// 添加一个键值对。val使用fmt.Sprint格式化成字符串。
func (w *Writer) AddElementf(key string, val interface{}) *Writer {
	return w.AddElement(key, fmt.Sprint(val))
}

// 添加注释
func (w *Writer) AddComment(comment string) *Writer {
	w.buf.WriteByte(w.symbol)
	w.buf.WriteString(comment)

	return w.NewLine()
}

// 将内容输出到io.Writer中
func (w *Writer) Flush() {
	w.buf.Flush()
}

// 将v的内容以ini格式的形式输出。
func Marshal(v interface{}) ([]byte, error) {
	buf := bytes.NewBufferString("")
	w := NewWriter(buf)
	tree, err := scan(v)
	if err != nil {
		return nil, err
	}

	if err = tree.marshal(w); err != nil {
		return nil, err
	}

	w.Flush()
	return buf.Bytes(), nil
}
