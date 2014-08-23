// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// session的内存存储模式
package memory

import (
	"sync"
	"time"

	"github.com/caixw/lib.go/session"
)

// 当前包的版本
const Version = "0.1.2.140823"

// implement session.Store
type store struct {
	sync.Mutex
	sessions map[string]*session.Session
}

var _ session.Store = &store{}

func New() *store {
	return &store{
		sessions: make(map[string]*session.Session),
	}
}

// implement session.Store.Get()
func (s *store) Get(sid string) (*session.Session, error) {
	s.Lock()
	defer s.Unlock()

	ret, found := s.sessions[sid]
	if found {
		return ret, nil
	}

	/* 声明新的session实例 */
	ret = session.NewSession(sid, make(map[interface{}]interface{}), s)

	s.sessions[sid] = ret
	return ret, nil
}

// implement session.Store.Save()
func (s *store) Save(sess *session.Session) error {
	// 本身就在内存中，无需多做什么操作
	return nil
}

// implement session.Store.Delete()
func (s *store) Delete(sid string) error {
	s.Lock()
	defer s.Unlock()

	sess, found := s.sessions[sid]
	if !found {
		return nil
	}

	session.FreeSession(sess)
	return nil
}

// implement session.Store.Release()
func (s *store) Release(sess *session.Session) error {
	return nil
}

// implement session.Store.GC()
func (s *store) GC(duration int) {
	s.Lock()
	defer s.Unlock()

	t := time.Now().Unix() - int64(duration) // 过期时间
	for k, v := range s.sessions {
		if session.SessionAccessed(v).Unix() < t {
			s.Delete(k)
		}
	}
}

// implement session.Store.Free()
func (s *store) Free() {
	s.Lock()
	defer s.Unlock()

	for _, v := range s.sessions {
		session.FreeSession(v)
	}

	s.sessions = nil
}
