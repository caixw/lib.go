// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package writer

import (
	"io"
)

// io.Writer的容器接口。该接口可以描述了，如何添加新
// 的io.Writer。同时本身又是一个io.Writer，可以通过
// io.Writer.Write()方法实现对内容的分发。
type WriterContainer interface {
	io.Writer

	// 添加io.Writer到当前的实例。
	AddWriter(io.Writer) error
}

// 对WriterContainer的默认实现。
type Container struct {
	ws []io.Writer
}

var _ WriterContainer = &Container{}

func NewContainer(writers ...io.Writer) *Container {
	return &Container{ws: writers}
}

// 当某一项出错时，将直接返回其信息，后续的都将中断。
// TODO(caixw) 保存error信息，将后面的分发完
func (c *Container) Write(bs []byte) (size int, err error) {
	for _, w := range c.ws {
		if size, err = w.Write(bs); err != nil {
			return
		}
	}

	return
}

func (c *Container) AddWriter(w io.Writer) error {
	c.ws = append(c.ws, w)
	return nil
}
