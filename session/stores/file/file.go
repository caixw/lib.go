// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// session的文件存储模式
package file

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/caixw/lib.go/session"
)

// TODO(caixw) 将普通的锁改为文件锁。

const (
	// 版本号
	Version = "0.1.2.140823"
	// 所有Session文件的权限
	mode = 0666
)

// implement session.Store
type store struct {
	sync.Mutex
	sessions map[string]*session.Session
	saveDir  string
}

var _ session.Store = &store{}

// 新建Store
//
// saveDir:session的保存路径，若目录不存在，会尝试创建。
func New(saveDir string) *store {
	_, err := os.Stat(saveDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(saveDir, mode)
		} else if !os.IsExist(err) {
			panic(err)
		}
	}

	// 确保最后个字符为os.PathSeparator
	lastRune := saveDir[len(saveDir)-1]
	if lastRune != os.PathSeparator && lastRune != '\\' {
		saveDir = saveDir + string(os.PathSeparator)
	}

	return &store{
		sessions: make(map[string]*session.Session),
		saveDir:  saveDir,
	}
}

// implement session.Store.Get()
func (s *store) Get(sid string) (*session.Session, error) {
	s.Lock()
	defer s.Unlock()

	filepath := s.saveDir + sid
	data := make(map[interface{}]interface{})

	f, err := os.Open(filepath)
	defer f.Close()

	if err != nil {
		if os.IsNotExist(err) {
			return session.NewSession(sid, data, s), nil
		}
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	de := gob.NewDecoder(bytes.NewBuffer(b))
	err = de.Decode(&data)
	if err != nil {
		return nil, err
	}

	return session.NewSession(sid, data, s), nil
}

// implement session.Store.Save()
func (s *store) Save(sess *session.Session) error {
	s.Lock()
	defer s.Unlock()

	str := bytes.NewBufferString("")
	de := gob.NewEncoder(str)
	de.Encode(session.SessionData(sess))

	filepath := s.saveDir + sess.ID()
	return ioutil.WriteFile(filepath, str.Bytes(), mode)
}

// implement session.Store.Delete()
func (s *store) Delete(sid string) error {
	s.Lock()
	defer s.Unlock()

	return os.Remove(s.saveDir + sid)
}

// implement session.Store.Release()
func (s *store) Release(sess *session.Session) error {
	s.Lock()
	defer s.Unlock()

	str := bytes.NewBufferString("")
	de := gob.NewEncoder(str)
	de.Encode(session.SessionData(sess))

	filepath := s.saveDir + sess.ID()
	return ioutil.WriteFile(filepath, str.Bytes(), mode)

	session.FreeSession(sess)
	return nil
}

// implement session.Store.GC()
func (s *store) GC(duration int) {
	s.Lock()
	defer s.Unlock()

	fs, err := ioutil.ReadDir(s.saveDir)
	if err != nil {
		panic(err)
	}

	t := time.Now().Unix() - int64(duration) // 过期时间
	for _, v := range fs {
		if v.ModTime().Unix() < t {
			os.Remove(s.saveDir + v.Name())
		}
	}
}

// implement session.Store.Free()
func (s *store) Free() {
	s.Lock()
	defer s.Unlock()

	os.RemoveAll(s.saveDir)
	os.MkdirAll(s.saveDir, mode)
}
