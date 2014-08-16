// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
	"net/url"
	"time"

	e "github.com/caixw/lib.go/errors"
)

type Instance struct {
	store      Store
	sessIDName string
	secure     bool
	lifetime   int
	ticker     *time.Ticker
}

// 实例一个新的Instance。会自动开始GC操作。
//
// sessionIDName：用于保存sessionid的变量名称，如果用cookie传递，则为cookie的名
// 称，不能包含特殊字符；
// lifetime：sessionid的生存时间，如果是cookie传递，则为cookie的max-age属性；
// secure: 是否只能用于https、ssl等安全链接，若为cookie传递，则为cookie的secure属性。
func New(store Store, sessionIDName string, lifetime int, secure bool) *Instance {
	inst := &Instance{
		sessIDName: sessionIDName,
		lifetime:   lifetime,
		secure:     secure,
		store:      store,
	}

	inst.gc()

	return inst
}

// 启动GC操作
func (i *Instance) gc() {
	if i.ticker != nil {
		i.ticker.Stop()
	}
	i.ticker = time.NewTicker(time.Duration(i.lifetime))

	go func() {
		for {
			_, ok := <-i.ticker.C
			if !ok {
				break
			}
			i.store.GC(i.lifetime)
		}
	}()
}

func (i *Instance) StartSessionWithForm(r *http.Request) (*Session, error) {
	var sessid string
	var err error

	tmp, found := r.Form[i.sessIDName]
	if !found {
		sessid, err = sessionID()
	} else {
		sessid = tmp[0]
	}

	if err != nil {
		return nil, e.New(0, nil, "初始化Session失败")
	}

	if sess, err := i.store.Get(sessid); err != nil {
		return nil, e.New(0, nil, "初始化Session失败")
	} else {
		return sess, nil
	}
}

// 开始一个新的Session
// 应该在任何输出之前调用，否则不会输出成功
func (i *Instance) StartSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	var sessid string
	var err error

	cookie, cerr := r.Cookie(i.sessIDName)
	if cerr != nil || cookie.Value == "" { // 不存在，新建一个sessionid
		sessid, err = sessionID()
	} else {
		sessid, err = url.QueryUnescape(cookie.Value)
	}

	if err != nil {
		return nil, e.New(0, err, "无法初始化新的session")
	}

	sess, err := i.store.Get(sessid)
	if err != nil {
		return nil, e.New(0, err, "无法初始化新的session")
	}

	i.setCookie(r, w, sessid, i.lifetime)
	return sess, nil
}

// 释放当前协程下的Session
func (i *Instance) FreeSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(i.sessIDName)
	if err != nil {
		return err
	}
	if cookie.Value == "" {
		return nil
	}

	sessid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return e.Newf(0, err, "无法释放Session:sid=[%v]", sessid)
	}

	err = i.store.Delete(sessid)
	if err != nil {
		return e.Newf(0, err, "无法释放Session:sid=[%v]", sessid)
	}

	i.setCookie(r, w, sessid, -1)

	return nil
}

// 释放整个SessionManager及相关数据。
func (i *Instance) Free() {
	if i.ticker != nil {
		i.ticker.Stop()
	}
	i.store.Free()
}

// 设置相应的cookie值
func (i *Instance) setCookie(r *http.Request, w http.ResponseWriter, sessid string, maxAge int) {
	cookie := &http.Cookie{
		Name:     i.sessIDName,
		Value:    url.QueryEscape(sessid),
		Secure:   i.secure,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   maxAge,
		// TODO(caixw) ie8以下只支持Expires而不支持max_age。http1.0只有只有expires，
		// 而在http1.1中expires属于废弃的属性，max-age才是正规的。
		Expires: time.Now().Add(time.Second * time.Duration(maxAge)),
	}
	http.SetCookie(w, cookie)
}
