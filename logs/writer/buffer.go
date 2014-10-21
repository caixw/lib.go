// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writer

import (
	"io"
)

// 带缓存功能的io.Writer，只有数量达到size
// 日志才会真正写入到w中。
type Buffer struct {
	size   int
	buffer [][]byte
	w      io.Writer
}

var _ WriterContainer = &Buffer{}

// 新建一个Buffer
// 当size小于1时，相当于其值为1
func NewBuffer(size int, w io.Writer) *Buffer {
	return &Buffer{size: size, w: w}
}

// WriterContainer.AddWriter
func (b *Buffer) AddWriter(w io.Writer) error {
	if ws, ok := b.w.(WriterContainer); ok {
		ws.AddWriter(w)
		return nil
	}

	if ws, ok := w.(WriterContainer); ok {
		ws.AddWriter(b.w)
		b.w = ws
		return nil
	}

	b.w = NewContainer(b.w, w)
	return nil
}

// io.Writer
func (b *Buffer) Write(bs []byte) (int, error) {
	if b.size <= 1 {
		return b.w.Write(bs)
	}

	b.buffer = append(b.buffer, bs)

	if len(b.buffer) < b.size {
		return len(bs), nil
	}

	return b.Flush()
}

// 分发所有的内容。
func (b *Buffer) Flush() (size int, err error) {
	for _, buf := range b.buffer {
		if size, err = b.w.Write(buf); err != nil {
			return
		}
	}

	b.buffer = b.buffer[:0]
	return
}
