// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writer

import (
	"errors"
	"io"
)

type Buffer struct {
	size   int       // 最大的缓存数量
	buffer [][]byte  // 缓存的消息
	w      io.Writer // 输出的writer
}

var _ WriteFlushAdder = &Buffer{}

// 新建一个Buffer
// 当size小于1时，相当于其值为1
func NewBuffer(w io.Writer, size int) *Buffer {
	return &Buffer{size: size, w: w, buffer: make([][]byte, 0, size)}
}

// WriterContainer.AddWriter
func (b *Buffer) Add(w io.Writer) error {
	if b.w == nil {
		b.w = w
		return nil
	}

	if ws, ok := b.w.(WriteAdder); ok {
		ws.Add(w)
		return nil
	}

	if ws, ok := w.(WriteAdder); ok {
		ws.Add(b.w)
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
	if b.w == nil {
		return 0, errors.New("并未指定输出环境，b.w指向空值")
	}

	for _, buf := range b.buffer {
		if size, err = b.w.Write(buf); err != nil {
			return
		}
	}

	b.buffer = b.buffer[:0]
	return
}
