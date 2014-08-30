// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package ini

import (
	"bufio"
	"bytes"
	"io"

	"github.com/caixw/lib.go/conv"
)

// 用于输出ini内容到指定的io.Writer
type Writer struct {
	buf    *bufio.Writer
	line   int
	symbol byte
}

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

// 将map[string]interface{}转换成ini格式字符串
func MarshalMap(m map[string]interface{}) ([]byte, error) {
	buf := bytes.NewBufferString("")
	w := NewWriter(buf)

	err := marshalMap(w, m)
	if err != nil {
		return nil, err
	}

	w.Flush()
	return buf.Bytes(), nil
}

func marshalMap(w *Writer, m map[string]interface{}) error {
	for index, val := range m {
		switch v := val.(type) {
		case map[string]interface{}:
			w.AddSection(index)
			err := marshalMap(w, v)
			if err != nil {
				return err
			}
		case string:
			w.AddElement(index, v)
		default:
			value, err := conv.String(val)
			if err != nil {
				return err
			}
			w.AddElement(index, value)
		}
	}
	return nil
}
