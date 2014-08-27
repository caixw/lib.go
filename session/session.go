// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"sync"
	"time"
)

// session可回收利用的对象池
var sessFree = sync.Pool{
	New: func() interface{} { return &Session{} },
}

// 针对Session的相关操作
type Session struct {
	sync.Mutex
	data     map[interface{}]interface{}
	store    Store
	id       string
	accessed time.Time
}

// sessionid
func (s *Session) ID() string {
	return s.id
}

// 获取某个Session的值
func (s *Session) Get(key interface{}) (val interface{}, found bool) {
	s.Lock()
	defer s.Unlock()
	val, found = s.data[key]
	return
}

// 同Get，但是在值不存在时，返回def作为默认值。
func (s *Session) MustGet(key, def interface{}) interface{} {
	s.Lock()
	defer s.Unlock()

	val, found := s.data[key]
	if !found {
		return def
	} else {
		return val
	}
}

// 设置某个Session的值
func (s *Session) Set(key, val interface{}) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = val
}

// 删除某个Session的值
func (s *Session) Delete(key interface{}) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
}

// 保存并删除当前Session的内容。
// 必须释放Session，否则可能造成内存泄漏。
func (s *Session) Release() error {
	s.Lock()
	defer s.Unlock()

	return s.store.Release(s)
}

// 保存Session所做的修改到Store中
func (s *Session) Save() error {
	s.Lock()
	defer s.Unlock()
	return s.store.Save(s)
}

// 获取最后的存取时间
func SessionAccessed(s *Session) time.Time {
	return s.accessed
}

// 获取Session中的数据。
// 供Store实现者调用
func SessionData(s *Session) map[interface{}]interface{} {
	return s.data
}

// 释放Session。
// 供Store实现者调用。
func FreeSession(s *Session) {
	s.Lock()
	defer s.Unlock()

	s.data = nil
	s.store = nil
	sessFree.Put(s)
}

// 新建Session。一般由store的实现者调用。
func NewSession(sid string, data map[interface{}]interface{}, s Store) *Session {
	sess := sessFree.Get().(*Session)
	sess.id = sid
	sess.data = data
	sess.store = s
	sess.accessed = time.Now()

	return sess
}

// 产生一个唯一的SessionID
func sessionID() (string, error) {
	ret := make([]byte, 64)
	n, err := io.ReadFull(rand.Reader, ret)
	if n == 0 {
		return "", errors.New("未读取到随机数")
	}

	h := md5.New()
	h.Write(ret)
	return hex.EncodeToString(h.Sum(nil)), err
}
