// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

// 各类Session存储系统需要实现的接口。
type Store interface {
	// 返回数据，不存在就创建
	Get(sid string) (sess *Session, err error)

	// 保存数据
	Save(sess *Session) error

	// 删除数据
	Delete(sid string) error

	// 保存Session中的内容，并将使sess对象处于不可用状态
	Release(sess *Session) error

	// 执行一次垃圾回收，时间小于duration都会被回收
	GC(duration int)

	// 释放整个空间
	Free()
}
